// Package snapshot provides functionality for capturing and comparing
// environment variable states at different points in time.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a captured state of environment variables at a specific time.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Label     string            `json:"label"`
	Env       map[string]string `json:"env"`
}

// New creates a new Snapshot with the given label and environment map.
// The timestamp is set to the current UTC time.
func New(label string, env map[string]string) *Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Env:       copy,
	}
}

// Save writes the snapshot to a JSON file at the given path.
func Save(s *Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file %q: %w", path, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}

// Get returns the value of the environment variable with the given key and
// a boolean indicating whether the key was present in the snapshot.
func (s *Snapshot) Get(key string) (string, bool) {
	v, ok := s.Env[key]
	return v, ok
}
