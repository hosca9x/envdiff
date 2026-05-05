package masker_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/masker"
)

func TestIsSensitive_KnownPatterns(t *testing.T) {
	m := masker.New()
	sensitiveKeys := []string{
		"DB_PASSWORD",
		"API_KEY",
		"AUTH_TOKEN",
		"AWS_SECRET",
		"PRIVATE_KEY",
		"DATABASE_URL",
		"app_secret",
		"stripe_api_key",
	}
	for _, key := range sensitiveKeys {
		if !m.IsSensitive(key) {
			t.Errorf("expected key %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	m := masker.New()
	safeKeys := []string{
		"APP_ENV",
		"PORT",
		"LOG_LEVEL",
		"REGION",
		"DEBUG",
	}
	for _, key := range safeKeys {
		if m.IsSensitive(key) {
			t.Errorf("expected key %q to NOT be sensitive", key)
		}
	}
}

func TestMaskMap_MasksSensitiveValues(t *testing.T) {
	m := masker.New()
	input := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "supersecret",
		"PORT":        "8080",
		"API_KEY":     "abc123",
	}
	result := m.MaskMap(input)

	if result["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV to be unmasked, got %q", result["APP_ENV"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT to be unmasked, got %q", result["PORT"])
	}
	if result["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD to be masked, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "***" {
		t.Errorf("expected API_KEY to be masked, got %q", result["API_KEY"])
	}
}

func TestMaskMap_DoesNotMutateOriginal(t *testing.T) {
	m := masker.New()
	input := map[string]string{
		"DB_PASSWORD": "supersecret",
	}
	_ = m.MaskMap(input)
	if input["DB_PASSWORD"] != "supersecret" {
		t.Error("original map was mutated")
	}
}

func TestNewWithOptions_CustomPatterns(t *testing.T) {
	m := masker.NewWithOptions([]string{"CUSTOM_SENSITIVE"}, "[REDACTED]")
	if !m.IsSensitive("MY_CUSTOM_SENSITIVE_KEY") {
		t.Error("expected custom pattern to match")
	}
	val := m.MaskValue("MY_CUSTOM_SENSITIVE_KEY", "value")
	if val != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", val)
	}
}
