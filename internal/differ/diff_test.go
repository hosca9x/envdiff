package differ_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/differ"
)

func TestDiff_NoChanges(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	tgt := map[string]string{"FOO": "bar", "BAZ": "qux"}
	entries := differ.Diff(src, tgt)
	for _, e := range entries {
		if e.ChangeType != differ.Unchanged {
			t.Errorf("expected Unchanged for key %q, got %s", e.Key, e.ChangeType)
		}
	}
}

func TestDiff_AddedKey(t *testing.T) {
	src := map[string]string{"FOO": "bar", "NEW": "value"}
	tgt := map[string]string{"FOO": "bar"}
	entries := differ.Diff(src, tgt)
	found := findEntry(entries, "NEW")
	if found == nil {
		t.Fatal("expected entry for NEW")
	}
	if found.ChangeType != differ.Added {
		t.Errorf("expected Added, got %s", found.ChangeType)
	}
	if found.NewValue != "value" {
		t.Errorf("unexpected NewValue: %q", found.NewValue)
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	tgt := map[string]string{"FOO": "bar", "OLD": "gone"}
	entries := differ.Diff(src, tgt)
	found := findEntry(entries, "OLD")
	if found == nil {
		t.Fatal("expected entry for OLD")
	}
	if found.ChangeType != differ.Removed {
		t.Errorf("expected Removed, got %s", found.ChangeType)
	}
	if found.OldValue != "gone" {
		t.Errorf("unexpected OldValue: %q", found.OldValue)
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	src := map[string]string{"FOO": "new"}
	tgt := map[string]string{"FOO": "old"}
	entries := differ.Diff(src, tgt)
	found := findEntry(entries, "FOO")
	if found == nil {
		t.Fatal("expected entry for FOO")
	}
	if found.ChangeType != differ.Changed {
		t.Errorf("expected Changed, got %s", found.ChangeType)
	}
	if found.OldValue != "old" || found.NewValue != "new" {
		t.Errorf("unexpected values: old=%q new=%q", found.OldValue, found.NewValue)
	}
}

func TestDiff_EmptyMaps(t *testing.T) {
	entries := differ.Diff(map[string]string{}, map[string]string{})
	if len(entries) != 0 {
		t.Errorf("expected no entries, got %d", len(entries))
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	src := map[string]string{"Z": "1", "A": "2", "M": "3"}
	tgt := map[string]string{}
	entries := differ.Diff(src, tgt)
	for i := 1; i < len(entries); i++ {
		if entries[i].Key < entries[i-1].Key {
			t.Errorf("entries not sorted: %q before %q", entries[i-1].Key, entries[i].Key)
		}
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
