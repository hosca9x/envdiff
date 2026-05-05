// Package differ provides functionality to compare two parsed .env files
// and produce a structured diff result highlighting added, removed, and changed keys.
package differ

// Status represents the diff status of a key.
type Status string

const (
	StatusAdded   Status = "added"
	StatusRemoved Status = "removed"
	StatusChanged Status = "changed"
	StatusSame    Status = "same"
)

// Entry represents a single diff entry for a key across two environments.
type Entry struct {
	Key      string
	Status   Status
	ValueA   string // value in the "base" env file
	ValueB   string // value in the "target" env file
}

// Result holds the full diff between two env maps.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if any entry is not StatusSame.
func (r *Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Status != StatusSame {
			return true
		}
	}
	return false
}

// Diff compares two env maps (key -> value) and returns a Result.
// mapA is considered the "base" and mapB the "target".
func Diff(mapA, mapB map[string]string) Result {
	seen := make(map[string]bool)
	var entries []Entry

	for k, vA := range mapA {
		seen[k] = true
		if vB, ok := mapB[k]; ok {
			if vA == vB {
				entries = append(entries, Entry{Key: k, Status: StatusSame, ValueA: vA, ValueB: vB})
			} else {
				entries = append(entries, Entry{Key: k, Status: StatusChanged, ValueA: vA, ValueB: vB})
			}
		} else {
			entries = append(entries, Entry{Key: k, Status: StatusRemoved, ValueA: vA, ValueB: ""})
		}
	}

	for k, vB := range mapB {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Status: StatusAdded, ValueA: "", ValueB: vB})
		}
	}

	return Result{Entries: entries}
}
