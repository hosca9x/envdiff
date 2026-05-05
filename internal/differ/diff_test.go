package differ_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/differ"
)

func TestDiff_NoChanges(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := differ.Diff(a, b)
	if result.HasChanges() {
		t.Error("expected no changes, but HasChanges() returned true")
	}
	for _, e := range result.Entries {
		if e.Status != differ.StatusSame {
			t.Errorf("key %q: expected status %q, got %q", e.Key, differ.StatusSame, e.Status)
		}
	}
}

func TestDiff_AddedKey(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "NEW_KEY": "newval"}

	result := differ.Diff(a, b)
	if !result.HasChanges() {
		t.Fatal("expected changes, but HasChanges() returned false")
	}
	found := findEntry(result.Entries, "NEW_KEY")
	if found == nil {
		t.Fatal("expected entry for NEW_KEY")
	}
	if found.Status != differ.StatusAdded {
		t.Errorf("expected status %q, got %q", differ.StatusAdded, found.Status)
	}
	if found.ValueB != "newval" {
		t.Errorf("expected ValueB %q, got %q", "newval", found.ValueB)
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	a := map[string]string{"FOO": "bar", "OLD_KEY": "oldval"}
	b := map[string]string{"FOO": "bar"}

	result := differ.Diff(a, b)
	found := findEntry(result.Entries, "OLD_KEY")
	if found == nil {
		t.Fatal("expected entry for OLD_KEY")
	}
	if found.Status != differ.StatusRemoved {
		t.Errorf("expected status %q, got %q", differ.StatusRemoved, found.Status)
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}

	result := differ.Diff(a, b)
	found := findEntry(result.Entries, "FOO")
	if found == nil {
		t.Fatal("expected entry for FOO")
	}
	if found.Status != differ.StatusChanged {
		t.Errorf("expected status %q, got %q", differ.StatusChanged, found.Status)
	}
	if found.ValueA != "old" || found.ValueB != "new" {
		t.Errorf("unexpected values: ValueA=%q ValueB=%q", found.ValueA, found.ValueB)
	}
}

func TestDiff_EmptyMaps(t *testing.T) {
	result := differ.Diff(map[string]string{}, map[string]string{})
	if result.HasChanges() {
		t.Error("expected no changes for empty maps")
	}
}

func findEntry(entries []differ.Entry, key string) *differ.Entry {
	for i := range entries {
		if entries[i].Key == key {
			return &entries[i]
		}
	}
	return nil
}
