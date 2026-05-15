package cloner_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/cloner"
)

func TestClone_CopiesSrc(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	out, err := cloner.Clone(src, nil, cloner.SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Errorf("expected src keys to be present, got %v", out)
	}
}

func TestClone_SkipExisting_KeepsDst(t *testing.T) {
	src := map[string]string{"KEY": "from_src"}
	dst := map[string]string{"KEY": "from_dst"}
	out, err := cloner.Clone(src, dst, cloner.SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "from_dst" {
		t.Errorf("expected dst value to be kept, got %q", out["KEY"])
	}
}

func TestClone_OverwriteExisting_UsesSrc(t *testing.T) {
	src := map[string]string{"KEY": "from_src"}
	dst := map[string]string{"KEY": "from_dst"}
	out, err := cloner.Clone(src, dst, cloner.OverwriteExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "from_src" {
		t.Errorf("expected src value to win, got %q", out["KEY"])
	}
}

func TestClone_ErrorOnConflict_ReturnsError(t *testing.T) {
	src := map[string]string{"KEY": "src_val"}
	dst := map[string]string{"KEY": "dst_val"}
	_, err := cloner.Clone(src, dst, cloner.ErrorOnConflict)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestClone_ErrorOnConflict_NoConflict(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{"B": "2"}
	out, err := cloner.Clone(src, dst, cloner.ErrorOnConflict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Errorf("expected both keys present, got %v", out)
	}
}

func TestClone_DoesNotMutateSrc(t *testing.T) {
	src := map[string]string{"X": "orig"}
	dst := map[string]string{}
	out, _ := cloner.Clone(src, dst, cloner.OverwriteExisting)
	out["X"] = "mutated"
	if src["X"] != "orig" {
		t.Errorf("src was mutated")
	}
}

func TestClone_DoesNotMutateDst(t *testing.T) {
	src := map[string]string{"X": "src_val"}
	dst := map[string]string{"Y": "dst_val"}
	out, _ := cloner.Clone(src, dst, cloner.OverwriteExisting)
	out["Y"] = "mutated"
	if dst["Y"] != "dst_val" {
		t.Errorf("dst was mutated")
	}
}

func TestClone_NilSrc_ReturnsError(t *testing.T) {
	_, err := cloner.Clone(nil, map[string]string{}, cloner.SkipExisting)
	if err == nil {
		t.Fatal("expected error for nil src")
	}
}

func TestMustClone_PanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got none")
		}
	}()
	cloner.MustClone(
		map[string]string{"K": "v1"},
		map[string]string{"K": "v2"},
		cloner.ErrorOnConflict,
	)
}
