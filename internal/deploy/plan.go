// Package deploy describes host and VM actions implied by the configuration.
package deploy

import "github.com/FortMox-Framework/Fortmox-Destop-Workstation/internal/config"

// Plan returns the host and VM actions implied by the loaded configuration.
func Plan(loaded *config.Loaded) []string {
	steps := []string{
		"host requirement checks and hardware detection are not implemented by deploy",
		"kernel IOMMU and host hardening changes are not implemented",
		"encrypted storage setup is not implemented",
		"network isolation and firewall configuration are not implemented",
	}
	if loaded.Config.Hardware.GPU.Passthrough {
		steps = append(steps, "GPU strategy planning only; no device changes for method "+loaded.Config.Hardware.GPU.Method)
	}
	order, err := loaded.Config.DeployOrder()
	if err == nil {
		for _, name := range order {
			steps = append(steps, "VM "+name+" is handled separately by vm create")
		}
	}
	return steps
}
