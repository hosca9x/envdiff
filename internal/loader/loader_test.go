package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestLoadFile_BasicParsing(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	m, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Errorf("unexpected map contents: %v", m)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := loader.LoadFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFiles_MultipleFiles(t *testing.T) {
	p1 := writeTempEnv(t, "A=1\n")
	p2 := writeTempEnv(t, "B=2\n")

	maps, err := loader.LoadFiles([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(maps) != 2 {
		t.Fatalf("expected 2 maps, got %d", len(maps))
	}
	if maps[0]["A"] != "1" {
		t.Errorf("expected maps[0][A]=1, got %q", maps[0]["A"])
	}
	if maps[1]["B"] != "2" {
		t.Errorf("expected maps[1][B]=2, got %q", maps[1]["B"])
	}
}

func TestLoadFiles_StopsOnError(t *testing.T) {
	p1 := writeTempEnv(t, "A=1\n")
	_, err := loader.LoadFiles([]string{p1, "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error when one file is missing")
	}
}

func TestMustLoadFile_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing file")
		}
	}()
	loader.MustLoadFile("/nonexistent/.env")
}
