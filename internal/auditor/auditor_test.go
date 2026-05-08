package auditor

import (
	"strings"
	"testing"
	"time"

	"github.com/user/envdiff/internal/differ"
)

func fixedNow() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func newTestAuditor() *Auditor {
	a := New("ci-bot", "staging")
	a.now = fixedNow
	return a
}

func TestAudit_AddedEntry(t *testing.T) {
	a := newTestAuditor()
	diffs := []differ.Entry{{Key: "NEW_KEY", Type: differ.Added, NewValue: "hello"}}
	entries := a.Audit(diffs)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.ChangeType != "added" {
		t.Errorf("expected change_type=added, got %q", e.ChangeType)
	}
	if e.NewValue != "hello" {
		t.Errorf("expected new_value=hello, got %q", e.NewValue)
	}
	if e.OldValue != "" {
		t.Errorf("expected empty old_value, got %q", e.OldValue)
	}
}

func TestAudit_RemovedEntry(t *testing.T) {
	a := newTestAuditor()
	diffs := []differ.Entry{{Key: "OLD_KEY", Type: differ.Removed, OldValue: "bye"}}
	entries := a.Audit(diffs)
	e := entries[0]
	if e.ChangeType != "removed" {
		t.Errorf("expected removed, got %q", e.ChangeType)
	}
	if e.OldValue != "bye" {
		t.Errorf("expected old_value=bye, got %q", e.OldValue)
	}
}

func TestAudit_ChangedEntry(t *testing.T) {
	a := newTestAuditor()
	diffs := []differ.Entry{{Key: "DB_HOST", Type: differ.Changed, OldValue: "localhost", NewValue: "prod.db"}}
	entries := a.Audit(diffs)
	e := entries[0]
	if e.ChangeType != "changed" {
		t.Errorf("expected changed, got %q", e.ChangeType)
	}
	if e.OldValue != "localhost" || e.NewValue != "prod.db" {
		t.Errorf("unexpected values: old=%q new=%q", e.OldValue, e.NewValue)
	}
}

func TestAudit_SetsActorAndEnvironment(t *testing.T) {
	a := newTestAuditor()
	entries := a.Audit([]differ.Entry{{Key: "X", Type: differ.Added, NewValue: "1"}})
	if entries[0].Actor != "ci-bot" {
		t.Errorf("expected actor=ci-bot, got %q", entries[0].Actor)
	}
	if entries[0].Environment != "staging" {
		t.Errorf("expected environment=staging, got %q", entries[0].Environment)
	}
}

func TestAudit_EmptyDiffs(t *testing.T) {
	a := newTestAuditor()
	entries := a.Audit(nil)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestSummary_Added(t *testing.T) {
	e := Entry{
		Timestamp:   fixedNow(),
		Actor:       "alice",
		Key:         "API_KEY",
		ChangeType:  "added",
		NewValue:    "secret",
		Environment: "prod",
	}
	s := Summary(e)
	if !strings.Contains(s, "added") || !strings.Contains(s, "API_KEY") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestSummary_Changed(t *testing.T) {
	e := Entry{
		Timestamp:   fixedNow(),
		Actor:       "bob",
		Key:         "DB_URL",
		ChangeType:  "changed",
		OldValue:    "old",
		NewValue:    "new",
		Environment: "staging",
	}
	s := Summary(e)
	if !strings.Contains(s, "old") || !strings.Contains(s, "new") {
		t.Errorf("unexpected summary: %q", s)
	}
}
