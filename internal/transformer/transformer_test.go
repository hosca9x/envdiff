package transformer_test

import (
	"errors"
	"testing"

	"github.com/your-org/envdiff/internal/transformer"
)

func TestApply_NoTransforms(t *testing.T) {
	input := map[string]string{"KEY": "value", "FOO": "bar"}
	tr := transformer.New()
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" || out["FOO"] != "bar" {
		t.Errorf("expected unchanged map, got %v", out)
	}
}

func TestApply_TrimSpace(t *testing.T) {
	input := map[string]string{"KEY": "  hello  ", "B": "\tworld\t"}
	tr := transformer.New(transformer.TrimSpace())
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello" {
		t.Errorf("KEY: expected %q got %q", "hello", out["KEY"])
	}
	if out["B"] != "world" {
		t.Errorf("B: expected %q got %q", "world", out["B"])
	}
}

func TestApply_ToUpper(t *testing.T) {
	input := map[string]string{"KEY": "hello"}
	tr := transformer.New(transformer.ToUpper())
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "HELLO" {
		t.Errorf("expected HELLO got %q", out["KEY"])
	}
}

func TestApply_ToLower(t *testing.T) {
	input := map[string]string{"KEY": "HELLO"}
	tr := transformer.New(transformer.ToLower())
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello" {
		t.Errorf("expected hello got %q", out["KEY"])
	}
}

func TestApply_ReplaceValue(t *testing.T) {
	input := map[string]string{"DB_URL": "localhost:5432", "CACHE": "localhost:6379"}
	tr := transformer.New(transformer.ReplaceValue("localhost", "prod-host"))
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_URL"] != "prod-host:5432" {
		t.Errorf("DB_URL: expected prod-host:5432 got %q", out["DB_URL"])
	}
	if out["CACHE"] != "prod-host:6379" {
		t.Errorf("CACHE: expected prod-host:6379 got %q", out["CACHE"])
	}
}

func TestApply_ChainedTransforms(t *testing.T) {
	input := map[string]string{"KEY": "  hello  "}
	tr := transformer.New(transformer.TrimSpace(), transformer.ToUpper())
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "HELLO" {
		t.Errorf("expected HELLO got %q", out["KEY"])
	}
}

func TestApply_ErrorPropagates(t *testing.T) {
	errFn := func(key, value string) (string, error) {
		if key == "BAD" {
			return "", errors.New("bad key")
		}
		return value, nil
	}
	input := map[string]string{"BAD": "val"}
	tr := transformer.New(errFn)
	_, err := tr.Apply(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	input := map[string]string{"KEY": "hello"}
	tr := transformer.New(transformer.ToUpper())
	_, _ = tr.Apply(input)
	if input["KEY"] != "hello" {
		t.Errorf("original map was mutated")
	}
}
