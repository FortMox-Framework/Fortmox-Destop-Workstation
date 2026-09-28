package power

import "testing"

func TestSelectGovernorUsesAvailableProfile(t *testing.T) {
	if got := selectGovernor("balanced", []string{"performance", "schedutil"}); got != "schedutil" {
		t.Fatalf("got %q, want schedutil", got)
	}
	if got := selectGovernor("performance", []string{"powersave"}); got != "" {
		t.Fatalf("got %q, want no supported governor", got)
	}
}
