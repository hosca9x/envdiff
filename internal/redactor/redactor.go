// Package redactor provides functionality for redacting sensitive values
// from env maps before logging, display, or transmission.
package redactor

import (
	"strings"
)

const defaultPlaceholder = "[REDACTED]"

// Redactor replaces sensitive values in env maps with a placeholder.
type Redactor struct {
	placeholder string
	patterns    []string
}

// Option configures a Redactor.
type Option func(*Redactor)

// WithPlaceholder sets a custom placeholder string.
func WithPlaceholder(p string) Option {
	return func(r *Redactor) {
		if p != "" {
			r.placeholder = p
		}
	}
}

// WithPatterns adds additional key patterns to treat as sensitive.
func WithPatterns(patterns ...string) Option {
	return func(r *Redactor) {
		r.patterns = append(r.patterns, patterns...)
	}
}

// defaultPatterns are substrings that indicate a key holds sensitive data.
var defaultPatterns = []string{
	"secret", "password", "passwd", "token", "apikey", "api_key",
	"auth", "credential", "private", "cert", "key",
}

// New creates a Redactor with default settings.
func New(opts ...Option) *Redactor {
	r := &Redactor{
		placeholder: defaultPlaceholder,
		patterns:    append([]string{}, defaultPatterns...),
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// IsSensitive reports whether the given key should be redacted.
func (r *Redactor) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// Redact returns a new map with sensitive values replaced by the placeholder.
// The original map is never mutated.
func (r *Redactor) Redact(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.IsSensitive(k) {
			out[k] = r.placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactValue returns the placeholder if the key is sensitive, otherwise the
// original value.
func (r *Redactor) RedactValue(key, value string) string {
	if r.IsSensitive(key) {
		return r.placeholder
	}
	return value
}
