package gpu

import "testing"

func TestPlanAutoUsesConservativeVirtIO(t *testing.T) {
	plan, err := Plan("auto", []string{"gpu-0"})
	if err != nil {
		t.Fatal(err)
	}
	if plan[0] != "GPU method: virtio" {
		t.Fatalf("unexpected strategy: %v", plan)
	}
}

func TestPlanRejectsUnknownMethod(t *testing.T) {
	if _, err := Plan("unknown", nil); err == nil {
		t.Fatal("expected an unsupported method error")
	}
}

func TestRequiresIOMMUOnlyForPhysicalAssignmentMethods(t *testing.T) {
	for method, want := range map[string]bool{
		"auto":   false,
		"virtio": false,
		"sr-iov": true,
		"vgpu":   true,
		"full":   true,
	} {
		if got := RequiresIOMMU(method); got != want {
			t.Errorf("RequiresIOMMU(%q) = %t, want %t", method, got, want)
		}
	}
}
