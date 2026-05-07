// Package filter provides utilities for filtering env map entries
// based on key patterns, prefixes, or custom predicates.
package filter

import "strings"

// Option configures how filtering is applied.
type Option func(*options)

type options struct {
	prefixes   []string
	keys       []string
	exclude    bool
}

// WithPrefixes returns only keys that start with any of the given prefixes.
func WithPrefixes(prefixes ...string) Option {
	return func(o *options) {
		o.prefixes = append(o.prefixes, prefixes...)
	}
}

// WithKeys returns only the explicitly listed keys.
func WithKeys(keys ...string) Option {
	return func(o *options) {
		o.keys = append(o.keys, keys...)
	}
}

// Exclude inverts the filter, removing matched keys instead of keeping them.
func Exclude() Option {
	return func(o *options) {
		o.exclude = true
	}
}

// Filter returns a new map containing only the entries that match the
// provided options. If no options are given, the original map is returned
// as-is (shallow copy).
func Filter(env map[string]string, opts ...Option) map[string]string {
	cfg := &options{}
	for _, o := range opts {
		o(cfg)
	}

	keySet := make(map[string]struct{}, len(cfg.keys))
	for _, k := range cfg.keys {
		keySet[k] = struct{}{}
	}

	result := make(map[string]string)
	for k, v := range env {
		matched := matchesKey(k, keySet, cfg.prefixes)
		if cfg.exclude {
			matched = !matched
		}
		if matched || (len(cfg.keys) == 0 && len(cfg.prefixes) == 0) {
			result[k] = v
		}
	}
	return result
}

func matchesKey(key string, keySet map[string]struct{}, prefixes []string) bool {
	if _, ok := keySet[key]; ok {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
