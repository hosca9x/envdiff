package transformer

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a single env value.
type TransformFunc func(key, value string) (string, error)

// Transformer applies a chain of TransformFuncs to an env map.
type Transformer struct {
	fns []TransformFunc
}

// New returns a Transformer with the given transform functions applied in order.
func New(fns ...TransformFunc) *Transformer {
	return &Transformer{fns: fns}
}

// Apply runs all transform functions over the input map and returns a new map.
// If any transform returns an error, Apply stops and returns that error.
func (t *Transformer) Apply(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for _, fn := range t.fns {
		for k, v := range out {
			newVal, err := fn(k, v)
			if err != nil {
				return nil, fmt.Errorf("transformer: key %q: %w", k, err)
			}
			out[k] = newVal
		}
	}
	return out, nil
}

// TrimSpace returns a TransformFunc that trims leading/trailing whitespace from values.
func TrimSpace() TransformFunc {
	return func(_, value string) (string, error) {
		return strings.TrimSpace(value), nil
	}
}

// ToUpper returns a TransformFunc that uppercases all values.
func ToUpper() TransformFunc {
	return func(_, value string) (string, error) {
		return strings.ToUpper(value), nil
	}
}

// ToLower returns a TransformFunc that lowercases all values.
func ToLower() TransformFunc {
	return func(_, value string) (string, error) {
		return strings.ToLower(value), nil
	}
}

// ReplaceValue returns a TransformFunc that replaces old with new in every value.
func ReplaceValue(old, replacement string) TransformFunc {
	return func(_, value string) (string, error) {
		return strings.ReplaceAll(value, old, replacement), nil
	}
}

// MustApply is like Apply but panics on error.
func (t *Transformer) MustApply(env map[string]string) map[string]string {
	out, err := t.Apply(env)
	if err != nil {
		panic(err)
	}
	return out
}
