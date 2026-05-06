package merger

import "fmt"

// String returns a human-readable name for a Strategy constant.
func (s Strategy) String() string {
	switch s {
	case FirstWins:
		return "FirstWins"
	case LastWins:
		return "LastWins"
	case ErrorOnConflict:
		return "ErrorOnConflict"
	default:
		return fmt.Sprintf("Strategy(%d)", int(s))
	}
}

// ParseStrategy converts a string such as "first", "last", or "error" into
// the corresponding Strategy constant. It is case-insensitive.
func ParseStrategy(s string) (Strategy, error) {
	switch s {
	case "first", "firstwins", "FirstWins":
		return FirstWins, nil
	case "last", "lastwins", "LastWins":
		return LastWins, nil
	case "error", "erroronconflict", "ErrorOnConflict":
		return ErrorOnConflict, nil
	default:
		return FirstWins, fmt.Errorf("merger: unknown strategy %q (want first|last|error)", s)
	}
}
