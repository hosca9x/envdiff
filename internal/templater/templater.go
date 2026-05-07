// Package templater provides functionality for expanding variable references
// within .env file values, resolving placeholders like ${VAR} or $VAR using
// a provided environment map.
package templater

import (
	"fmt"
	"regexp"
	"strings"
)

// ErrMissingVar is returned when a referenced variable is not found and
// strict mode is enabled.
type ErrMissingVar struct {
	Key string
	Ref string
}

func (e *ErrMissingVar) Error() string {
	return fmt.Sprintf("key %q references undefined variable %q", e.Key, e.Ref)
}

var placeholderRe = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls expansion behaviour.
type Options struct {
	// Strict causes Expand to return an error when a referenced variable
	// is not present in the environment map.
	Strict bool
	// Fallback is returned for missing references when Strict is false.
	Fallback string
}

// Expand replaces variable references in every value of env using the
// variables defined in env itself (self-referential expansion).
// It returns a new map and leaves the original untouched.
func Expand(env map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		expanded, err := expandValue(k, v, env, opts)
		if err != nil {
			return nil, err
		}
		out[k] = expanded
	}
	return out, nil
}

// MustExpand is like Expand but panics on error.
func MustExpand(env map[string]string, opts Options) map[string]string {
	out, err := Expand(env, opts)
	if err != nil {
		panic(err)
	}
	return out
}

func expandValue(key, value string, env map[string]string, opts Options) (string, error) {
	var expandErr error
	result := placeholderRe.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		ref := strings.TrimPrefix(strings.Trim(match, "${}"), "$")
		if resolved, ok := env[ref]; ok {
			return resolved
		}
		if opts.Strict {
			expandErr = &ErrMissingVar{Key: key, Ref: ref}
			return match
		}
		return opts.Fallback
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}
