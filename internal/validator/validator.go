// Package validator provides validation utilities for .env file entries,
// checking for common issues such as missing values, invalid key formats,
// and suspicious patterns that may indicate misconfiguration.
package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Issue represents a single validation problem found in an env map.
type Issue struct {
	Key      string
	Severity string // "error" or "warning"
	Message  string
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(i.Severity), i.Key, i.Message)
}

var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Validate checks the provided env map for common issues and returns a
// slice of Issues. An empty slice means no problems were found.
func Validate(env map[string]string) []Issue {
	var issues []Issue

	for key, value := range env {
		// Check key format
		if !validKeyPattern.MatchString(key) {
			issues = append(issues, Issue{
				Key:      key,
				Severity: "error",
				Message:  "key contains invalid characters; must match [A-Za-z_][A-Za-z0-9_]*",
			})
		}

		// Check for empty values
		if strings.TrimSpace(value) == "" {
			issues = append(issues, Issue{
				Key:      key,
				Severity: "warning",
				Message:  "value is empty",
			})
		}

		// Check for unresolved placeholder patterns like ${VAR} or %VAR%
		if strings.Contains(value, "${") || strings.Contains(value, "%") {
			issues = append(issues, Issue{
				Key:      key,
				Severity: "warning",
				Message:  "value may contain an unresolved variable placeholder",
			})
		}
	}

	return issues
}

// HasErrors returns true if any of the provided issues have severity "error".
func HasErrors(issues []Issue) bool {
	for _, i := range issues {
		if i.Severity == "error" {
			return true
		}
	}
	return false
}
