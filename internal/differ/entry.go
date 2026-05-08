package differ

// ChangeType represents the kind of difference between two env maps.
type ChangeType string

const (
	// Added indicates a key present in source but not in target.
	Added ChangeType = "added"
	// Removed indicates a key present in target but not in source.
	Removed ChangeType = "removed"
	// Changed indicates a key present in both but with different values.
	Changed ChangeType = "changed"
	// Unchanged indicates a key present in both with the same value.
	Unchanged ChangeType = "unchanged"
)

// Entry represents a single diff result for one key.
type Entry struct {
	Key        string
	ChangeType ChangeType
	// OldValue is the value from the target (baseline) map.
	OldValue string
	// NewValue is the value from the source map.
	NewValue string
}

// IsChanged returns true if the entry represents any kind of change.
func (e Entry) IsChanged() bool {
	return e.ChangeType != Unchanged
}
