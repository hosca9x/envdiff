// Package pinner provides utilities for locking (pinning) specific keys in an
// env map to predetermined values.
//
// Pinning is useful when deploying across multiple environments and certain
// keys — such as APP_ENV, LOG_LEVEL, or FEATURE_FLAG_X — must always carry a
// canonical value regardless of what the source .env file contains.
//
// # Basic usage
//
//	out, err := pinner.Pin(env, pins, pinner.OverwriteExisting)
//
// # Strategies
//
// Three conflict-resolution strategies are available:
//
//   - SkipExisting   – keep the value already present in env.
//   - OverwriteExisting – replace env's value with the pinned value.
//   - ErrorOnConflict – return an error if the key exists with a different value.
//
// The original map is never mutated; Pin always returns a new map.
package pinner
