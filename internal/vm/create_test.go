package vm

import (
	"reflect"
	"testing"

	"github.com/FortMoxDesktop-Team/FortmoxDestop-Workstation/internal/config"
)

func TestCreateArgsReadTemplateResources(t *testing.T) {
	template, err := config.LoadTemplate("../..", "config/vm-templates/clean-vm.yaml")
	if err != nil {
		t.Fatal(err)
	}
	got, err := CreateArgs("clean", template)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"create", "100", "--name", "clean", "--cores", "4", "--memory", "8192", "--scsi0", "local-lvm:50G", "--net0", "virtio,bridge=vmbr0,firewall=1", "--ostype", "l26", "--agent", "1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
