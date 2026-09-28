package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loaded is a parsed config plus the project root that relative paths
// (such as template_file) are resolved against.
type Loaded struct {
	Config Config
	Root   string
	Path   string
}

// Load reads and strictly parses system.yaml. The project root is the parent
// of the config directory containing the file (config/system.yaml -> project root).
func Load(path string) (*Loaded, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(abs)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var c Config
	if err := decodeBytesStrict(raw, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &Loaded{
		Config: c,
		Root:   filepath.Dir(filepath.Dir(abs)),
		Path:   abs,
	}, nil
}

// Parse parses config from bytes with root as the project root. Useful in tests.
func Parse(raw []byte, root string) (*Loaded, error) {
	var c Config
	if err := decodeBytesStrict(raw, &c); err != nil {
		return nil, err
	}
	return &Loaded{Config: c, Root: root}, nil
}

// decodeBytesStrict rejects YAML keys that are absent from the destination schema.
func decodeBytesStrict(raw []byte, v any) error {
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	dec.KnownFields(true)
	if err := dec.Decode(v); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}
