package redactor_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/redactor"
)

func TestIsSensitive_KnownPatterns(t *testing.T) {
	r := redactor.New()
	sensitive := []string{
		"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "PRIVATE_KEY",
		"SECRET", "AWS_SECRET_ACCESS_KEY", "CERT_DATA",
	}
	for _, key := range sensitive {
		if !r.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	r := redactor.New()
	safe := []string{"APP_ENV", "PORT", "HOST", "LOG_LEVEL", "TIMEOUT"}
	for _, key := range safe {
		if r.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestRedact_ReplacesSensitiveValues(t *testing.T) {
	r := redactor.New()
	env := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
	}

	result := r.Redact(env)

	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be redacted, got %q", result["APP_ENV"])
	}
	if result["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY should be redacted, got %q", result["API_KEY"])
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	r := redactor.New()
	env := map[string]string{"DB_PASSWORD": "s3cr3t"}
	_ = r.Redact(env)
	if env["DB_PASSWORD"] != "s3cr3t" {
		t.Error("original map was mutated")
	}
}

func TestWithPlaceholder_CustomValue(t *testing.T) {
	r := redactor.New(redactor.WithPlaceholder("***"))
	env := map[string]string{"API_KEY": "abc"}
	result := r.Redact(env)
	if result["API_KEY"] != "***" {
		t.Errorf("expected '***', got %q", result["API_KEY"])
	}
}

func TestWithPatterns_AddsCustomPattern(t *testing.T) {
	r := redactor.New(redactor.WithPatterns("internal"))
	if !r.IsSensitive("INTERNAL_URL") {
		t.Error("expected INTERNAL_URL to be sensitive with custom pattern")
	}
}

func TestRedactValue_SensitiveAndSafe(t *testing.T) {
	r := redactor.New()
	if got := r.RedactValue("DB_PASSWORD", "hunter2"); got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
	if got := r.RedactValue("APP_ENV", "staging"); got != "staging" {
		t.Errorf("expected 'staging', got %q", got)
	}
}
