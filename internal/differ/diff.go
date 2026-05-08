// Package differ computes the difference between two env maps.
package differ

import "sort"

// Diff computes the difference between a source and target env map.
// source is the "new" state; target is the "old" / baseline state.
// Returns a slice of Entry values sorted by key.
func Diff(source, target map[string]string) []Entry {
	keys := make(map[string]struct{})
	for k := range source {
		keys[k] = struct{}{}
	}
	for k := range target {
		keys[k] = struct{}{}
	}

	entries := make([]Entry, 0, len(keys))
	for k := range keys {
		srcVal, inSrc := source[k]
		tgtVal, inTgt := target[k]

		switch {
		case inSrc && !inTgt:
			entries = append(entries, Entry{
				Key:        k,
				ChangeType: Added,
				NewValue:   srcVal,
			})
		case !inSrc && inTgt:
			entries = append(entries, Entry{
				Key:        k,
				ChangeType: Removed,
				OldValue:   tgtVal,
			})
		case srcVal != tgtVal:
			entries = append(entries, Entry{
				Key:        k,
				ChangeType: Changed,
				OldValue:   tgtVal,
				NewValue:   srcVal,
			})
		default:
			entries = append(entries, Entry{
				Key:        k,
				ChangeType: Unchanged,
				OldValue:   tgtVal,
				NewValue:   srcVal,
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}
