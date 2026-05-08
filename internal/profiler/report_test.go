package profiler_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/profiler"
)

func TestWriteReport_ContainsKeyStats(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "",
	}
	p := profiler.Analyze(env)
	var buf bytes.Buffer
	if err := profiler.WriteReport(&buf, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{
		"Total Keys",
		"Empty Values",
		"Sensitive Keys",
		"Avg Value Length",
		"Longest Key",
		"Shortest Key",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestWriteReport_PrefixGroupsListed(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://",
	}
	p := profiler.Analyze(env)
	var buf bytes.Buffer
	_ = profiler.WriteReport(&buf, p)
	out := buf.String()
	if !strings.Contains(out, "APP") {
		t.Error("expected APP prefix in output")
	}
	if !strings.Contains(out, "DB") {
		t.Error("expected DB prefix in output")
	}
}

func TestWriteReport_SharedPrefixes(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	p := profiler.Analyze(env)
	var buf bytes.Buffer
	_ = profiler.WriteReport(&buf, p)
	out := buf.String()
	if !strings.Contains(out, "Shared Prefixes") {
		t.Error("expected Shared Prefixes section in output")
	}
}

func TestWriteReport_EmptyMap(t *testing.T) {
	p := profiler.Analyze(map[string]string{})
	var buf bytes.Buffer
	if err := profiler.WriteReport(&buf, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Total Keys       : 0") {
		t.Errorf("expected zero total keys in output, got: %s", out)
	}
}
