package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single key-value pair from an .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
	LineNum int
}

// EnvMap is an ordered collection of env entries keyed by variable name.
type EnvMap struct {
	Entries []Entry
	Index   map[string]int // key -> position in Entries
}

// Parse reads an .env file from r and returns an EnvMap.
func Parse(r io.Reader) (*EnvMap, error) {
	em := &EnvMap{
		Index: make(map[string]int),
	}

	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines
		if line == "" {
			continue
		}

		// Comment-only lines
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Inline comment stripping
		comment := ""
		if idx := strings.Index(line, " #"); idx != -1 {
			comment = strings.TrimSpace(line[idx+2:])
			line = strings.TrimSpace(line[:idx])
		}

		eqIdx := strings.IndexByte(line, '=')
		if eqIdx < 1 {
			return nil, fmt.Errorf("line %d: invalid format %q", lineNum, line)
		}

		key := strings.TrimSpace(line[:eqIdx])
		value := strings.TrimSpace(line[eqIdx+1:])
		value = stripQuotes(value)

		if _, exists := em.Index[key]; exists {
			return nil, fmt.Errorf("line %d: duplicate key %q", lineNum, key)
		}

		em.Index[key] = len(em.Entries)
		em.Entries = append(em.Entries, Entry{
			Key:     key,
			Value:   value,
			Comment: comment,
			LineNum: lineNum,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return em, nil
}

// Get returns the Entry for the given key and whether it was found.
func (em *EnvMap) Get(key string) (Entry, bool) {
	if idx, ok := em.Index[key]; ok {
		return em.Entries[idx], true
	}
	return Entry{}, false
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
