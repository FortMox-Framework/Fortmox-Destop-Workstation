package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Template holds the header fields of a VM template that the CLI needs for
// planning and validation. The rest of each template file is not typed yet;
// it is passed through when the VM is created.
type Template struct {
	VM struct {
		Name  string `yaml:"name"`
		VMID  int    `yaml:"vmid"`
		Specs struct {
			CPU struct {
				Cores   int `yaml:"cores"`
				Sockets int `yaml:"sockets"`
			} `yaml:"cpu"`
			Memory struct {
				Min int `yaml:"min"`
				Max int `yaml:"max"`
			} `yaml:"memory"`
			Disk struct {
				Size    string `yaml:"size"`
				Storage string `yaml:"storage"`
			} `yaml:"disk"`
		} `yaml:"specs"`
	} `yaml:"vm"`
}

// LoadTemplate reads a template file (relative paths resolve against root).
func LoadTemplate(root, file string) (*Template, error) {
	p := file
	if !filepath.IsAbs(p) {
		p = filepath.Join(root, file)
	}
	raw, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var t Template
	if err := yaml.Unmarshal(raw, &t); err != nil {
		return nil, fmt.Errorf("parse %s: %w", file, err)
	}
	return &t, nil
}
