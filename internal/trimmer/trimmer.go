package trimmer

import (
	"fmt"
	"strings"
)

// Strategy controls how trimming is applied to env map values.
type Strategy int

const (
	// TrimAll removes leading and trailing whitespace from all values.
	TrimAll Strategy = iota
	// TrimLeading removes only leading whitespace.
	TrimLeading
	// TrimTrailing removes only trailing whitespace.
	TrimTrailing
	// TrimInner collapses internal runs of whitespace to a single space.
	TrimInner
)

// Result holds the outcome of a trim operation for a single key.
type Result struct {
	Key      string
	Original string
	Trimmed  string
	Changed  bool
}

// Trim applies the given strategy to every value in env and returns a new
// map along with a slice of Result describing what changed.
func Trim(env map[string]string, strategy Strategy) (map[string]string, []Result, error) {
	if env == nil {
		return nil, nil, fmt.Errorf("trimmer: env map must not be nil")
	}

	out := make(map[string]string, len(env))
	var results []Result

	for k, v := range env {
		trimmed := applyStrategy(v, strategy)
		out[k] = trimmed
		if trimmed != v {
			results = append(results, Result{
				Key:      k,
				Original: v,
				Trimmed:  trimmed,
				Changed:  true,
			})
		}
	}

	return out, results, nil
}

// MustTrim is like Trim but panics on error.
func MustTrim(env map[string]string, strategy Strategy) (map[string]string, []Result) {
	out, results, err := Trim(env, strategy)
	if err != nil {
		panic(err)
	}
	return out, results
}

func applyStrategy(v string, s Strategy) string {
	switch s {
	case TrimLeading:
		return strings.TrimLeft(v, " \t")
	case TrimTrailing:
		return strings.TrimRight(v, " \t")
	case TrimInner:
		fields := strings.Fields(v)
		return strings.Join(fields, " ")
	default: // TrimAll
		return strings.TrimSpace(v)
	}
}
