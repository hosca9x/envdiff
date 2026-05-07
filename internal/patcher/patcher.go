// Package patcher applies a set of diff entries back to an env map,
// producing a patched result that can be written out as a new .env file.
package patcher

import (
	"fmt"

	"github.com/user/envdiff/internal/differ"
)

// Option controls how the patcher behaves when conflicts arise.
type Option func(*options)

type options struct {
	skipRemoved bool
	skipAdded   bool
}

// SkipRemoved instructs the patcher not to delete keys that are marked as removed.
func SkipRemoved() Option {
	return func(o *options) { o.skipRemoved = true }
}

// SkipAdded instructs the patcher not to insert keys that are marked as added.
func SkipAdded() Option {
	return func(o *options) { o.skipAdded = true }
}

// Apply takes a base env map and a slice of DiffEntry values and returns a new
// map with the patch applied. The original map is never mutated.
//
// An error is returned if a Changed or Removed entry references a key that does
// not exist in base, indicating the patch is out of sync with the target file.
func Apply(base map[string]string, entries []differ.DiffEntry, opts ...Option) (map[string]string, error) {
	cfg := &options{}
	for _, o := range opts {
		o(cfg)
	}

	// Shallow-copy base so we never mutate the caller's map.
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}

	for _, e := range entries {
		switch e.Type {
		case differ.Added:
			if !cfg.skipAdded {
				result[e.Key] = e.NewValue
			}
		case differ.Removed:
			if !cfg.skipRemoved {
				if _, ok := result[e.Key]; !ok {
					return nil, fmt.Errorf("patcher: cannot remove key %q: not present in base", e.Key)
				}
				delete(result, e.Key)
			}
		case differ.Changed:
			if _, ok := result[e.Key]; !ok {
				return nil, fmt.Errorf("patcher: cannot change key %q: not present in base", e.Key)
			}
			result[e.Key] = e.NewValue
		}
	}

	return result, nil
}
