package config

import (
	"fmt"
	"sort"
)

// OPNsenseName is the key of the firewall VM in the vms map.
const OPNsenseName = "opnsense"

// VMEnabled reports whether a VM will be deployed. OPNsense is special: the
// config has two switches for it (vms.opnsense.enabled and
// microvm.opnsense_firewall.enabled); it deploys if either is on.
func (c *Config) VMEnabled(name string) bool {
	vm, ok := c.VMs[name]
	if name == OPNsenseName {
		return (ok && vm.Enabled) || (c.MicroVM.Enabled && c.MicroVM.OPNsenseFirewall.Enabled)
	}
	return ok && vm.Enabled
}

// TemplateFile returns the template path for a VM, falling back to the
// microvm section for OPNsense when vms.opnsense has none.
func (c *Config) TemplateFile(name string) string {
	if vm, ok := c.VMs[name]; ok && vm.TemplateFile != "" {
		return vm.TemplateFile
	}
	if name == OPNsenseName {
		return c.MicroVM.OPNsenseFirewall.TemplateFile
	}
	return ""
}

// DeployOrder returns the enabled VMs in an order that satisfies depends_on.
// Ties are broken alphabetically so the result is deterministic. It returns
// an error on a dependency cycle. Dependencies on disabled or unknown VMs are
// skipped here and reported by Validate instead.
func (c *Config) DeployOrder() ([]string, error) {
	var names []string
	for name := range c.VMs {
		if c.VMEnabled(name) {
			names = append(names, name)
		}
	}
	// OPNsense may be enabled only via microvm.* and absent from vms:.
	if _, ok := c.VMs[OPNsenseName]; !ok && c.VMEnabled(OPNsenseName) {
		names = append(names, OPNsenseName)
	}
	sort.Strings(names)

	enabled := map[string]bool{}
	for _, n := range names {
		enabled[n] = true
	}
	deps := func(n string) []string {
		if vm, ok := c.VMs[n]; ok {
			return vm.DependsOn
		}
		return nil
	}

	const (
		unseen = iota
		visiting
		done
	)
	state := map[string]int{}
	var order []string
	var visit func(n string, path []string) error
	visit = func(n string, path []string) error {
		switch state[n] {
		case done:
			return nil
		case visiting:
			return fmt.Errorf("dependency cycle: %v -> %s", path, n)
		}
		state[n] = visiting
		ds := append([]string(nil), deps(n)...)
		sort.Strings(ds)
		for _, d := range ds {
			if !enabled[d] {
				continue // reported by Validate, not a planning error
			}
			if err := visit(d, append(path, n)); err != nil {
				return err
			}
		}
		state[n] = done
		order = append(order, n)
		return nil
	}
	for _, n := range names {
		if err := visit(n, nil); err != nil {
			return nil, err
		}
	}
	return order, nil
}
