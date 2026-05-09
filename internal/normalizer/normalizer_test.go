package normalizer

import (
	"testing"
)

func TestNormalize_UpperSnake(t *testing.T) {
	input := map[string]string{
		"myAppKey": "value1",
		"db-host":  "localhost",
		"api.port": "8080",
	}
	got, err := Normalize(input, Options{Strategy: UpperSnake})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := map[string]string{
		"MY_APP_KEY": "value1",
		"DB_HOST":    "localhost",
		"API_PORT":   "8080",
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("key %q: got %q, want %q", k, got[k], want)
		}
	}
}

func TestNormalize_LowerSnake(t *testing.T) {
	input := map[string]string{"MyKey": "v", "DB-HOST": "h"}
	got, err := Normalize(input, Options{Strategy: LowerSnake})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["my_key"] != "v" {
		t.Errorf("expected my_key=v, got %v", got)
	}
	if got["db_host"] != "h" {
		t.Errorf("expected db_host=h, got %v", got)
	}
}

func TestNormalize_TrimPrefix(t *testing.T) {
	input := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "9000",
		"DB_PASS":  "secret",
	}
	got, err := Normalize(input, Options{Strategy: TrimPrefix, Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost")
	}
	if got["PORT"] != "9000" {
		t.Errorf("expected PORT=9000")
	}
	// Key without prefix should remain unchanged.
	if got["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret")
	}
}

func TestNormalize_DoesNotMutateOriginal(t *testing.T) {
	input := map[string]string{"myKey": "val"}
	_, err := Normalize(input, Options{Strategy: UpperSnake})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := input["myKey"]; !ok {
		t.Error("original map was mutated")
	}
}

func TestNormalize_NilMap_ReturnsError(t *testing.T) {
	_, err := Normalize(nil, Options{})
	if err == nil {
		t.Error("expected error for nil map")
	}
}

func TestMustNormalize_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil map")
		}
	}()
	MustNormalize(nil, Options{})
}
