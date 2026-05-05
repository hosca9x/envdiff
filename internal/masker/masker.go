// Package masker provides utilities for masking sensitive values
// in environment variable maps before display or output.
package masker

import "strings"

const defaultMask = "***"

// DefaultSensitiveKeys contains common patterns for sensitive environment variable names.
var DefaultSensitiveKeys = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"CREDENTIAL",
	"AUTH",
	"ACCESS_KEY",
	"DATABASE_URL",
	"DSN",
}

// Masker masks sensitive values in environment variable maps.
type Masker struct {
	sensitivePatterns []string
	maskValue         string
}

// New creates a Masker with the default sensitive key patterns and mask value.
func New() *Masker {
	return &Masker{
		sensitivePatterns: DefaultSensitiveKeys,
		maskValue:         defaultMask,
	}
}

// NewWithOptions creates a Masker with custom patterns and mask value.
func NewWithOptions(patterns []string, maskValue string) *Masker {
	return &Masker{
		sensitivePatterns: patterns,
		maskValue:         maskValue,
	}
}

// IsSensitive returns true if the given key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range m.sensitivePatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// MaskMap returns a copy of the provided map with sensitive values replaced by the mask.
func (m *Masker) MaskMap(env map[string]string) map[string]string {
	masked := make(map[string]string, len(env))
	for k, v := range env {
		if m.IsSensitive(k) {
			masked[k] = m.maskValue
		} else {
			masked[k] = v
		}
	}
	return masked
}

// MaskValue returns the masked value if the key is sensitive, otherwise returns the original value.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return m.maskValue
	}
	return value
}
