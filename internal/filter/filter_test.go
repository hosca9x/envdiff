package filter_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/filter"
)

var sampleEnv = map[string]string{
	"APP_NAME":    "myapp",
	"APP_PORT":    "8080",
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"LOG_LEVEL":   "info",
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	result := filter.Filter(sampleEnv)
	if len(result) != len(sampleEnv) {
		t.Errorf("expected %d entries, got %d", len(sampleEnv), len(result))
	}
}

func TestFilter_WithPrefix_ReturnsMatching(t *testing.T) {
	result := filter.Filter(sampleEnv, filter.WithPrefixes("APP_"))
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
	if _, ok := result["APP_NAME"]; !ok {
		t.Error("expected APP_NAME in result")
	}
	if _, ok := result["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in result")
	}
}

func TestFilter_WithKeys_ReturnsExact(t *testing.T) {
	result := filter.Filter(sampleEnv, filter.WithKeys("DB_HOST", "LOG_LEVEL"))
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("unexpected value for DB_HOST: %s", result["DB_HOST"])
	}
}

func TestFilter_Exclude_RemovesMatching(t *testing.T) {
	result := filter.Filter(sampleEnv, filter.WithPrefixes("DB_"), filter.Exclude())
	if _, ok := result["DB_HOST"]; ok {
		t.Error("DB_HOST should have been excluded")
	}
	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded")
	}
	if len(result) != 3 {
		t.Errorf("expected 3 entries after exclusion, got %d", len(result))
	}
}

func TestFilter_DoesNotMutateOriginal(t *testing.T) {
	original := map[string]string{"A": "1", "B": "2", "C": "3"}
	result := filter.Filter(original, filter.WithKeys("A"))
	result["EXTRA"] = "injected"
	if _, ok := original["EXTRA"]; ok {
		t.Error("filter mutated the original map")
	}
}

func TestFilter_MultiplePrefixes(t *testing.T) {
	result := filter.Filter(sampleEnv, filter.WithPrefixes("APP_", "LOG_"))
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
}
