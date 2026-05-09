package renamer

import (
	"testing"
)

func TestRename_BasicRule(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value"}
	out, res, err := Rename(env, []Rule{{From: "OLD_KEY", To: "NEW_KEY"}}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", out["NEW_KEY"])
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been removed")
	}
	if res.Renamed["OLD_KEY"] != "NEW_KEY" {
		t.Errorf("expected renamed entry OLD_KEY->NEW_KEY")
	}
}

func TestRename_SkipsMissingKey(t *testing.T) {
	env := map[string]string{"A": "1"}
	_, res, err := Rename(env, []Rule{{From: "MISSING", To: "X"}}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING in skipped, got %v", res.Skipped)
	}
}

func TestRename_ConflictNoOverwrite(t *testing.T) {
	env := map[string]string{"OLD": "v1", "NEW": "v2"}
	out, res, err := Rename(env, []Rule{{From: "OLD", To: "NEW"}}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflict) != 1 || res.Conflict[0] != "NEW" {
		t.Errorf("expected conflict on NEW, got %v", res.Conflict)
	}
	// original keys must be untouched
	if out["OLD"] != "v1" || out["NEW"] != "v2" {
		t.Error("env should be unchanged on conflict without overwrite")
	}
}

func TestRename_ConflictWithOverwrite(t *testing.T) {
	env := map[string]string{"OLD": "v1", "NEW": "v2"}
	out, res, err := Rename(env, []Rule{{From: "OLD", To: "NEW"}}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflict) != 0 {
		t.Errorf("expected no conflicts with overwrite, got %v", res.Conflict)
	}
	if out["NEW"] != "v1" {
		t.Errorf("expected NEW=v1 after overwrite, got %q", out["NEW"])
	}
}

func TestRename_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	original := map[string]string{"FOO": "bar"}
	Rename(env, []Rule{{From: "FOO", To: "BAZ"}}, false)
	if env["FOO"] != original["FOO"] {
		t.Error("input env was mutated")
	}
}

func TestRename_NilEnvReturnsError(t *testing.T) {
	_, _, err := Rename(nil, []Rule{{From: "A", To: "B"}}, false)
	if err == nil {
		t.Error("expected error for nil env")
	}
}

func TestRename_SameFromTo_IsNoop(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	out, res, _ := Rename(env, []Rule{{From: "KEY", To: "KEY"}}, false)
	if len(res.Renamed) != 0 {
		t.Error("expected no renames for same from/to")
	}
	if out["KEY"] != "val" {
		t.Error("key should still be present")
	}
}
