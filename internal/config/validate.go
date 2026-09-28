package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Severity of a validation finding.
type Severity int

const (
	// Warning reports a validation concern that does not invalidate the config.
	Warning Severity = iota
	// Error reports a config issue that prevents safe use.
	Error
)

// String returns the lowercase display name of the severity.
func (s Severity) String() string {
	if s == Error {
		return "error"
	}
	return "warning"
}

// Issue is one validation finding. Field is the dotted config path.
type Issue struct {
	Severity Severity
	Field    string
	Message  string
}

// String formats an issue as severity, config field, and explanation.
func (i Issue) String() string {
	return fmt.Sprintf("%s: %s: %s", i.Severity, i.Field, i.Message)
}

// HasErrors reports whether any issue is an error.
func HasErrors(issues []Issue) bool {
	for _, i := range issues {
		if i.Severity == Error {
			return true
		}
	}
	return false
}

var hostnameRe = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)
var resolutionRe = regexp.MustCompile(`^[0-9]+x[0-9]+$`)

type validator struct {
	issues []Issue
}

func (v *validator) errf(field, format string, a ...any) {
	v.issues = append(v.issues, Issue{Error, field, fmt.Sprintf(format, a...)})
}

func (v *validator) warnf(field, format string, a ...any) {
	v.issues = append(v.issues, Issue{Warning, field, fmt.Sprintf(format, a...)})
}

func (v *validator) oneOf(field, val string, allowed ...string) {
	for _, a := range allowed {
		if val == a {
			return
		}
	}
	v.errf(field, "%q is not valid; use one of: %s", val, strings.Join(allowed, ", "))
}

// Validate checks values, cross-field rules and (when l.Root is set) that
// template files exist and agree with the config. It never touches the system.
func (l *Loaded) Validate() []Issue {
	c := &l.Config
	v := &validator{}

	// system
	v.oneOf("system.type", c.System.Type, "auto", "laptop", "desktop")
	if !hostnameRe.MatchString(c.System.Hostname) {
		v.errf("system.hostname", "%q is not a valid hostname", c.System.Hostname)
	}
	if c.System.Timezone == "" {
		v.errf("system.timezone", "must not be empty")
	} else if _, err := time.LoadLocation(c.System.Timezone); err != nil {
		v.errf("system.timezone", "%q is not a known timezone", c.System.Timezone)
	}

	// hardware
	v.oneOf("hardware.cpu.scaling_governor", c.Hardware.CPU.ScalingGovernor, "auto", "balanced", "performance", "powersave")
	v.oneOf("hardware.gpu.method", c.Hardware.GPU.Method, "auto", "virtio", "sr-iov", "vgpu", "full")
	v.oneOf("hardware.storage.encryption_type", c.Hardware.Storage.EncryptionType, "luks2", "luks1")
	if c.Hardware.GPU.Passthrough && !c.Hardware.CPU.IOMMU && c.Hardware.GPU.Method != "virtio" {
		v.errf("hardware.cpu.iommu", "GPU passthrough (method %q) needs IOMMU; set hardware.cpu.iommu: true or use method virtio", c.Hardware.GPU.Method)
	}

	// wayland
	v.oneOf("wayland.compositor", c.Wayland.Compositor, "auto", "hikari", "woodland", "dwl")
	if c.Wayland.Resolution != "auto" && !resolutionRe.MatchString(c.Wayland.Resolution) {
		v.errf("wayland.resolution", "%q must be \"auto\" or WIDTHxHEIGHT, for example 1920x1080", c.Wayland.Resolution)
	}

	// power
	v.oneOf("power.profile", c.Power.Profile, "auto", "balanced", "performance")
	v.oneOf("power.idle_states", c.Power.IdleStates, "auto", "balanced", "deep", "performance")

	// security
	v.oneOf("security.hardening_level", c.Security.HardeningLevel, "paranoid", "standard")
	fw := c.Security.Firewall
	v.oneOf("security.firewall.input_policy", fw.InputPolicy, "DROP", "REJECT", "ACCEPT")
	v.oneOf("security.firewall.output_policy", fw.OutputPolicy, "DROP", "REJECT", "ACCEPT")
	v.oneOf("security.firewall.forward_policy", fw.ForwardPolicy, "DROP", "REJECT", "ACCEPT")
	if c.Security.HardeningLevel == "paranoid" && fw.InputPolicy == "ACCEPT" {
		v.warnf("security.firewall.input_policy", "ACCEPT on input weakens hardening_level paranoid")
	}

	// network
	v.oneOf("network.isolation_mode", c.Network.IsolationMode, "shared", "isolated", "air-gap")
	for i, s := range c.Network.DNSServers {
		if net.ParseIP(s) == nil {
			v.errf(fmt.Sprintf("network.dns_servers[%d]", i), "%q is not an IP address", s)
		}
	}

	// storage_encryption
	se := c.StorageEncryption
	if se.KeySize != 256 && se.KeySize != 512 {
		v.errf("storage_encryption.key_size", "%d is not valid; use 256 or 512", se.KeySize)
	}
	if se.PBKDFIterations < 1000 {
		v.errf("storage_encryption.pbkdf_iterations", "%d is too low", se.PBKDFIterations)
	}
	v.vaults(se.Vaults)

	// logging, backup
	v.oneOf("logging.level", c.Logging.Level, "debug", "info", "warn", "error")
	if c.Backup.Enabled {
		v.oneOf("backup.schedule", c.Backup.Schedule, "daily", "weekly", "monthly")
		if c.Backup.Retention < 1 {
			v.errf("backup.retention", "must be at least 1 when backup is enabled")
		}
		if !filepath.IsAbs(c.Backup.Destination) {
			v.errf("backup.destination", "%q must be an absolute path", c.Backup.Destination)
		}
	}

	l.validateVMs(v)
	return v.issues
}

func (v *validator) vaults(vaults map[string]Vault) {
	names := make([]string, 0, len(vaults))
	for n := range vaults {
		names = append(names, n)
	}
	sort.Strings(names)
	seen := map[string]string{}
	for _, n := range names {
		vt := vaults[n]
		field := "storage_encryption.vaults." + n
		if !vt.Enabled {
			continue
		}
		if !filepath.IsAbs(vt.Path) {
			v.errf(field+".path", "%q must be an absolute path", vt.Path)
			continue
		}
		if other, dup := seen[vt.Path]; dup {
			v.errf(field+".path", "%q is already used by vault %q", vt.Path, other)
		}
		seen[vt.Path] = n
	}
}

func (l *Loaded) validateVMs(v *validator) {
	c := &l.Config

	// Two switches for OPNsense: warn when they disagree.
	if vm, ok := c.VMs[OPNsenseName]; ok {
		mv := c.MicroVM.Enabled && c.MicroVM.OPNsenseFirewall.Enabled
		if vm.Enabled != mv {
			v.warnf("vms.opnsense.enabled", "vms.opnsense.enabled (%t) and microvm.opnsense_firewall.enabled (%t) disagree; OPNsense deploys if either is on", vm.Enabled, mv)
		}
	}

	names := make([]string, 0, len(c.VMs))
	for n := range c.VMs {
		names = append(names, n)
	}
	sort.Strings(names)

	// depends_on must name known, enabled VMs.
	for _, n := range names {
		if !c.VMEnabled(n) {
			continue
		}
		for _, d := range c.VMs[n].DependsOn {
			field := fmt.Sprintf("vms.%s.depends_on", n)
			if _, known := c.VMs[d]; !known && d != OPNsenseName {
				v.errf(field, "unknown VM %q", d)
			} else if !c.VMEnabled(d) {
				v.errf(field, "%s is enabled but depends on %q, which is disabled", n, d)
			}
		}
	}
	if _, err := c.DeployOrder(); err != nil {
		v.errf("vms", "%v", err)
	}

	if c.MicroVM.Enabled && c.MicroVM.OPNsenseFirewall.Enabled && c.MicroVM.OPNsenseFirewall.VMID < 1 {
		v.errf("microvm.opnsense_firewall.vmid", "must be a positive number")
	}

	if l.Root == "" {
		return
	}

	// Templates: must exist, and their vmids must be unique.
	byID := map[int]string{}
	all := append([]string(nil), names...)
	if _, ok := c.VMs[OPNsenseName]; !ok && c.VMEnabled(OPNsenseName) {
		all = append(all, OPNsenseName)
	}
	for _, n := range all {
		if !c.VMEnabled(n) {
			continue
		}
		field := fmt.Sprintf("vms.%s.template_file", n)
		file := c.TemplateFile(n)
		if file == "" {
			v.errf(field, "no template_file set")
			continue
		}
		t, err := LoadTemplate(l.Root, file)
		if err != nil {
			if os.IsNotExist(err) {
				v.errf(field, "template not found: %s", file)
			} else {
				v.errf(field, "%v", err)
			}
			continue
		}
		if t.VM.VMID < 1 {
			v.errf(field, "%s has no vm.vmid", file)
			continue
		}
		if other, dup := byID[t.VM.VMID]; dup {
			v.errf(field, "vmid %d in %s is already used by %s", t.VM.VMID, file, other)
		}
		byID[t.VM.VMID] = file
		if n == OPNsenseName && c.MicroVM.OPNsenseFirewall.VMID != 0 && t.VM.VMID != c.MicroVM.OPNsenseFirewall.VMID {
			v.errf("microvm.opnsense_firewall.vmid", "config says %d but %s says %d", c.MicroVM.OPNsenseFirewall.VMID, file, t.VM.VMID)
		}
	}
}
