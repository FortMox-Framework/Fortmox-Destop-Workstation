// Package hardware detects host CPU, GPU, IOMMU, and form-factor information.
package hardware

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Result contains hardware information detected from Linux host interfaces.
type Result struct {
	CPUModel       string
	LogicalCPUs    int
	Virtualization bool
	IOMMU          bool
	FormFactor     string
	GPUs           []string
}

// Detect reads CPU flags, IOMMU groups, DMI form factor, and PCI display devices.
func Detect() Result {
	result := Result{LogicalCPUs: runtime.NumCPU(), FormFactor: detectFormFactor()}
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			key, value, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			switch strings.TrimSpace(key) {
			case "model name", "Hardware":
				if result.CPUModel == "" {
					result.CPUModel = strings.TrimSpace(value)
				}
			case "flags", "Features":
				flags := " " + strings.TrimSpace(value) + " "
				result.Virtualization = result.Virtualization || strings.Contains(flags, " vmx ") || strings.Contains(flags, " svm ")
			}
		}
	}
	groups, _ := filepath.Glob("/sys/kernel/iommu_groups/[0-9]*")
	result.IOMMU = len(groups) > 0
	if output, err := exec.Command("lspci", "-nn").Output(); err == nil {
		for _, line := range strings.Split(string(output), "\n") {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "vga compatible controller") || strings.Contains(lower, "3d controller") || strings.Contains(lower, "display controller") {
				result.GPUs = append(result.GPUs, strings.TrimSpace(line))
			}
		}
	}
	return result
}

func detectFormFactor() string {
	if data, err := os.ReadFile("/sys/class/dmi/id/chassis_type"); err == nil {
		// DMI chassis codes 8-14 describe portable systems; 3-7 are common desktops.
		switch strings.TrimSpace(string(data)) {
		case "8", "9", "10", "11", "12", "14":
			return "laptop"
		case "3", "4", "5", "6", "7", "13", "15", "16", "17":
			return "desktop"
		}
	}
	if _, err := os.Stat("/sys/class/power_supply/BAT0"); err == nil {
		return "laptop"
	}
	return "unknown"
}
