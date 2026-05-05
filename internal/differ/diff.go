// Package differ computes differences between two env maps.
package differ

// ChangeType represents the kind of change detected.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// DiffEntry represents a single key-level difference between two env maps.
type DiffEntry struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Diff computes the differences between a base env map and a target env map.
// Keys present only in target are Added, keys only in base are Removed,
// and keys present in both with different values are Changed.
func Diff(base, target map[string]string) []DiffEntry {
	var entries []DiffEntry

	for k, tv := range target {
		if bv, ok := base[k]; !ok {
			entries = append(entries, DiffEntry{
				Key:      k,
				Type:     Added,
				NewValue: tv,
			})
		} else if bv != tv {
			entries = append(entries, DiffEntry{
				Key:      k,
				Type:     Changed,
				OldValue: bv,
				NewValue: tv,
			})
		}
	}

	for k, bv := range base {
		if _, ok := target[k]; !ok {
			entries = append(entries, DiffEntry{
				Key:      k,
				Type:     Removed,
				OldValue: bv,
			})
		}
	}

	return entries
}
