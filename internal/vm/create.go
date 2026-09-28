// Package vm builds Proxmox VM commands from FortMox templates.
package vm

import (
	"fmt"
	"strings"

	"github.com/FortMoxDesktop-Team/FortmoxDestop-Workstation/internal/config"
)

// CreateArgs builds qm create arguments from the currently supported template fields.
// It maps VMID, name, CPU, memory, disk, and a default bridge only.
func CreateArgs(name string, template *config.Template) ([]string, error) {
	if template.VM.VMID < 1 {
		return nil, fmt.Errorf("template for %s has no valid vmid", name)
	}
	if template.VM.Specs.CPU.Cores < 1 || template.VM.Specs.Memory.Max < 1 {
		return nil, fmt.Errorf("template for %s needs positive CPU cores and memory", name)
	}
	disk := template.VM.Specs.Disk
	if disk.Storage == "" || disk.Size == "" {
		return nil, fmt.Errorf("template for %s needs disk storage and size", name)
	}
	vmName := strings.TrimSuffix(template.VM.Name, "-vm")
	if vmName == "" {
		vmName = name
	}
	return []string{"create", fmt.Sprint(template.VM.VMID), "--name", vmName, "--cores", fmt.Sprint(template.VM.Specs.CPU.Cores), "--memory", fmt.Sprint(template.VM.Specs.Memory.Max), "--scsi0", disk.Storage + ":" + disk.Size, "--net0", "virtio,bridge=vmbr0,firewall=1", "--ostype", "l26", "--agent", "1"}, nil
}
