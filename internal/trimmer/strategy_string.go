package trimmer

import (
	"fmt"
	"strings"
)

var strategyNames = map[Strategy]string{
	TrimAll:      "all",
	TrimLeading:  "leading",
	TrimTrailing: "trailing",
	TrimInner:    "inner",
}

var strategyValues = map[string]Strategy{
	"all":      TrimAll,
	"leading":  TrimLeading,
	"trailing": TrimTrailing,
	"inner":    TrimInner,
}

// String returns the lowercase name of the strategy.
func (s Strategy) String() string {
	if name, ok := strategyNames[s]; ok {
		return name
	}
	return fmt.Sprintf("Strategy(%d)", int(s))
}

// ParseStrategy converts a string to a Strategy.
// It is case-insensitive. Returns an error for unknown values.
func ParseStrategy(s string) (Strategy, error) {
	norm := strings.ToLower(strings.TrimSpace(s))
	if st, ok := strategyValues[norm]; ok {
		return st, nil
	}
	return TrimAll, fmt.Errorf("trimmer: unknown strategy %q", s)
}
