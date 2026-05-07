package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/snapshot"
)

func TestNew_CopiesEnv(t *testing.T) {
	original := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.New("test", original)

	original["FOO"] = "mutated"

	if s.Env["FOO"] != "bar" {
		t.Errorf("expected snapshot to be independent of original map, got %q", s.Env["FOO"])
	}
}

func TestNew_SetsLabelAndTimestamp(t *testing.T) {
	s := snapshot.New("staging", map[string]string{})

	if s.Label != "staging" {
		t.Errorf("expected label %q, got %q", "staging", s.Label)
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"API_KEY":  "secret123",
	}

	original := snapshot.New("prod", env)

	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(original, tmp); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := snapshot.Load(tmp)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Label != original.Label {
		t.Errorf("label mismatch: got %q, want %q", loaded.Label, original.Label)
	}
	if len(loaded.Env) != len(original.Env) {
		t.Errorf("env length mismatch: got %d, want %d", len(loaded.Env), len(original.Env))
	}
	for k, v := range original.Env {
		if loaded.Env[k] != v {
			t.Errorf("key %q: got %q, want %q", k, loaded.Env[k], v)
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	s := snapshot.New("test", map[string]string{"X": "1"})
	err := snapshot.Save(s, "/nonexistent/dir/snap.json")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestNew_EmptyEnv(t *testing.T) {
	s := snapshot.New("empty", map[string]string{})
	if s.Env == nil {
		t.Error("expected non-nil env map")
	}
	if len(s.Env) != 0 {
		t.Errorf("expected empty env, got %d keys", len(s.Env))
	}
}

func writeTempSnapshot(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()
	_, _ = f.WriteString(content)
	return f.Name()
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := writeTempSnapshot(t, "{invalid json")
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
