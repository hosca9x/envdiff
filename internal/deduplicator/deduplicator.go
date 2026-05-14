package deduplicator

import "fmt"

// Strategy controls how duplicate keys across multiple env maps are resolved.
type Strategy int

const (
	// FirstWins keeps the value from the first map that defines a key.
	FirstWins Strategy = iota
	// LastWins keeps the value from the last map that defines a key.
	LastWins
	// ErrorOnDuplicate returns an error if any key appears in more than one map.
	ErrorOnDuplicate
)

// Result holds the deduplicated map and a record of which keys were duplicated.
type Result struct {
	Env        map[string]string
	Duplicates map[string][]string // key -> list of values seen
}

// Deduplicate merges multiple env maps, resolving duplicate keys using the
// given strategy. It never mutates any of the input maps.
func Deduplicate(strategy Strategy, maps ...map[string]string) (Result, error) {
	out := make(map[string]string)
	duplicates := make(map[string][]string)
	seen := make(map[string]bool)

	for _, m := range maps {
		for k, v := range m {
			if !seen[k] {
				out[k] = v
				seen[k] = true
				continue
			}
			// Key already present — record the duplicate value.
			if _, tracked := duplicates[k]; !tracked {
				duplicates[k] = []string{out[k]}
			}
			duplicates[k] = append(duplicates[k], v)

			switch strategy {
			case FirstWins:
				// keep existing value — do nothing
			case LastWins:
				out[k] = v
			case ErrorOnDuplicate:
				return Result{}, fmt.Errorf("deduplicator: duplicate key %q", k)
			}
		}
	}

	return Result{Env: out, Duplicates: duplicates}, nil
}

// MustDeduplicate is like Deduplicate but panics on error.
func MustDeduplicate(strategy Strategy, maps ...map[string]string) Result {
	r, err := Deduplicate(strategy, maps...)
	if err != nil {
		panic(err)
	}
	return r
}
