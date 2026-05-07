// Package patcher applies a set of diff entries produced by the differ package
// to a base env map, yielding a patched map ready for export.
//
// # Overview
//
// Given a base map and a []differ.DiffEntry, Apply produces a new map with all
// additions, removals, and changes reflected. The base map is never modified.
//
// # Options
//
// SkipRemoved — retain keys that the diff marks as removed.
// SkipAdded   — omit keys that the diff marks as added.
//
// # Error handling
//
// Apply returns an error if a Changed or Removed entry references a key that is
// absent from the base map, which signals that the patch is stale or was
// generated against a different file.
//
// # Typical usage
//
//	entries := differ.Diff(source, target)
//	patched, err := patcher.Apply(target, entries, patcher.SkipRemoved())
//	if err != nil {
//		log.Fatal(err)
//	}
package patcher
