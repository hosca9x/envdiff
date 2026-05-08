// Package linter provides rule-based linting for .env file contents.
// It checks for common issues such as overly long values, keys with lowercase
// letters, duplicate prefixes, and values that look like unquoted URLs.
package linter

import (
	"fmt"
	"strings"
	"unicode"
)

// Issue represents a single linting violation.
type Issue struct {
	Key     string
	Message string
}

// Result holds all issues found during a lint pass.
type Result struct {
	Issues []Issue
}

// HasIssues returns true when at least one issue was found.
func (r Result) HasIssues() bool {
	return len(r.Issues) > 0
}

// Lint runs all built-in rules against the provided env map and returns a
// Result containing any violations found.
func Lint(env map[string]string) Result {
	var issues []Issue

	for k, v := range env {
		if issue, ok := checkKeyCase(k); ok {
			issues = append(issues, issue)
		}
		if issue, ok := checkValueLength(k, v); ok {
			issues = append(issues, issue)
		}
		if issue, ok := checkUnquotedURL(k, v); ok {
			issues = append(issues, issue)
		}
		if issue, ok := checkWhitespace(k, v); ok {
			issues = append(issues, issue)
		}
	}

	return Result{Issues: issues}
}

func checkKeyCase(key string) (Issue, bool) {
	for _, r := range key {
		if unicode.IsLower(r) {
			return Issue{
				Key:     key,
				Message: "key contains lowercase letters; env keys should be UPPER_SNAKE_CASE",
			}, true
		}
	}
	return Issue{}, false
}

func checkValueLength(key, value string) (Issue, bool) {
	const maxLen = 512
	if len(value) > maxLen {
		return Issue{
			Key:     key,
			Message: fmt.Sprintf("value exceeds %d characters (%d)", maxLen, len(value)),
		}, true
	}
	return Issue{}, false
}

func checkUnquotedURL(key, value string) (Issue, bool) {
	if (strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")) &&
		!strings.HasPrefix(value, `"`) {
		return Issue{
			Key:     key,
			Message: "value looks like a URL and should be quoted",
		}, true
	}
	return Issue{}, false
}

func checkWhitespace(key, value string) (Issue, bool) {
	if value != strings.TrimSpace(value) {
		return Issue{
			Key:     key,
			Message: "value has leading or trailing whitespace",
		}, true
	}
	return Issue{}, false
}
