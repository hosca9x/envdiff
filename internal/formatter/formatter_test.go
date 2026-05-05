package formatter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/formatter"
)

func TestWriteText_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatText, &buf)
	if err := f.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected no-diff message, got: %q", buf.String())
	}
}

func TestWriteText_Added(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatText, &buf)
	diffs := []differ.DiffEntry{
		{Key: "NEW_KEY", Type: differ.Added, NewValue: "value"},
	}
	f.Write(diffs)
	if !strings.Contains(buf.String(), "+ NEW_KEY=value") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestWriteText_Removed(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatText, &buf)
	diffs := []differ.DiffEntry{
		{Key: "OLD_KEY", Type: differ.Removed, OldValue: "old"},
	}
	f.Write(diffs)
	if !strings.Contains(buf.String(), "- OLD_KEY=old") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestWriteText_Changed(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatText, &buf)
	diffs := []differ.DiffEntry{
		{Key: "HOST", Type: differ.Changed, OldValue: "localhost", NewValue: "prod.example.com"},
	}
	f.Write(diffs)
	out := buf.String()
	if !strings.Contains(out, "~ HOST") || !strings.Contains(out, "localhost") || !strings.Contains(out, "prod.example.com") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestWriteJSON_Output(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatJSON, &buf)
	diffs := []differ.DiffEntry{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "db.prod"},
	}
	if err := f.Write(diffs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"key\"") || !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected JSON output with key, got: %q", out)
	}
}

func TestWriteText_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatText, &buf)
	diffs := []differ.DiffEntry{
		{Key: "Z_KEY", Type: differ.Added, NewValue: "z"},
		{Key: "A_KEY", Type: differ.Added, NewValue: "a"},
	}
	f.Write(diffs)
	out := buf.String()
	if strings.Index(out, "A_KEY") > strings.Index(out, "Z_KEY") {
		t.Errorf("expected sorted output, got: %q", out)
	}
}
