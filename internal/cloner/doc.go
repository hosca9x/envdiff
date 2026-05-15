// Package cloner provides utilities for deep-copying .env maps and merging
// them into a destination map with configurable conflict resolution.
//
// # Strategies
//
// Three strategies control what happens when a key exists in both the source
// and the destination:
//
//   - SkipExisting   – the destination value is preserved.
//   - OverwriteExisting – the source value replaces the destination value.
//   - ErrorOnConflict   – Clone returns an error immediately.
//
// Neither the source nor the destination map is ever mutated; Clone always
// returns a freshly allocated map.
//
// # Usage
//
//	out, err := cloner.Clone(src, dst, cloner.SkipExisting)
//	out, err := cloner.Clone(src, dst, cloner.OverwriteExisting)
//	out  = cloner.MustClone(src, dst, cloner.SkipExisting) // panics on error
package cloner
