package profiler_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/profiler"
)

func TestAnalyze_TotalKeys(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://",
	}
	p := profiler.Analyze(env)
	if p.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", p.TotalKeys)
	}
}

func TestAnalyze_EmptyValues(t *testing.T) {
	env := map[string]string{
		"KEY_A": "",
		"KEY_B": "value",
		"KEY_C": "",
	}
	p := profiler.Analyze(env)
	if p.EmptyValues != 2 {
		t.Errorf("expected 2 empty values, got %d", p.EmptyValues)
	}
}

func TestAnalyze_SensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret123",
		"API_TOKEN":   "tok_abc",
		"APP_NAME":    "myapp",
	}
	p := profiler.Analyze(env)
	if p.SensitiveKeys != 2 {
		t.Errorf("expected 2 sensitive keys, got %d", p.SensitiveKeys)
	}
}

func TestAnalyze_PrefixGroups(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://",
	}
	p := profiler.Analyze(env)
	if p.PrefixGroups["APP"] != 2 {
		t.Errorf("expected APP prefix count 2, got %d", p.PrefixGroups["APP"])
	}
	if p.PrefixGroups["DB"] != 1 {
		t.Errorf("expected DB prefix count 1, got %d", p.PrefixGroups["DB"])
	}
}

func TestAnalyze_AvgValueLength(t *testing.T) {
	env := map[string]string{
		"A": "ab",   // 2
		"B": "abcd", // 4
	}
	p := profiler.Analyze(env)
	if p.AvgValueLength != 3.0 {
		t.Errorf("expected avg 3.0, got %f", p.AvgValueLength)
	}
}

func TestAnalyze_LongestAndShortestKey(t *testing.T) {
	env := map[string]string{
		"AB":         "x",
		"ABCDEFGHIJ": "y",
		"ABCDE":      "z",
	}
	p := profiler.Analyze(env)
	if p.LongestKey != "ABCDEFGHIJ" {
		t.Errorf("expected longest key ABCDEFGHIJ, got %s", p.LongestKey)
	}
	if p.ShortestKey != "AB" {
		t.Errorf("expected shortest key AB, got %s", p.ShortestKey)
	}
}

func TestAnalyze_EmptyMap(t *testing.T) {
	p := profiler.Analyze(map[string]string{})
	if p.TotalKeys != 0 {
		t.Errorf("expected 0 total keys, got %d", p.TotalKeys)
	}
	if p.PrefixGroups == nil {
		t.Error("expected non-nil PrefixGroups")
	}
}

func TestAnalyze_DuplicatePrefixes(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://",
	}
	p := profiler.Analyze(env)
	if len(p.DuplicatePrefixes) != 1 || p.DuplicatePrefixes[0] != "APP" {
		t.Errorf("expected [APP] duplicate prefixes, got %v", p.DuplicatePrefixes)
	}
}
