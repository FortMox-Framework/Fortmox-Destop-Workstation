// Package gpu plans GPU passthrough strategies without changing host state.
package gpu

import (
	"fmt"
	"strings"
)

var methods = map[string]string{
	"virtio": "virtual display device; no physical GPU assignment",
	"sr-iov": "assign a virtual function; requires device SR-IOV support and host configuration",
	"vgpu":   "assign a vendor vGPU; requires a supported GPU and licensed host driver",
	"full":   "assign a physical GPU to one VM through VFIO",
}

// RequiresIOMMU reports whether a method assigns a physical GPU or device function.
func RequiresIOMMU(method string) bool {
	if method == "auto" {
		method = "virtio"
	}
	return method == "sr-iov" || method == "vgpu" || method == "full"
}

// Plan describes a configured strategy without changing host or VM state.
func Plan(method string, devices []string) ([]string, error) {
	if method == "auto" {
		method = "virtio"
	}
	description, ok := methods[method]
	if !ok {
		return nil, fmt.Errorf("unsupported GPU method %q (choose auto, virtio, sr-iov, vgpu or full)", method)
	}
	plan := []string{"GPU method: " + method, "Strategy: " + description}
	if len(devices) == 0 {
		plan = append(plan, "Devices: none configured or detected")
	} else {
		plan = append(plan, "Devices:")
		for _, device := range devices {
			plan = append(plan, "  "+strings.TrimSpace(device))
		}
	}
	if method == "full" && len(devices) > 1 {
		plan = append(plan, "Warning: full passthrough normally assigns each physical GPU to one VM; review IOMMU groups before applying")
	}
	return plan, nil
}
