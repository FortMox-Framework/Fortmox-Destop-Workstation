package verify

import (
	"reflect"
	"testing"
)

func TestParseVMIDs(t *testing.T) {
	got := parseVMIDs(" VMID NAME STATUS\n50 firewall running\n100 clean stopped\n")
	want := map[int]bool{50: true, 100: true}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
