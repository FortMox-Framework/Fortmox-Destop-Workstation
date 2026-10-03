// Package config defines the typed schema for config/system.yaml and the VM
// template fields currently consumed by the CLI.
package config

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config mirrors config/system.yaml. Unknown keys are rejected at load time.
type Config struct {
	System            System            `yaml:"system"`
	Hardware          Hardware          `yaml:"hardware"`
	Wayland           Wayland           `yaml:"wayland"`
	Power             Power             `yaml:"power"`
	Security          Security          `yaml:"security"`
	Network           Network           `yaml:"network"`
	StorageEncryption StorageEncryption `yaml:"storage_encryption"`
	MicroVM           MicroVM           `yaml:"microvm"`
	VMs               map[string]VM     `yaml:"vms"`
	Advanced          Advanced          `yaml:"advanced"`
	Logging           Logging           `yaml:"logging"`
	Backup            Backup            `yaml:"backup"`
	Development       Development       `yaml:"development"`
}

// System contains host identity and form-factor settings.
type System struct {
	Type     string `yaml:"type"` // auto, laptop, desktop
	Hostname string `yaml:"hostname"`
	Timezone string `yaml:"timezone"`
	Debug    bool   `yaml:"debug"`
}

// Hardware groups CPU, GPU, and storage settings.
type Hardware struct {
	CPU     CPU     `yaml:"cpu"`
	GPU     GPU     `yaml:"gpu"`
	Storage Storage `yaml:"storage"`
}

// CPU contains processor, IOMMU, and scaling settings.
type CPU struct {
	Cores           AutoInt `yaml:"cores"`
	IOMMU           bool    `yaml:"iommu"`
	ScalingGovernor string  `yaml:"scaling_governor"` // auto, balanced, performance, powersave
}

// GPU describes passthrough intent and selected devices.
type GPU struct {
	Passthrough bool     `yaml:"passthrough"`
	Method      string   `yaml:"method"` // auto, virtio, sr-iov, vgpu, full
	Devices     []string `yaml:"devices"`
}

// Storage contains host storage encryption and optimization settings.
type Storage struct {
	Encryption     bool   `yaml:"encryption"`
	EncryptionType string `yaml:"encryption_type"` // luks2, luks1
	Compression    bool   `yaml:"compression"`
	TrimEnabled    bool   `yaml:"trim_enabled"`
}

// Wayland contains host compositor and display settings.
type Wayland struct {
	Compositor string `yaml:"compositor"` // auto, hikari, woodland, dwl
	Stacking   bool   `yaml:"stacking"`
	Resolution string `yaml:"resolution"` // auto or WxH
	HDR        bool   `yaml:"hdr"`
}

// Power contains CPU and device power-management preferences.
type Power struct {
	Profile            string   `yaml:"profile"` // auto, balanced, performance
	TurboBoost         AutoBool `yaml:"turbo_boost"`
	GPUDynamicPower    bool     `yaml:"gpu_dynamic_power"`
	AdaptiveBrightness bool     `yaml:"adaptive_brightness"`
	IdleStates         string   `yaml:"idle_states"` // auto, balanced, deep, performance
}

// Security contains host hardening and firewall settings.
type Security struct {
	HardeningLevel  string          `yaml:"hardening_level"` // paranoid, standard
	SELinux         bool            `yaml:"selinux"`
	AppArmor        bool            `yaml:"apparmor"`
	KernelHardening KernelHardening `yaml:"kernel_hardening"`
	Firewall        HostFirewall    `yaml:"firewall"`
}

// KernelHardening accepts both the list-of-single-key-maps form used in
// system.yaml today and a plain mapping.
type KernelHardening struct {
	DmesgRestrict     bool `yaml:"dmesg_restrict"`
	KptrRestrict      bool `yaml:"kptr_restrict"`
	RestrictNamespace bool `yaml:"restrict_namespaces"`
	PanicOnOops       bool `yaml:"panic_on_oops"`
}

// UnmarshalYAML accepts either a mapping or a list of single-key mappings.
func (k *KernelHardening) UnmarshalYAML(n *yaml.Node) error {
	type plain KernelHardening
	switch n.Kind {
	case yaml.MappingNode:
		return decodeStrict(n, (*plain)(k))
	case yaml.SequenceNode:
		for _, item := range n.Content {
			if item.Kind != yaml.MappingNode {
				return fmt.Errorf("line %d: kernel_hardening entries must be key: value pairs", item.Line)
			}
			if err := decodeStrict(item, (*plain)(k)); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("line %d: kernel_hardening must be a list or mapping", n.Line)
}

// HostFirewall contains default traffic policies and management access rules.
type HostFirewall struct {
	InputPolicy            string `yaml:"input_policy"`   // DROP, REJECT, ACCEPT
	OutputPolicy           string `yaml:"output_policy"`  // DROP, REJECT, ACCEPT
	ForwardPolicy          string `yaml:"forward_policy"` // DROP, REJECT, ACCEPT
	SSHLocalhostOnly       bool   `yaml:"ssh_localhost_only"`
	ProxmoxUILocalhostOnly bool   `yaml:"proxmox_ui_localhost_only"`
}

// Network contains VM isolation and DNS settings.
type Network struct {
	IsolationMode string   `yaml:"isolation_mode"` // shared, isolated, air-gap
	DNSServers    []string `yaml:"dns_servers"`
}

// StorageEncryption contains encryption parameters and named vault definitions.
type StorageEncryption struct {
	Cipher          string           `yaml:"cipher"`
	KeySize         int              `yaml:"key_size"` // 256 or 512
	PBKDFIterations int              `yaml:"pbkdf_iterations"`
	Vaults          map[string]Vault `yaml:"vaults"`
}

// Vault describes one encrypted storage area.
type Vault struct {
	Enabled                bool   `yaml:"enabled"`
	Path                   string `yaml:"path"`
	Access                 string `yaml:"access"`
	EncryptionKeyProtected bool   `yaml:"encryption_key_protected"`
}

// MicroVM contains lightweight firewall VM settings.
type MicroVM struct {
	Enabled             bool                `yaml:"enabled"`
	OPNsenseFirewall    OPNsenseFirewall    `yaml:"opnsense_firewall"`
	AdditionalFirewalls AdditionalFirewalls `yaml:"additional_firewalls"`
}

// OPNsenseFirewall configures the OPNsense firewall VM.
type OPNsenseFirewall struct {
	Enabled      bool   `yaml:"enabled"`
	VMID         int    `yaml:"vmid"`
	TemplateFile string `yaml:"template_file"`
	Priority     int    `yaml:"priority"`
}

// AdditionalFirewalls contains optional additional firewall instances.
type AdditionalFirewalls struct {
	Enabled   bool  `yaml:"enabled"`
	Instances []any `yaml:"instances"`
}

// VM is one entry under vms:. DependsOn accepts a single name or a list.
type VM struct {
	Enabled      bool       `yaml:"enabled"`
	TemplateFile string     `yaml:"template_file"`
	DependsOn    StringList `yaml:"depends_on"`
}

// Advanced contains optional host features and experimental settings.
type Advanced struct {
	ModuleLoadingRestrictions bool `yaml:"module_loading_restrictions"`
	SecureBoot                bool `yaml:"secure_boot"`
	AuditLogging              bool `yaml:"audit_logging"`
	Watchdog                  bool `yaml:"watchdog"`
}

// Logging contains host logging and audit preferences.
type Logging struct {
	Level           string `yaml:"level"` // debug, info, warn, error
	Syslog          bool   `yaml:"syslog"`
	Journal         bool   `yaml:"journal"`
	AuditOperations bool   `yaml:"audit_operations"`
}

// Backup contains scheduled backup settings.
type Backup struct {
	Enabled     bool   `yaml:"enabled"`
	Schedule    string `yaml:"schedule"` // daily, weekly, monthly
	Destination string `yaml:"destination"`
	Retention   int    `yaml:"retention"`
}

// Development contains testing and dry-run controls.
type Development struct {
	TestingMode bool `yaml:"testing_mode"`
	DryRun      bool `yaml:"dry_run"`
}

// AutoInt is either the word "auto" or a positive integer.
type AutoInt struct {
	Auto  bool
	Value int
}

// UnmarshalYAML accepts the scalar "auto" or a positive integer.
func (a *AutoInt) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.ScalarNode {
		return fmt.Errorf("line %d: expected \"auto\" or a number", n.Line)
	}
	if strings.EqualFold(n.Value, "auto") {
		*a = AutoInt{Auto: true}
		return nil
	}
	v, err := strconv.Atoi(n.Value)
	if err != nil || v < 1 {
		return fmt.Errorf("line %d: expected \"auto\" or a positive number, got %q", n.Line, n.Value)
	}
	*a = AutoInt{Value: v}
	return nil
}

// AutoBool is "auto", true or false.
type AutoBool struct {
	Auto  bool
	Value bool
}

// UnmarshalYAML accepts the scalar "auto", true, or false.
func (a *AutoBool) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.ScalarNode {
		return fmt.Errorf("line %d: expected auto, true or false", n.Line)
	}
	switch strings.ToLower(n.Value) {
	case "auto":
		*a = AutoBool{Auto: true}
	case "true":
		*a = AutoBool{Value: true}
	case "false":
		*a = AutoBool{Value: false}
	default:
		return fmt.Errorf("line %d: expected auto, true or false, got %q", n.Line, n.Value)
	}
	return nil
}

// StringList accepts either a single string or a list of strings.
type StringList []string

// UnmarshalYAML accepts either one dependency name or a sequence of names.
func (s *StringList) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.ScalarNode:
		if n.Value == "" {
			*s = nil
			return nil
		}
		*s = StringList{n.Value}
		return nil
	case yaml.SequenceNode:
		out := make(StringList, 0, len(n.Content))
		for _, c := range n.Content {
			if c.Kind != yaml.ScalarNode {
				return fmt.Errorf("line %d: expected a list of names", c.Line)
			}
			out = append(out, c.Value)
		}
		*s = out
		return nil
	}
	return fmt.Errorf("line %d: expected a name or a list of names", n.Line)
}

// decodeStrict decodes a node into v and rejects unknown keys.
func decodeStrict(n *yaml.Node, v any) error {
	// A node cannot be decoded strictly directly, so re-encode and decode.
	raw, err := yaml.Marshal(n)
	if err != nil {
		return err
	}
	return decodeBytesStrict(raw, v)
}
