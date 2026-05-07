package patcher_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/patcher"
)

func TestApply_AddedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	entries := []differ.DiffEntry{
		{Type: differ.Added, Key: "NEW_KEY", NewValue: "new"},
	}
	result, err := patcher.Apply(base, entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["NEW_KEY"] != "new" {
		t.Errorf("expected NEW_KEY=new, got %q", result["NEW_KEY"])
	}
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar to be preserved, got %q", result["FOO"])
	}
}

func TestApply_RemovedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD": "gone"}
	entries := []differ.DiffEntry{
		{Type: differ.Removed, Key: "OLD", OldValue: "gone"},
	}
	result, err := patcher.Apply(base, entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["OLD"]; ok {
		t.Error("expected OLD to be removed")
	}
}

func TestApply_ChangedKey(t *testing.T) {
	base := map[string]string{"FOO": "old"}
	entries := []differ.DiffEntry{
		{Type: differ.Changed, Key: "FOO", OldValue: "old", NewValue: "new"},
	}
	result, err := patcher.Apply(base, entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %q", result["FOO"])
	}
}

func TestApply_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	entries := []differ.DiffEntry{
		{Type: differ.Added, Key: "EXTRA", NewValue: "val"},
	}
	_, _ = patcher.Apply(base, entries)
	if _, ok := base["EXTRA"]; ok {
		t.Error("Apply mutated the base map")
	}
}

func TestApply_SkipRemoved(t *testing.T) {
	base := map[string]string{"KEEP": "yes", "REMOVE": "no"}
	entries := []differ.DiffEntry{
		{Type: differ.Removed, Key: "REMOVE", OldValue: "no"},
	}
	result, err := patcher.Apply(base, entries, patcher.SkipRemoved())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["REMOVE"] != "no" {
		t.Error("expected REMOVE to be kept when SkipRemoved is set")
	}
}

func TestApply_SkipAdded(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	entries := []differ.DiffEntry{
		{Type: differ.Added, Key: "NEW", NewValue: "val"},
	}
	result, err := patcher.Apply(base, entries, patcher.SkipAdded())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["NEW"]; ok {
		t.Error("expected NEW to be skipped when SkipAdded is set")
	}
}

func TestApply_ErrorOnMissingChangedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	entries := []differ.DiffEntry{
		{Type: differ.Changed, Key: "MISSING", OldValue: "x", NewValue: "y"},
	}
	_, err := patcher.Apply(base, entries)
	if err == nil {
		t.Error("expected error when changing a key not in base")
	}
}

func TestApply_ErrorOnMissingRemovedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	entries := []differ.DiffEntry{
		{Type: differ.Removed, Key: "GHOST", OldValue: "val"},
	}
	_, err := patcher.Apply(base, entries)
	if err == nil {
		t.Error("expected error when removing a key not in base")
	}
}
