package validator_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/validator"
)

func TestValidate_ValidEnv(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"PORT":         "8080",
		"_PRIVATE":     "value",
	}
	issues := validator.Validate(env)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestValidate_InvalidKeyFormat(t *testing.T) {
	env := map[string]string{
		"123INVALID": "value",
	}
	issues := validator.Validate(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "error" {
		t.Errorf("expected error severity, got %s", issues[0].Severity)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	env := map[string]string{
		"MY_KEY": "",
	}
	issues := validator.Validate(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "warning" {
		t.Errorf("expected warning severity, got %s", issues[0].Severity)
	}
}

func TestValidate_UnresolvedPlaceholder(t *testing.T) {
	env := map[string]string{
		"API_URL": "https://example.com/${PATH}",
	}
	issues := validator.Validate(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "warning" {
		t.Errorf("expected warning, got %s", issues[0].Severity)
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"GOOD_KEY":  "good_value",
		"bad-key":   "",
		"EMPTY_KEY": "",
	}
	issues := validator.Validate(env)
	// bad-key triggers error + warning (invalid + empty); EMPTY_KEY triggers warning
	if len(issues) < 3 {
		t.Errorf("expected at least 3 issues, got %d", len(issues))
	}
}

func TestHasErrors_WithError(t *testing.T) {
	issues := []validator.Issue{
		{Key: "X", Severity: "warning", Message: "empty"},
		{Key: "Y", Severity: "error", Message: "bad key"},
	}
	if !validator.HasErrors(issues) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_OnlyWarnings(t *testing.T) {
	issues := []validator.Issue{
		{Key: "X", Severity: "warning", Message: "empty"},
	}
	if validator.HasErrors(issues) {
		t.Error("expected HasErrors to return false")
	}
}

func TestHasErrors_Empty(t *testing.T) {
	if validator.HasErrors(nil) {
		t.Error("expected HasErrors to return false for nil slice")
	}
	if validator.HasErrors([]validator.Issue{}) {
		t.Error("expected HasErrors to return false for empty slice")
	}
}

func TestIssue_String(t *testing.T) {
	i := validator.Issue{Key: "MY_KEY", Severity: "error", Message: "bad format"}
	got := i.String()
	expected := "[ERROR] MY_KEY: bad format"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestIssue_String_Warning(t *testing.T) {
	i := validator.Issue{Key: "SOME_KEY", Severity: "warning", Message: "empty value"}
	got := i.String()
	expected := "[WARNING] SOME_KEY: empty value"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
