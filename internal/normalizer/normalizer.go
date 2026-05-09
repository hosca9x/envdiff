package normalizer

import (
	"errors"
	"strings"
)

// Strategy defines how keys should be normalized.
type Strategy int

const (
	// UpperSnake converts keys to UPPER_SNAKE_CASE (e.g. myKey -> MY_KEY).
	UpperSnake Strategy = iota
	// LowerSnake converts keys to lower_snake_case (e.g. MyKey -> my_key).
	LowerSnake
	// TrimPrefix strips a given prefix from all keys.
	TrimPrefix
)

// Options controls normalization behaviour.
type Options struct {
	Strategy Strategy
	// Prefix is used when Strategy == TrimPrefix.
	Prefix string
}

// Normalize applies the chosen strategy to every key in env and returns a new
// map. The original map is never mutated. An error is returned when env is nil.
func Normalize(env map[string]string, opts Options) (map[string]string, error) {
	if env == nil {
		return nil, errors.New("normalizer: env map must not be nil")
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[applyStrategy(k, opts)] = v
	}
	return out, nil
}

// MustNormalize is like Normalize but panics on error.
func MustNormalize(env map[string]string, opts Options) map[string]string {
	result, err := Normalize(env, opts)
	if err != nil {
		panic(err)
	}
	return result
}

func applyStrategy(key string, opts Options) string {
	switch opts.Strategy {
	case UpperSnake:
		return toUpperSnake(key)
	case LowerSnake:
		return strings.ToLower(toUpperSnake(key))
	case TrimPrefix:
		return strings.TrimPrefix(key, opts.Prefix)
	default:
		return key
	}
}

// toUpperSnake converts camelCase / PascalCase / kebab-case to UPPER_SNAKE_CASE.
func toUpperSnake(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, ".", "_")
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' && i > 0 && s[i-1] != '_' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToUpper(b.String())
}
