package reconciler_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/reconciler"
)

func TestReconcile_AddsMissingKeys(t *testing.T) {
	target := map[string]string{"A": "1"}
	source := map[string]string{"A": "1", "B": "2"}

	res, err := reconciler.Reconcile(target, source, reconciler.StrategySourceWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["B"] != "2" {
		t.Errorf("expected B=2, got %q", res.Merged["B"])
	}
	if len(res.Applied) != 1 || res.Applied[0].Type != differ.Added {
		t.Errorf("expected one Applied entry of type Added")
	}
}

func TestReconcile_SourceWins_OverwritesChanged(t *testing.T) {
	target := map[string]string{"A": "old"}
	source := map[string]string{"A": "new"}

	res, err := reconciler.Reconcile(target, source, reconciler.StrategySourceWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["A"] != "new" {
		t.Errorf("expected A=new, got %q", res.Merged["A"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected one applied change")
	}
}

func TestReconcile_TargetWins_KeepsExisting(t *testing.T) {
	target := map[string]string{"A": "old"}
	source := map[string]string{"A": "new"}

	res, err := reconciler.Reconcile(target, source, reconciler.StrategyTargetWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["A"] != "old" {
		t.Errorf("expected A=old, got %q", res.Merged["A"])
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected one skipped change")
	}
}

func TestReconcile_AddOnly_SkipsChanges(t *testing.T) {
	target := map[string]string{"A": "1"}
	source := map[string]string{"A": "99", "B": "2"}

	res, err := reconciler.Reconcile(target, source, reconciler.StrategyAddOnly)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["A"] != "1" {
		t.Errorf("expected A unchanged, got %q", res.Merged["A"])
	}
	if res.Merged["B"] != "2" {
		t.Errorf("expected B=2 to be added")
	}
}

func TestReconcile_DoesNotMutateTarget(t *testing.T) {
	target := map[string]string{"A": "1"}
	source := map[string]string{"A": "2", "B": "3"}

	_, err := reconciler.Reconcile(target, source, reconciler.StrategySourceWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target["A"] != "1" {
		t.Errorf("target was mutated: A=%q", target["A"])
	}
	if _, ok := target["B"]; ok {
		t.Errorf("target was mutated: unexpected key B")
	}
}

func TestReconcile_NilTargetReturnsError(t *testing.T) {
	_, err := reconciler.Reconcile(nil, map[string]string{}, reconciler.StrategySourceWins)
	if err == nil {
		t.Error("expected error for nil target")
	}
}

func TestReconcile_NilSourceReturnsError(t *testing.T) {
	_, err := reconciler.Reconcile(map[string]string{}, nil, reconciler.StrategySourceWins)
	if err == nil {
		t.Error("expected error for nil source")
	}
}
