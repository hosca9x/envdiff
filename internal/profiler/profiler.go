package profiler

import (
	"sort"
	"strings"
)

// Profile holds statistical analysis of an env map.
type Profile struct {
	TotalKeys       int
	EmptyValues     int
	SensitiveKeys   int
	PrefixGroups    map[string]int
	AvgValueLength  float64
	LongestKey      string
	ShortestKey     string
	DuplicatePrefixes []string
}

var sensitivePatterns = []string{
	"secret", "password", "passwd", "token", "api_key",
	"auth", "private", "credential", "cert", "key",
}

// Analyze computes a Profile from the given env map.
func Analyze(env map[string]string) Profile {
	if len(env) == 0 {
		return Profile{PrefixGroups: map[string]int{}}
	}

	p := Profile{
		PrefixGroups: make(map[string]int),
	}

	totalLen := 0
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	p.LongestKey = keys[0]
	p.ShortestKey = keys[0]

	for _, k := range keys {
		v := env[k]
		p.TotalKeys++
		totalLen += len(v)

		if v == "" {
			p.EmptyValues++
		}
		if isSensitive(k) {
			p.SensitiveKeys++
		}
		if len(k) > len(p.LongestKey) {
			p.LongestKey = k
		}
		if len(k) < len(p.ShortestKey) {
			p.ShortestKey = k
		}

		prefix := extractPrefix(k)
		p.PrefixGroups[prefix]++
	}

	if p.TotalKeys > 0 {
		p.AvgValueLength = float64(totalLen) / float64(p.TotalKeys)
	}

	for prefix, count := range p.PrefixGroups {
		if count > 1 {
			p.DuplicatePrefixes = append(p.DuplicatePrefixes, prefix)
		}
	}
	sort.Strings(p.DuplicatePrefixes)

	return p
}

func isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, pat := range sensitivePatterns {
		if strings.Contains(lower, pat) {
			return true
		}
	}
	return false
}

func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
