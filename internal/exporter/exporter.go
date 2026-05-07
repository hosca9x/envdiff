// Package exporter provides functionality to export env maps to various file formats.
package exporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Format represents the output format for export.
type Format int

const (
	FormatEnv  Format = iota // KEY=VALUE format
	FormatExport             // export KEY=VALUE format (shell-compatible)
	FormatDocker             // Docker --env-file compatible format
)

// Exporter writes env maps to an output destination.
type Exporter struct {
	format Format
	quote  bool
}

// New returns an Exporter using the default env format.
func New() *Exporter {
	return &Exporter{format: FormatEnv, quote: false}
}

// NewWithOptions returns an Exporter with the given format and quoting preference.
func NewWithOptions(format Format, quote bool) *Exporter {
	return &Exporter{format: format, quote: quote}
}

// Write writes the env map to w in the configured format.
// Keys are written in sorted order for deterministic output.
func (e *Exporter) Write(w io.Writer, env map[string]string) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		if e.quote {
			v = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
		}
		var line string
		switch e.format {
		case FormatExport:
			line = fmt.Sprintf("export %s=%s\n", k, v)
		case FormatDocker:
			line = fmt.Sprintf("%s=%s\n", k, v)
		default:
			line = fmt.Sprintf("%s=%s\n", k, v)
		}
		if _, err := fmt.Fprint(w, line); err != nil {
			return fmt.Errorf("exporter: write key %q: %w", k, err)
		}
	}
	return nil
}

// WriteFile writes the env map to the file at path, creating or truncating it.
func (e *Exporter) WriteFile(path string, env map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("exporter: create file %q: %w", path, err)
	}
	defer f.Close()
	return e.Write(f, env)
}
