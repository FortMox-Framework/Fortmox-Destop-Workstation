// Package power plans and applies CPU frequency governor profiles.
package power

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Change identifies one CPU governor file and the value selected for it.
type Change struct {
	Path     string
	Governor string
}

// Plan resolves a profile to governors supported by every CPU policy.
func Plan(profile string) ([]Change, error) {
	if profile != "performance" && profile != "balanced" && profile != "powersave" {
		return nil, fmt.Errorf("unknown power profile %q (choose performance, balanced or powersave)", profile)
	}
	paths, err := filepath.Glob("/sys/devices/system/cpu/cpu[0-9]*/cpufreq/scaling_governor")
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no CPU frequency governor controls found")
	}
	changes := make([]Change, 0, len(paths))
	for _, path := range paths {
		available := ""
		if raw, err := os.ReadFile(filepath.Join(filepath.Dir(path), "scaling_available_governors")); err == nil {
			available = string(raw)
		}
		governor := selectGovernor(profile, strings.Fields(available))
		if governor == "" {
			return nil, fmt.Errorf("no supported governor found for %s", path)
		}
		changes = append(changes, Change{Path: path, Governor: governor})
	}
	return changes, nil
}

// selectGovernor chooses the first supported governor for the requested profile.
func selectGovernor(profile string, available []string) string {
	candidates := map[string][]string{
		"performance": {"performance"},
		"balanced":    {"schedutil", "ondemand", "powersave"},
		"powersave":   {"powersave", "conservative"},
	}
	// Some drivers expose only performance and powersave; use powersave as the balanced fallback.
	for _, candidate := range candidates[profile] {
		for _, supported := range available {
			if candidate == supported {
				return candidate
			}
		}
	}
	if len(available) == 0 {
		return profile
	}
	return ""
}

// Apply writes the planned governors; callers must enforce their privilege policy.
func Apply(changes []Change) error {
	for _, change := range changes {
		if err := os.WriteFile(change.Path, []byte(change.Governor+"\n"), 0); err != nil {
			return fmt.Errorf("set %s to %s: %w", change.Path, change.Governor, err)
		}
	}
	return nil
}
