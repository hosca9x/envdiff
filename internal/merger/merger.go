// Package merger provides functionality to merge multiple .env files
// into a single unified map, with configurable conflict resolution strategies.
package merger

import "fmt"

// Strategy defines how conflicts are resolved when the same key appears
// in multiple source files.
type Strategy int

const (
	// FirstWins keeps the value from the first file that defined the key.
	FirstWins Strategy = iota
	// LastWins keeps the value from the last file that defined the key.
	LastWins
	// ErrorOnConflict returns an error if the same key appears with different values.
	ErrorOnConflict
)

// ConflictError describes a key that appeared in multiple files with differing values.
type ConflictError struct {
	Key    string
	First  string
	Second string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("merger: conflict on key %q: %q vs %q", e.Key, e.First, e.Second)
}

// Merge combines multiple env maps into one according to the given Strategy.
// The maps are processed in the order they are provided.
func Merge(strategy Strategy, maps ...map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			existing, exists := result[k]
			switch {
			case !exists:
				result[k] = v
			case strategy == FirstWins:
				// keep existing, do nothing
			case strategy == LastWins:
				result[k] = v
			case strategy == ErrorOnConflict:
				if existing != v {
					return nil, &ConflictError{Key: k, First: existing, Second: v}
				}
			}
		}
	}

	return result, nil
}

// MustMerge is like Merge but panics on error.
func MustMerge(strategy Strategy, maps ...map[string]string) map[string]string {
	result, err := Merge(strategy, maps...)
	if err != nil {
		panic(err)
	}
	return result
}
