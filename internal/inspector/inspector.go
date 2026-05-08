// Package inspector provides utilities for inspecting and summarizing
// the structure and statistics of an env map, such as key counts,
// value type hints, and prefix groupings.
package inspector

import (
	"sort"
	"strings"
)

// Report holds the inspection results for an env map.
type Report struct {
	TotalKeys    int
	EmptyValues  int
	PrefixGroups map[string][]string
	ValueHints   map[string]string
}

// Inspect analyzes the given env map and returns a Report.
func Inspect(env map[string]string) Report {
	r := Report{
		TotalKeys:    len(env),
		PrefixGroups: make(map[string][]string),
		ValueHints:   make(map[string]string),
	}

	for k, v := range env {
		if v == "" {
			r.EmptyValues++
		}

		prefix := extractPrefix(k)
		r.PrefixGroups[prefix] = append(r.PrefixGroups[prefix], k)
		r.ValueHints[k] = hintFor(v)
	}

	for prefix := range r.PrefixGroups {
		sort.Strings(r.PrefixGroups[prefix])
	}

	return r
}

// extractPrefix returns the portion of a key before the first underscore,
// or the full key if no underscore is present.
func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}

// hintFor returns a short string describing the likely type of a value.
func hintFor(v string) string {
	if v == "" {
		return "empty"
	}
	v = strings.TrimSpace(v)
	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		return "url"
	}
	if strings.EqualFold(v, "true") || strings.EqualFold(v, "false") {
		return "bool"
	}
	allDigits := true
	for _, c := range v {
		if c < '0' || c > '9' {
			allDigits = false
			break
		}
	}
	if allDigits {
		return "number"
	}
	return "string"
}
