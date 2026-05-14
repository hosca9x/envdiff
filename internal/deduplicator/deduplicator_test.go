package deduplicator

import (
	"testing"
)

func TestDeduplicate_NoDuplicates(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"C": "3"}

	r, err := Deduplicate(FirstWins, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Env) != 3 {
		t.Errorf("expected 3 keys, got %d", len(r.Env))
	}
	if len(r.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", r.Duplicates)
	}
}

func TestDeduplicate_FirstWins(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	r, err := Deduplicate(FirstWins, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", r.Env["KEY"])
	}
	if len(r.Duplicates["KEY"]) != 2 {
		t.Errorf("expected 2 duplicate values recorded, got %v", r.Duplicates["KEY"])
	}
}

func TestDeduplicate_LastWins(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	r, err := Deduplicate(LastWins, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", r.Env["KEY"])
	}
}

func TestDeduplicate_ErrorOnDuplicate(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	_, err := Deduplicate(ErrorOnDuplicate, a, b)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeduplicate_ErrorOnDuplicate_NoConflict(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}

	r, err := Deduplicate(ErrorOnDuplicate, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(r.Env))
	}
}

func TestDeduplicate_DoesNotMutateInputs(t *testing.T) {
	a := map[string]string{"KEY": "original"}
	b := map[string]string{"KEY": "override"}

	_, _ = Deduplicate(LastWins, a, b)

	if a["KEY"] != "original" {
		t.Errorf("input map a was mutated")
	}
}

func TestMustDeduplicate_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got none")
		}
	}()
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	MustDeduplicate(ErrorOnDuplicate, a, b)
}
