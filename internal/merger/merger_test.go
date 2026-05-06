package merger_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/merger"
)

func TestMerge_FirstWins(t *testing.T) {
	a := map[string]string{"KEY": "alpha", "ONLY_A": "yes"}
	b := map[string]string{"KEY": "beta", "ONLY_B": "yes"}

	got, err := merger.Merge(merger.FirstWins, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "alpha" {
		t.Errorf("FirstWins: expected 'alpha', got %q", got["KEY"])
	}
	if got["ONLY_B"] != "yes" {
		t.Errorf("expected ONLY_B to be merged in")
	}
}

func TestMerge_LastWins(t *testing.T) {
	a := map[string]string{"KEY": "alpha"}
	b := map[string]string{"KEY": "beta"}

	got, err := merger.Merge(merger.LastWins, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "beta" {
		t.Errorf("LastWins: expected 'beta', got %q", got["KEY"])
	}
}

func TestMerge_ErrorOnConflict_NoConflict(t *testing.T) {
	a := map[string]string{"KEY": "same", "A": "1"}
	b := map[string]string{"KEY": "same", "B": "2"}

	got, err := merger.Merge(merger.ErrorOnConflict, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "same" {
		t.Errorf("expected 'same', got %q", got["KEY"])
	}
}

func TestMerge_ErrorOnConflict_WithConflict(t *testing.T) {
	a := map[string]string{"KEY": "alpha"}
	b := map[string]string{"KEY": "beta"}

	_, err := merger.Merge(merger.ErrorOnConflict, a, b)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	ce, ok := err.(*merger.ConflictError)
	if !ok {
		t.Fatalf("expected *ConflictError, got %T", err)
	}
	if ce.Key != "KEY" {
		t.Errorf("expected conflict key 'KEY', got %q", ce.Key)
	}
}

func TestMerge_EmptyInputs(t *testing.T) {
	got, err := merger.Merge(merger.FirstWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestMerge_DoesNotMutateInputs(t *testing.T) {
	a := map[string]string{"KEY": "alpha"}
	b := map[string]string{"KEY": "beta"}

	_, _ = merger.Merge(merger.LastWins, a, b)

	if a["KEY"] != "alpha" {
		t.Errorf("input map a was mutated")
	}
}
