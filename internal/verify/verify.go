// Package verify performs read-only config, host, VM, and firewall checks.
package verify

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/FortMox-Framework/Fortmox-Destop-Workstation/internal/config"
	"github.com/FortMox-Framework/Fortmox-Destop-Workstation/internal/hardware"
)

// Result describes one verification check and its pass, warn, or fail status.
type Result struct {
	Status  string
	Message string
}

// Run performs available local and Proxmox checks without changing system state.
func Run(loaded *config.Loaded) []Result {
	var results []Result
	add := func(status, format string, args ...any) {
		results = append(results, Result{Status: status, Message: fmt.Sprintf(format, args...)})
	}
	issues := loaded.Validate()
	for _, issue := range issues {
		status := "pass"
		if issue.Severity == config.Error {
			status = "fail"
		} else {
			status = "warn"
		}
		add(status, "config %s", issue)
	}
	if !config.HasErrors(issues) {
		add("pass", "system config schema and enabled VM templates are valid")
	}

	detected := hardware.Detect()
	if detected.Virtualization {
		add("pass", "CPU virtualization extensions are available")
	} else {
		add("warn", "CPU virtualization extensions were not detected")
	}
	if detected.IOMMU {
		add("pass", "IOMMU groups are present")
	} else if loaded.Config.Hardware.GPU.Passthrough {
		add("fail", "GPU passthrough is configured but no IOMMU groups were found")
	} else {
		add("warn", "no IOMMU groups were found")
	}
	if _, err := exec.LookPath("pveversion"); err != nil {
		add("warn", "Proxmox VE tools are unavailable; running in a non-Proxmox or incomplete environment")
	} else if output, err := exec.Command("pveversion", "-v").Output(); err != nil {
		add("fail", "pveversion failed: %v", err)
	} else {
		add("pass", "Proxmox VE detected (%s)", strings.TrimSpace(strings.SplitN(string(output), "\n", 2)[0]))
	}
	if _, err := exec.LookPath("pvesm"); err == nil {
		if output, err := exec.Command("pvesm", "status").CombinedOutput(); err != nil {
			add("fail", "Proxmox storage status failed: %s", strings.TrimSpace(string(output)))
		} else {
			add("pass", "Proxmox storage status is available")
		}
	} else {
		add("warn", "pvesm is unavailable; Proxmox storage checks were skipped")
	}
	if _, err := exec.LookPath("ip"); err == nil {
		if output, err := exec.Command("ip", "link", "show", "vmbr0").CombinedOutput(); err != nil {
			add("warn", "default VM bridge vmbr0 was not found: %s", strings.TrimSpace(string(output)))
		} else {
			add("pass", "default VM bridge vmbr0 exists")
		}
	} else {
		add("warn", "ip is unavailable; network bridge check was skipped")
	}

	all := 0
	for _, file := range loaded.Config.VMs {
		if file.TemplateFile != "" {
			all++
			if _, err := os.Stat(filepath.Join(loaded.Root, file.TemplateFile)); err != nil {
				add("fail", "VM template missing: %s", file.TemplateFile)
			} else {
				add("pass", "VM template exists: %s", file.TemplateFile)
			}
		}
	}
	if all < 5 {
		add("warn", "only %d VM templates are configured; expected five workload templates", all)
	}

	qmPath, qmErr := exec.LookPath("qm")
	if qmErr != nil {
		add("warn", "qm is unavailable; Proxmox VM existence checks were skipped")
	} else if output, err := exec.Command(qmPath, "list").Output(); err != nil {
		add("fail", "cannot list Proxmox VMs: %v", err)
	} else {
		ids := parseVMIDs(string(output))
		for name := range loaded.Config.VMs {
			if !loaded.Config.VMEnabled(name) {
				continue
			}
			template, err := config.LoadTemplate(loaded.Root, loaded.Config.TemplateFile(name))
			if err != nil {
				continue
			}
			if ids[template.VM.VMID] {
				add("pass", "enabled VM %s (vmid %d) exists", name, template.VM.VMID)
			} else {
				add("fail", "enabled VM %s (vmid %d) is missing", name, template.VM.VMID)
			}
		}
	}

	if loaded.Config.MicroVM.Enabled && loaded.Config.MicroVM.OPNsenseFirewall.Enabled {
		firewallID := loaded.Config.MicroVM.OPNsenseFirewall.VMID
		if qmErr != nil {
			add("warn", "firewall VM %d status was not checked because qm is unavailable", firewallID)
		} else if output, err := exec.Command(qmPath, "status", strconv.Itoa(firewallID)).CombinedOutput(); err != nil {
			add("fail", "firewall VM %d is not available: %s", firewallID, strings.TrimSpace(string(output)))
		} else if !strings.Contains(string(output), "running") {
			add("fail", "firewall VM %d is not running", firewallID)
		} else {
			add("pass", "firewall VM %d is running", firewallID)
		}
		if _, err := exec.LookPath("pve-firewall"); err != nil {
			add("warn", "pve-firewall is unavailable; Proxmox firewall service was not checked")
		} else if output, err := exec.Command("pve-firewall", "status").CombinedOutput(); err != nil || !strings.Contains(strings.ToLower(string(output)), "running") {
			add("fail", "Proxmox firewall is not running: %s", strings.TrimSpace(string(output)))
		} else {
			add("pass", "Proxmox firewall service is running")
		}
	}
	return results
}

// parseVMIDs extracts numeric IDs from qm list output and ignores its header.
func parseVMIDs(output string) map[int]bool {
	ids := map[int]bool{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		id, err := strconv.Atoi(fields[0])
		if err == nil {
			ids[id] = true
		}
	}
	return ids
}
