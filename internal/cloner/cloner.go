package cloner

import "fmt"

// Strategy controls how key conflicts are handled during cloning.
type Strategy int

const (
	// SkipExisting keeps the destination value when a key already exists.
	SkipExisting Strategy = iota
	// OverwriteExisting replaces the destination value with the source value.
	OverwriteExisting
	// ErrorOnConflict returns an error if a key exists in both source and destination.
	ErrorOnConflict
)

// Clone copies keys from src into a new map, optionally merging with dst
// according to the given strategy. Neither src nor dst is mutated.
func Clone(src map[string]string, dst map[string]string, strategy Strategy) (map[string]string, error) {
	if src == nil {
		return nil, fmt.Errorf("cloner: src must not be nil")
	}

	result := make(map[string]string, len(dst)+len(src))

	// Seed result with destination values.
	for k, v := range dst {
		result[k] = v
	}

	for k, v := range src {
		existing, exists := result[k]
		switch {
		case !exists:
			result[k] = v
		case strategy == SkipExisting:
			// keep existing value — do nothing
		case strategy == OverwriteExisting:
			result[k] = v
		case strategy == ErrorOnConflict:
			return nil, fmt.Errorf("cloner: conflict on key %q (src=%q, dst=%q)", k, v, existing)
		}
	}

	return result, nil
}

// MustClone is like Clone but panics on error.
func MustClone(src, dst map[string]string, strategy Strategy) map[string]string {
	out, err := Clone(src, dst, strategy)
	if err != nil {
		panic(err)
	}
	return out
}
