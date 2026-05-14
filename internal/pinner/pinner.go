package pinner

import (
	"errors"
	"fmt"
	"sort"
)

// Strategy controls how pinning behaves when a key already has a pinned value.
type Strategy int

const (
	// SkipExisting keeps the already-pinned value unchanged.
	SkipExisting Strategy = iota
	// OverwriteExisting replaces the pinned value with the new one.
	OverwriteExisting
	// ErrorOnConflict returns an error if a key is pinned more than once.
	ErrorOnConflict
)

// ErrConflict is returned by Pin when a key is already pinned and the
// strategy is ErrorOnConflict.
var ErrConflict = errors.New("pinner: key already pinned")

// Pin locks specific keys in env to the values provided in pins.
// Keys present in env but absent from pins are left untouched.
// Keys present in pins but absent from env are added.
// The original map is never mutated; a new map is returned.
func Pin(env map[string]string, pins map[string]string, strategy Strategy) (map[string]string, error) {
	if env == nil {
		return nil, errors.New("pinner: env must not be nil")
	}
	if pins == nil {
		return nil, errors.New("pinner: pins must not be nil")
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for k, pinVal := range pins {
		existing, exists := out[k]
		if !exists {
			out[k] = pinVal
			continue
		}
		if existing == pinVal {
			continue
		}
		switch strategy {
		case SkipExisting:
			// keep existing value
		case OverwriteExisting:
			out[k] = pinVal
		case ErrorOnConflict:
			return nil, fmt.Errorf("%w: %q (current=%q, pin=%q)", ErrConflict, k, existing, pinVal)
		}
	}
	return out, nil
}

// MustPin is like Pin but panics on error.
func MustPin(env map[string]string, pins map[string]string, strategy Strategy) map[string]string {
	out, err := Pin(env, pins, strategy)
	if err != nil {
		panic(err)
	}
	return out
}

// PinnedKeys returns a sorted slice of keys that differ between env and pins
// (i.e. keys that would be overwritten by an OverwriteExisting pin operation).
func PinnedKeys(env map[string]string, pins map[string]string) []string {
	var conflicting []string
	for k, pinVal := range pins {
		if existing, ok := env[k]; ok && existing != pinVal {
			conflicting = append(conflicting, k)
		}
	}
	sort.Strings(conflicting)
	return conflicting
}
