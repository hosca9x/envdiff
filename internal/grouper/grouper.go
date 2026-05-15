package grouper

import (
	"fmt"
	"sort"
	"strings"
)

// Strategy defines how keys are grouped.
type Strategy int

const (
	// ByPrefix groups keys by their underscore-delimited prefix (e.g. DB_HOST → DB).
	ByPrefix Strategy = iota
	// ByLength groups keys into buckets: short (<= 8), medium (<= 20), long (> 20).
	ByLength
)

// Group holds a named collection of key-value pairs.
type Group struct {
	Name string
	Keys map[string]string
}

// GroupBy partitions env into named groups according to the chosen strategy.
// Keys that do not match any natural group are placed in "_other".
func GroupBy(env map[string]string, s Strategy) ([]Group, error) {
	if env == nil {
		return nil, fmt.Errorf("grouper: env map must not be nil")
	}

	buckets := map[string]map[string]string{}

	for k, v := range env {
		name := bucketName(k, s)
		if _, ok := buckets[name]; !ok {
			buckets[name] = map[string]string{}
		}
		buckets[name][k] = v
	}

	names := make([]string, 0, len(buckets))
	for n := range buckets {
		names = append(names, n)
	}
	sort.Strings(names)

	groups := make([]Group, 0, len(names))
	for _, n := range names {
		groups = append(groups, Group{Name: n, Keys: buckets[n]})
	}
	return groups, nil
}

// MustGroupBy is like GroupBy but panics on error.
func MustGroupBy(env map[string]string, s Strategy) []Group {
	g, err := GroupBy(env, s)
	if err != nil {
		panic(err)
	}
	return g
}

func bucketName(key string, s Strategy) string {
	switch s {
	case ByPrefix:
		if idx := strings.Index(key, "_"); idx > 0 {
			return key[:idx]
		}
		return "_other"
	case ByLength:
		switch {
		case len(key) <= 8:
			return "short"
		case len(key) <= 20:
			return "medium"
		default:
			return "long"
		}
	}
	return "_other"
}
