package sorter

import (
	"sort"
	"strings"
)

// Strategy defines how keys should be sorted.
type Strategy int

const (
	// Alphabetical sorts keys in ascending lexicographic order.
	Alphabetical Strategy = iota
	// Reverse sorts keys in descending lexicographic order.
	Reverse
	// GroupedByPrefix sorts keys grouped by their prefix (e.g. DB_, APP_).
	GroupedByPrefix
)

// Sort returns a new map with keys sorted according to the given strategy.
// The returned slice of keys reflects the desired ordering.
func Sort(env map[string]string, strategy Strategy) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	switch strategy {
	case Reverse:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	case GroupedByPrefix:
		keys = sortByPrefix(keys)
	default:
		sort.Strings(keys)
	}

	return keys
}

// sortByPrefix groups keys by their underscore-delimited prefix, then sorts
// within each group alphabetically. Keys without a prefix are placed last.
func sortByPrefix(keys []string) []string {
	groups := make(map[string][]string)
	order := []string{}
	seen := map[string]bool{}

	for _, k := range keys {
		prefix := extractPrefix(k)
		if !seen[prefix] {
			seen[prefix] = true
			order = append(order, prefix)
		}
		groups[prefix] = append(groups[prefix], k)
	}

	sort.Strings(order)

	result := make([]string, 0, len(keys))
	for _, prefix := range order {
		grp := groups[prefix]
		sort.Strings(grp)
		result = append(result, grp...)
	}
	return result
}

// extractPrefix returns the portion of a key before the first underscore.
// If no underscore exists, the whole key is treated as its own prefix.
func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
