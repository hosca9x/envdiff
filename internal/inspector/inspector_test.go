package inspector

import (
	"testing"
)

func TestInspect_TotalKeys(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://localhost/db",
	}
	r := Inspect(env)
	if r.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", r.TotalKeys)
	}
}

func TestInspect_EmptyValues(t *testing.T) {
	env := map[string]string{
		"KEY_A": "",
		"KEY_B": "value",
		"KEY_C": "",
	}
	r := Inspect(env)
	if r.EmptyValues != 2 {
		t.Errorf("expected 2 empty values, got %d", r.EmptyValues)
	}
}

func TestInspect_PrefixGroups(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "dbhost",
	}
	r := Inspect(env)
	if len(r.PrefixGroups["APP"]) != 2 {
		t.Errorf("expected 2 keys under APP prefix, got %d", len(r.PrefixGroups["APP"]))
	}
	if len(r.PrefixGroups["DB"]) != 1 {
		t.Errorf("expected 1 key under DB prefix, got %d", len(r.PrefixGroups["DB"]))
	}
}

func TestInspect_ValueHints(t *testing.T) {
	env := map[string]string{
		"ENABLED":  "true",
		"PORT":     "3000",
		"BASE_URL": "https://example.com",
		"NAME":     "myapp",
		"EMPTY":    "",
	}
	r := Inspect(env)
	cases := map[string]string{
		"ENABLED":  "bool",
		"PORT":     "number",
		"BASE_URL": "url",
		"NAME":     "string",
		"EMPTY":    "empty",
	}
	for key, want := range cases {
		got := r.ValueHints[key]
		if got != want {
			t.Errorf("key %q: expected hint %q, got %q", key, want, got)
		}
	}
}

func TestInspect_EmptyMap(t *testing.T) {
	r := Inspect(map[string]string{})
	if r.TotalKeys != 0 {
		t.Errorf("expected 0 total keys, got %d", r.TotalKeys)
	}
	if r.EmptyValues != 0 {
		t.Errorf("expected 0 empty values, got %d", r.EmptyValues)
	}
	if len(r.PrefixGroups) != 0 {
		t.Errorf("expected empty prefix groups")
	}
}

func TestExtractPrefix_NoUnderscore(t *testing.T) {
	if got := extractPrefix("HOSTNAME"); got != "HOSTNAME" {
		t.Errorf("expected HOSTNAME, got %q", got)
	}
}

func TestExtractPrefix_WithUnderscore(t *testing.T) {
	if got := extractPrefix("APP_SECRET"); got != "APP" {
		t.Errorf("expected APP, got %q", got)
	}
}
