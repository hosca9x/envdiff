package sorter

import (
	"fmt"
	"strings"
)

// String returns the lowercase string representation of a Strategy.
func (s Strategy) String() string {
	switch s {
	case Alphabetical:
		return "alphabetical"
	case Reverse:
		return "reverse"
	case GroupedByPrefix:
		return "grouped"
	default:
		return fmt.Sprintf("strategy(%d)", int(s))
	}
}

// ParseStrategy converts a string to a Strategy value.
// Accepted values (case-insensitive): "alphabetical", "reverse", "grouped".
func ParseStrategy(s string) (Strategy, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "alphabetical", "alpha", "asc":
		return Alphabetical, nil
	case "reverse", "desc":
		return Reverse, nil
	case "grouped", "prefix":
		return GroupedByPrefix, nil
	default:
		return Alphabetical, fmt.Errorf("unknown sort strategy %q: must be one of alphabetical, reverse, grouped", s)
	}
}
