package config

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const (
	realConfig = "../../config/system.yaml"
	realRoot   = "../.."
)

func loadReal(t *testing.T, mutate func(string) string) (*Loaded, error) {
	t.Helper()
	raw, err := os.ReadFile(realConfig)
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	if mutate != nil {
		mutated := mutate(s)
		if mutated == s {
			t.Fatal("mutation did not change the config; test text is stale")
		}
		s = mutated
	}
	return Parse([]byte(s), realRoot)
}

func hasIssue(issues []Issue, sev Severity, field string) bool {
	for _, i := range issues {
		if i.Severity == sev && i.Field == field {
			return true
		}
	}
	return false
}

func TestRealConfigLoadsWithoutErrors(t *testing.T) {
	l, err := Load(realConfig)
	if err != nil {
		t.Fatal(err)
	}
	issues := l.Validate()
	if HasErrors(issues) {
		t.Fatalf("shipped config has errors: %v", issues)
	}
	if hasIssue(issues, Warning, "vms.opnsense.enabled") {
		t.Errorf("shipped OPNsense switches should agree, got %v", issues)
	}
}

func TestKernelHardeningListForm(t *testing.T) {
	l, err := loadReal(t, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := KernelHardening{DmesgRestrict: true, KptrRestrict: true, RestrictNamespace: true, PanicOnOops: false}
	if got := l.Config.Security.KernelHardening; got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestKernelHardeningMapAndBadKey(t *testing.T) {
	var k KernelHardening
	if err := yaml.Unmarshal([]byte("dmesg_restrict: true\npanic_on_oops: true\n"), &k); err != nil {
		t.Fatal(err)
	}
	if !k.DmesgRestrict || !k.PanicOnOops {
		t.Errorf("mapping form not applied: %+v", k)
	}
	if err := yaml.Unmarshal([]byte("- bogus: true\n"), &k); err == nil {
		t.Error("expected an error for an unknown kernel_hardening key")
	}
}

func TestDeployOrder(t *testing.T) {
	l, err := loadReal(t, nil)
	if err != nil {
		t.Fatal(err)
	}
	got, err := l.Config.DeployOrder()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"opnsense", "clean", "gaming", "research"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v (tools is disabled)", got, want)
	}
}

func TestUnknownKeyRejected(t *testing.T) {
	_, err := loadReal(t, func(s string) string { return s + "\nbogus_key: 1\n" })
	if err == nil {
		t.Fatal("expected an error for an unknown top-level key")
	}
}

func TestBadEnum(t *testing.T) {
	l, err := loadReal(t, func(s string) string {
		return strings.Replace(s, "method: auto", "method: bogus", 1)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !hasIssue(l.Validate(), Error, "hardware.gpu.method") {
		t.Error("expected an error on hardware.gpu.method")
	}
}

func TestDependencyOnDisabledVM(t *testing.T) {
	l, err := loadReal(t, func(s string) string {
		s = strings.Replace(s, "    enabled: true   # Keep aligned with microvm.opnsense_firewall.enabled", "    enabled: false  # Keep aligned with microvm.opnsense_firewall.enabled", 1)
		return strings.Replace(s, "  opnsense_firewall:\n    enabled: true", "  opnsense_firewall:\n    enabled: false", 1)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !hasIssue(l.Validate(), Error, "vms.clean.depends_on") {
		t.Error("expected clean to complain that opnsense is disabled")
	}
}

func TestDependencyCycle(t *testing.T) {
	l, err := loadReal(t, func(s string) string {
		return strings.Replace(s,
			"    enabled: true   # Keep aligned with microvm.opnsense_firewall.enabled\n",
			"    enabled: true   # Keep aligned with microvm.opnsense_firewall.enabled\n    depends_on: clean\n", 1)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !hasIssue(l.Validate(), Error, "vms") {
		t.Error("expected a dependency cycle error")
	}
}

func TestMissingTemplate(t *testing.T) {
	l, err := loadReal(t, func(s string) string {
		return strings.Replace(s, "config/vm-templates/clean-vm.yaml", "config/vm-templates/nope.yaml", 1)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !hasIssue(l.Validate(), Error, "vms.clean.template_file") {
		t.Error("expected a missing template error for clean")
	}
}

func TestDuplicateVaultPath(t *testing.T) {
	l, err := loadReal(t, func(s string) string {
		return strings.Replace(s, "path: /var/lib/vz-clean", "path: /var/lib/vz-gaming", 1)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !hasIssue(l.Validate(), Error, "storage_encryption.vaults.gaming_vm.path") {
		t.Errorf("expected a duplicate vault path error, got %v", l.Validate())
	}
}

func TestAutoInt(t *testing.T) {
	type doc struct {
		C AutoInt `yaml:"c"`
	}
	cases := []struct {
		in      string
		want    AutoInt
		wantErr bool
	}{
		{"c: auto", AutoInt{Auto: true}, false},
		{"c: 8", AutoInt{Value: 8}, false},
		{"c: 0", AutoInt{}, true},
		{"c: many", AutoInt{}, true},
	}
	for _, tc := range cases {
		var d doc
		err := yaml.Unmarshal([]byte(tc.in), &d)
		if (err != nil) != tc.wantErr {
			t.Errorf("%q: err=%v, wantErr=%v", tc.in, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && d.C != tc.want {
			t.Errorf("%q: got %+v, want %+v", tc.in, d.C, tc.want)
		}
	}
}

func TestAutoBool(t *testing.T) {
	type doc struct {
		C AutoBool `yaml:"c"`
	}
	cases := []struct {
		in      string
		want    AutoBool
		wantErr bool
	}{
		{"c: auto", AutoBool{Auto: true}, false},
		{"c: true", AutoBool{Value: true}, false},
		{"c: false", AutoBool{}, false},
		{"c: maybe", AutoBool{}, true},
	}
	for _, tc := range cases {
		var d doc
		err := yaml.Unmarshal([]byte(tc.in), &d)
		if (err != nil) != tc.wantErr {
			t.Errorf("%q: err=%v, wantErr=%v", tc.in, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && d.C != tc.want {
			t.Errorf("%q: got %+v, want %+v", tc.in, d.C, tc.want)
		}
	}
}

func TestStringList(t *testing.T) {
	type doc struct {
		D StringList `yaml:"d"`
	}
	var one, many doc
	if err := yaml.Unmarshal([]byte("d: opnsense"), &one); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string(one.D), []string{"opnsense"}) {
		t.Errorf("scalar form: got %v", one.D)
	}
	if err := yaml.Unmarshal([]byte("d: [a, b]"), &many); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string(many.D), []string{"a", "b"}) {
		t.Errorf("list form: got %v", many.D)
	}
}
