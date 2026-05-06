package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/differ"
	"github.com/yourorg/envdiff/internal/reporter"
)

func sampleEntries() []differ.DiffEntry {
	return []differ.DiffEntry{
		{Key: "NEW_KEY", Status: differ.Added, TargetValue: "val1"},
		{Key: "OLD_KEY", Status: differ.Removed, SourceValue: "old"},
		{Key: "CHANGED_KEY", Status: differ.Changed, SourceValue: "a", TargetValue: "b"},
		{Key: "ANOTHER_NEW", Status: differ.Added, TargetValue: "val2"},
	}
}

func TestSummarize_Counts(t *testing.T) {
	entries := sampleEntries()
	s := reporter.Summarize(entries)

	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Added != 2 {
		t.Errorf("expected Added=2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", s.Removed)
	}
	if s.Changed != 1 {
		t.Errorf("expected Changed=1, got %d", s.Changed)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := reporter.Summarize(nil)
	if s.Total != 0 || s.Added != 0 || s.Removed != 0 || s.Changed != 0 {
		t.Errorf("expected all zeros for empty entries, got %+v", s)
	}
}

func TestWriteSummary_Output(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)

	if err := r.WriteSummary(sampleEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "+2 added") {
		t.Errorf("expected '+2 added' in output, got: %s", out)
	}
	if !strings.Contains(out, "-1 removed") {
		t.Errorf("expected '-1 removed' in output, got: %s", out)
	}
	if !strings.Contains(out, "~1 changed") {
		t.Errorf("expected '~1 changed' in output, got: %s", out)
	}
}

func TestWriteKeyList_SortedAndGrouped(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)

	if err := r.WriteKeyList(sampleEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d: %v", len(lines), lines)
	}

	// ADDED keys come first and should be sorted
	if !strings.HasPrefix(lines[0], "[ADDED]") {
		t.Errorf("expected first line to be ADDED, got: %s", lines[0])
	}
	if !strings.Contains(lines[0], "ANOTHER_NEW") {
		t.Errorf("expected ANOTHER_NEW before NEW_KEY (sorted), got: %s", lines[0])
	}
}

func TestWriteKeyList_Empty(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)

	if err := r.WriteKeyList(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty entries, got: %s", buf.String())
	}
}
