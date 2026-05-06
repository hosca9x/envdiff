package sorter

import (
	"testing"
)

var sampleEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "envdiff",
	"APP_VERSION": "1.0",
	"ZEBRA":       "true",
	"ALPHA":       "false",
}

func TestSort_Alphabetical(t *testing.T) {
	keys := Sort(sampleEnv, Alphabetical)
	for i := 1; i < len(keys); i++ {
		if keys[i-1] > keys[i] {
			t.Errorf("expected ascending order, got %q before %q", keys[i-1], keys[i])
		}
	}
}

func TestSort_Reverse(t *testing.T) {
	keys := Sort(sampleEnv, Reverse)
	for i := 1; i < len(keys); i++ {
		if keys[i-1] < keys[i] {
			t.Errorf("expected descending order, got %q before %q", keys[i-1], keys[i])
		}
	}
}

func TestSort_GroupedByPrefix(t *testing.T) {
	keys := Sort(sampleEnv, GroupedByPrefix)
	// APP_ group should appear before DB_ group (alphabetically)
	appIdx, dbIdx := -1, -1
	for i, k := range keys {
		if k == "APP_NAME" {
			appIdx = i
		}
		if k == "DB_HOST" {
			dbIdx = i
		}
	}
	if appIdx == -1 || dbIdx == -1 {
		t.Fatal("expected APP_NAME and DB_HOST to be present")
	}
	if appIdx > dbIdx {
		t.Errorf("expected APP_ group before DB_ group, got APP_NAME at %d, DB_HOST at %d", appIdx, dbIdx)
	}
}

func TestSort_EmptyMap(t *testing.T) {
	keys := Sort(map[string]string{}, Alphabetical)
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}

func TestExtractPrefix_WithUnderscore(t *testing.T) {
	if p := extractPrefix("DB_HOST"); p != "DB" {
		t.Errorf("expected DB, got %q", p)
	}
}

func TestExtractPrefix_WithoutUnderscore(t *testing.T) {
	if p := extractPrefix("ZEBRA"); p != "ZEBRA" {
		t.Errorf("expected ZEBRA, got %q", p)
	}
}
