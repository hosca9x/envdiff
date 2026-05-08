// Package differ provides utilities for computing the difference between
// two environment variable maps.
//
// # Overview
//
// The primary entry point is the [Diff] function, which accepts a source
// (new state) and a target (baseline) map and returns a sorted slice of
// [Entry] values describing each key's change status.
//
// # Change Types
//
// Each entry is annotated with one of four [ChangeType] values:
//
//   - [Added]     – key exists in source but not in target
//   - [Removed]   – key exists in target but not in source
//   - [Changed]   – key exists in both but values differ
//   - [Unchanged] – key exists in both with identical values
//
// # Example
//
//	entries := differ.Diff(source, target)
//	for _, e := range entries {
//		fmt.Println(e.Key, e.ChangeType)
//	}
package differ
