package comparator_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/comparator"
)

func TestCompare_IdenticalMaps(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "1", "B": "2"}

	r := comparator.Compare(src, dst)

	if len(r.MatchingKeys) != 2 {
		t.Errorf("expected 2 matching keys, got %d", len(r.MatchingKeys))
	}
	if r.SimilarityScore != 1.0 {
		t.Errorf("expected similarity 1.0, got %f", r.SimilarityScore)
	}
	if len(r.DriftedKeys) != 0 || len(r.OnlyInSource) != 0 || len(r.OnlyInTarget) != 0 {
		t.Error("expected no drifted or exclusive keys")
	}
}

func TestCompare_AllDifferent(t *testing.T) {
	src := map[string]string{"A": "old"}
	dst := map[string]string{"A": "new"}

	r := comparator.Compare(src, dst)

	if len(r.DriftedKeys) != 1 {
		t.Errorf("expected 1 drifted key, got %d", len(r.DriftedKeys))
	}
	if r.SimilarityScore != 0.0 {
		t.Errorf("expected similarity 0.0, got %f", r.SimilarityScore)
	}
}

func TestCompare_OnlyInSource(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "1"}

	r := comparator.Compare(src, dst)

	if len(r.OnlyInSource) != 1 || r.OnlyInSource[0] != "B" {
		t.Errorf("expected OnlyInSource=[B], got %v", r.OnlyInSource)
	}
}

func TestCompare_OnlyInTarget(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{"A": "1", "C": "3"}

	r := comparator.Compare(src, dst)

	if len(r.OnlyInTarget) != 1 || r.OnlyInTarget[0] != "C" {
		t.Errorf("expected OnlyInTarget=[C], got %v", r.OnlyInTarget)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	r := comparator.Compare(map[string]string{}, map[string]string{})

	if r.SimilarityScore != 0.0 {
		t.Errorf("expected similarity 0.0 for empty maps, got %f", r.SimilarityScore)
	}
	if len(r.MatchingKeys) != 0 {
		t.Error("expected no matching keys for empty maps")
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	src := map[string]string{"A": "1", "B": "old", "C": "3"}
	dst := map[string]string{"A": "1", "B": "new", "D": "4"}

	r := comparator.Compare(src, dst)

	if len(r.MatchingKeys) != 1 {
		t.Errorf("expected 1 matching key, got %d", len(r.MatchingKeys))
	}
	if len(r.DriftedKeys) != 1 {
		t.Errorf("expected 1 drifted key, got %d", len(r.DriftedKeys))
	}
	if len(r.OnlyInSource) != 1 {
		t.Errorf("expected 1 source-only key, got %d", len(r.OnlyInSource))
	}
	if len(r.OnlyInTarget) != 1 {
		t.Errorf("expected 1 target-only key, got %d", len(r.OnlyInTarget))
	}
	// 1 matching out of 4 total entries
	expected := 0.25
	if r.SimilarityScore != expected {
		t.Errorf("expected similarity %f, got %f", expected, r.SimilarityScore)
	}
}
