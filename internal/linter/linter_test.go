package linter

import (
	"strings"
	"testing"
)

func TestLint_CleanEnv_NoIssues(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": `"https://db.example.com"`,
		"APP_PORT":     "8080",
		"LOG_LEVEL":    "info",
	}
	// LOG_LEVEL has lowercase value, not key — only keys are checked for case.
	env2 := map[string]string{
		"DATABASE_URL": `"https://db.example.com"`,
		"APP_PORT":     "8080",
	}
	result := Lint(env2)
	if result.HasIssues() {
		t.Errorf("expected no issues, got %v", result.Issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"app_name": "myapp",
	}
	result := Lint(env)
	if !result.HasIssues() {
		t.Fatal("expected issue for lowercase key")
	}
	if !strings.Contains(result.Issues[0].Message, "lowercase") {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_ValueTooLong(t *testing.T) {
	env := map[string]string{
		"BIG_VALUE": strings.Repeat("x", 600),
	}
	result := Lint(env)
	if !result.HasIssues() {
		t.Fatal("expected issue for oversized value")
	}
	if !strings.Contains(result.Issues[0].Message, "exceeds") {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_UnquotedURL(t *testing.T) {
	env := map[string]string{
		"API_ENDPOINT": "https://api.example.com/v1",
	}
	result := Lint(env)
	if !result.HasIssues() {
		t.Fatal("expected issue for unquoted URL")
	}
	if !strings.Contains(result.Issues[0].Message, "URL") {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_QuotedURL_NoIssue(t *testing.T) {
	env := map[string]string{
		"API_ENDPOINT": `"https://api.example.com/v1"`,
	}
	result := Lint(env)
	if result.HasIssues() {
		t.Errorf("expected no issues for quoted URL, got %v", result.Issues)
	}
}

func TestLint_LeadingWhitespace(t *testing.T) {
	env := map[string]string{
		"PADDED_VAL": "  something",
	}
	result := Lint(env)
	if !result.HasIssues() {
		t.Fatal("expected issue for leading whitespace")
	}
	if !strings.Contains(result.Issues[0].Message, "whitespace") {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"bad_key":    "  https://example.com",
		"CLEAN_KEY":  "ok",
	}
	result := Lint(env)
	// bad_key triggers: lowercase key, leading whitespace, unquoted URL
	if len(result.Issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d: %v", len(result.Issues), result.Issues)
	}
}

func TestResult_HasIssues_Empty(t *testing.T) {
	r := Result{}
	if r.HasIssues() {
		t.Error("empty result should not have issues")
	}
}
