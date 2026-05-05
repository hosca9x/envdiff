// Package loader provides functionality for loading .env files from
// the filesystem and preparing them for diffing or reconciliation.
package loader

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/parser"
)

// LoadFile reads a .env file from the given path and returns a map of
// key-value pairs. It returns an error if the file cannot be read or parsed.
func LoadFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loader: reading file %q: %w", path, err)
	}
	return parser.Parse(string(data))
}

// LoadFiles loads multiple .env files and returns a slice of maps in the
// same order as the provided paths. If any file fails to load, the error
// is returned immediately and no further files are processed.
func LoadFiles(paths []string) ([]map[string]string, error) {
	results := make([]map[string]string, 0, len(paths))
	for _, p := range paths {
		m, err := LoadFile(p)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

// MustLoadFile is like LoadFile but panics on error. Useful in tests or
// init functions where error handling is not required.
func MustLoadFile(path string) map[string]string {
	m, err := LoadFile(path)
	if err != nil {
		panic(err)
	}
	return m
}
