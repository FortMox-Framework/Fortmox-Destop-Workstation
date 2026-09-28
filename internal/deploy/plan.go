// Package deploy describes host and VM actions implied by the configuration.
package deploy

	import "github.com/FortMox-Framework/Fortmox-Destop-Workstation/internal/config"

// Plan returns the host and VM actions implied by the loaded configuration.
func Plan(loaded *config.Loaded) []string {
	steps := []string{
		"validate host requirements and detect hardware",
		"configure kernel IOMMU and hardening settings",
		"prepare configured storage encryption vaults",
		"configure isolated virtual networking and firewall",
	}
	if loaded.Config.Hardware.GPU.Passthrough {
		steps = append(steps, "prepare GPU passthrough using method "+loaded.Config.Hardware.GPU.Method)
	}
	order, err := loaded.Config.DeployOrder()
	if err == nil {
		for _, name := range order {
			steps = append(steps, "create configured VM "+name)
		}
	}
	return steps
}
