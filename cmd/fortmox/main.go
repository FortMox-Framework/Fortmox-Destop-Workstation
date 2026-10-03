// Command fortmox manages a FortMox Proxmox workstation.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/FortMox-Framework/Fortmox-Destop-Workstation/internal/config"
	"github.com/FortMox-Framework/Fortmox-Destop-Workstation/internal/deploy"
	"github.com/FortMox-Framework/Fortmox-Destop-Workstation/internal/gpu"
	"github.com/FortMox-Framework/Fortmox-Destop-Workstation/internal/hardware"
	"github.com/FortMox-Framework/Fortmox-Destop-Workstation/internal/power"
	"github.com/FortMox-Framework/Fortmox-Destop-Workstation/internal/verify"
	"github.com/FortMox-Framework/Fortmox-Destop-Workstation/internal/vm"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// newRootCommand constructs the CLI tree and shares the config flag with commands.
func newRootCommand() *cobra.Command {
	var configPath string
	root := &cobra.Command{
		Use:           "fortmox",
		Short:         "Manage a FortMox Proxmox workstation",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().StringVarP(&configPath, "config", "c", "config/system.yaml", "path to system.yaml")
	root.AddCommand(newConfigCommand(&configPath), newDetectCommand(), newVerifyCommand(&configPath), newDeployCommand(&configPath), newVMCommand(&configPath), newGPUCommand(&configPath), newPowerCommand(&configPath))
	return root
}

func newConfigCommand(path *string) *cobra.Command {
	command := &cobra.Command{Use: "config", Short: "Validate and plan from system.yaml"}
	command.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "Validate the system config and VM templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			loaded, err := config.Load(*path)
			if err != nil {
				return err
			}
			issues := loaded.Validate()
			for _, issue := range issues {
				fmt.Fprintln(cmd.OutOrStdout(), issue)
			}
			if config.HasErrors(issues) {
				return fmt.Errorf("configuration contains errors")
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s: ok (%d warning(s))\n", *path, len(issues))
			return nil
		},
	}, newPlanCommand(path))
	return command
}

func newPlanCommand(path *string) *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Show enabled VMs in dependency order",
		RunE: func(cmd *cobra.Command, args []string) error {
			loaded, err := loadValidConfig(*path)
			if err != nil {
				return err
			}
			order, err := loaded.Config.DeployOrder()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Deploy order:")
			for index, name := range order {
				template, loadErr := config.LoadTemplate(loaded.Root, loaded.Config.TemplateFile(name))
				if loadErr != nil {
					return loadErr
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s  vmid %d, %d cores, %d-%d MB RAM, %s disk\n", index+1, name, template.VM.VMID, template.VM.Specs.CPU.Cores, template.VM.Specs.Memory.Min, template.VM.Specs.Memory.Max, template.VM.Specs.Disk.Size)
			}
			return nil
		},
	}
}

func newDetectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Report CPU, GPU, IOMMU and form-factor information",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := hardware.Detect()
			fmt.Fprintf(cmd.OutOrStdout(), "Form factor: %s\nCPU: %s (%d logical CPUs)\nVirtualization: %t\nIOMMU: %t\n", result.FormFactor, result.CPUModel, result.LogicalCPUs, result.Virtualization, result.IOMMU)
			if len(result.GPUs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "GPUs: none detected (lspci unavailable or no display devices)")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "GPUs:")
				for _, gpu := range result.GPUs {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", gpu)
				}
			}
			return nil
		},
	}
}

func newVerifyCommand(path *string) *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Run read-only host, template, VM and firewall checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			loaded, err := config.Load(*path)
			if err != nil {
				return err
			}
			results := verify.Run(loaded)
			failed := false
			for _, result := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s\n", strings.ToUpper(result.Status), result.Message)
				failed = failed || result.Status == "fail"
			}
			if failed {
				return fmt.Errorf("verification found failures")
			}
			return nil
		},
	}
}

func newDeployCommand(path *string) *cobra.Command {
	var dryRun bool
	command := &cobra.Command{
		Use:   "deploy",
		Short: "Plan host changes without applying them (apply support is not implemented)",
		RunE: func(cmd *cobra.Command, args []string) error {
			loaded, err := loadValidConfig(*path)
			if err != nil {
				return err
			}
			for _, step := range deploy.Plan(loaded) {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", step)
			}
			if !dryRun {
				return fmt.Errorf("host mutation is not available in this CLI build; re-run with --dry-run to review the plan")
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Dry run only; no host changes made.")
			return nil
		},
	}
	command.Flags().BoolVar(&dryRun, "dry-run", true, "show planned operations without changing the host")
	return command
}

func newVMCommand(path *string) *cobra.Command {
	var dryRun bool
	create := &cobra.Command{
		Use:   "create [vm-name]",
		Short: "Create enabled VMs using the supported template fields",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			loaded, err := loadValidConfig(*path)
			if err != nil {
				return err
			}
			names, err := loaded.Config.DeployOrder()
			if err != nil {
				return err
			}
			if len(args) == 1 {
				names = []string{args[0]}
			}
			for _, name := range names {
				if !loaded.Config.VMEnabled(name) {
					return fmt.Errorf("VM %q is disabled in config", name)
				}
				template, err := config.LoadTemplate(loaded.Root, loaded.Config.TemplateFile(name))
				if err != nil {
					return err
				}
				args, err := vm.CreateArgs(name, template)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "qm %s\n", strings.Join(args, " "))
				if !dryRun {
					if _, err := exec.LookPath("qm"); err != nil {
						return fmt.Errorf("qm is required to create VMs")
					}
					if output, err := exec.Command("qm", args...).CombinedOutput(); err != nil {
						return fmt.Errorf("qm %s: %s: %w", strings.Join(args, " "), strings.TrimSpace(string(output)), err)
					}
				}
			}
			if dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), "Dry run only; no VMs created.")
			}
			return nil
		},
	}
	create.Flags().BoolVar(&dryRun, "dry-run", true, "print qm commands without creating VMs")
	parent := &cobra.Command{Use: "vm", Short: "Manage configured virtual machines"}
	parent.AddCommand(create)
	return parent
}

func newGPUCommand(path *string) *cobra.Command {
	var method string
	command := &cobra.Command{
		Use:   "gpu",
		Short: "Plan GPU passthrough using configured devices and method",
		RunE: func(cmd *cobra.Command, args []string) error {
			loaded, err := loadValidConfig(*path)
			if err != nil {
				return err
			}
			if !loaded.Config.Hardware.GPU.Passthrough {
				fmt.Fprintln(cmd.OutOrStdout(), "GPU passthrough is disabled in system.yaml; no GPU actions planned.")
				return nil
			}
			selected := loaded.Config.Hardware.GPU.Method
			if method != "" {
				selected = method
			}
			devices := loaded.Config.Hardware.GPU.Devices
			if len(devices) == 0 {
				devices = hardware.Detect().GPUs
			}
			plan, err := gpu.Plan(selected, devices)
			if err != nil {
				return err
			}
			for _, line := range plan {
				fmt.Fprintln(cmd.OutOrStdout(), line)
			}
			return nil
		},
	}
	command.Flags().StringVar(&method, "method", "", "override method: auto, virtio, sr-iov, vgpu, full")
	return command
}

func newPowerCommand(path *string) *cobra.Command {
	var apply bool
	command := &cobra.Command{
		Use:   "power",
		Short: "Plan or apply a CPU frequency governor profile",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			loaded, err := loadValidConfig(*path)
			if err != nil {
				return err
			}
			profile := loaded.Config.Power.Profile
			if profile == "auto" {
				profile = "performance"
				if hardware.Detect().FormFactor == "laptop" {
					profile = "balanced"
				}
			}
			if len(args) == 1 {
				profile = args[0]
			}
			changes, err := power.Plan(profile)
			if err != nil {
				return err
			}
			for _, change := range changes {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", change.Path, change.Governor)
			}
			if !apply {
				fmt.Fprintln(cmd.OutOrStdout(), "Dry run only; pass --apply to change governors.")
				return nil
			}
			if os.Geteuid() != 0 {
				return fmt.Errorf("--apply requires root")
			}
			return power.Apply(changes)
		},
	}
	command.Flags().BoolVar(&apply, "apply", false, "apply the profile (requires root)")
	command.Use += " [performance|balanced|powersave]"
	return command
}

func loadValidConfig(path string) (*config.Loaded, error) {
	loaded, err := config.Load(path)
	if err != nil {
		return nil, err
	}
	issues := loaded.Validate()
	if config.HasErrors(issues) {
		for _, issue := range issues {
			fmt.Fprintln(os.Stderr, issue)
		}
		return nil, fmt.Errorf("configuration contains errors; run fortmox config validate")
	}
	return loaded, nil
}
