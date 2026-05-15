package grouper

import (
	"fmt"
	"strings"
)

// String returns the canonical string representation of a Strategy.
func (s Strategy) String() string {
	switch s {
	case ByPrefix:
		return "prefix"
	case ByLength:
		return "length"
	default:
		return fmt.Sprintf("Strategy(%d)", int(s))
	}
}

// ParseStrategy converts a string to a Strategy value.
// Accepted values (case-insensitive): "prefix", "length".
func ParseStrategy(s string) (Strategy, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "prefix":
		return ByPrefix, nil
	case "length":
		return ByLength, nil
	default:
		return 0, fmt.Errorf("grouper: unknown strategy %q (want \"prefix\" or \"length\")", s)
	}
}
