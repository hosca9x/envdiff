package grouper

import (
	"testing"
)

func TestGroupBy_ByPrefix_BasicPartition(t *testing.T) {
	env := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"APP_NAME":    "envdiff",
		"APP_VERSION": "1.0",
		"STANDALONE":  "yes",
	}
	groups, err := GroupBy(env, ByPrefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	groupMap := toMap(groups)
	if len(groupMap["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(groupMap["DB"]))
	}
	if len(groupMap["APP"]) != 2 {
		t.Errorf("expected 2 APP keys, got %d", len(groupMap["APP"]))
	}
	if len(groupMap["_other"]) != 1 {
		t.Errorf("expected 1 _other key, got %d", len(groupMap["_other"]))
	}
}

func TestGroupBy_ByLength_Buckets(t *testing.T) {
	env := map[string]string{
		"A":                       "1",  // short
		"MEDIUM_KEY":              "2",  // medium (10 chars)
		"A_VERY_LONG_KEY_NAME_XY": "3",  // long (23 chars)
	}
	groups, err := GroupBy(env, ByLength)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	groupMap := toMap(groups)
	if _, ok := groupMap["short"]["A"]; !ok {
		t.Error("expected A in short bucket")
	}
	if _, ok := groupMap["medium"]["MEDIUM_KEY"]; !ok {
		t.Error("expected MEDIUM_KEY in medium bucket")
	}
	if _, ok := groupMap["long"]["A_VERY_LONG_KEY_NAME_XY"]; !ok {
		t.Error("expected long key in long bucket")
	}
}

func TestGroupBy_NilMap_ReturnsError(t *testing.T) {
	_, err := GroupBy(nil, ByPrefix)
	if err == nil {
		t.Error("expected error for nil map")
	}
}

func TestGroupBy_EmptyMap_ReturnsEmptySlice(t *testing.T) {
	groups, err := GroupBy(map[string]string{}, ByPrefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestGroupBy_SortedGroupNames(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "1",
		"A_KEY": "2",
		"M_KEY": "3",
	}
	groups, err := GroupBy(env, ByPrefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if groups[0].Name != "A" || groups[1].Name != "M" || groups[2].Name != "Z" {
		t.Errorf("groups not sorted: %v", groupNames(groups))
	}
}

func TestMustGroupBy_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil map")
		}
	}()
	MustGroupBy(nil, ByPrefix)
}

// helpers

func toMap(groups []Group) map[string]map[string]string {
	m := map[string]map[string]string{}
	for _, g := range groups {
		m[g.Name] = g.Keys
	}
	return m
}

func groupNames(groups []Group) []string {
	names := make([]string, len(groups))
	for i, g := range groups {
		names[i] = g.Name
	}
	return names
}
