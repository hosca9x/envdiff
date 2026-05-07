package exporter_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/exporter"
)

func TestWrite_DefaultFormat(t *testing.T) {
	e := exporter.New()
	var buf bytes.Buffer
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := e.Write(&buf, env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO=bar\n") {
		t.Errorf("expected FOO=bar in output, got:\n%s", out)
	}
	if !strings.Contains(out, "BAZ=qux\n") {
		t.Errorf("expected BAZ=qux in output, got:\n%s", out)
	}
}

func TestWrite_ExportFormat(t *testing.T) {
	e := exporter.NewWithOptions(exporter.FormatExport, false)
	var buf bytes.Buffer
	env := map[string]string{"PORT": "8080"}
	if err := e.Write(&buf, env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "export PORT=8080\n" {
		t.Errorf("expected 'export PORT=8080', got %q", got)
	}
}

func TestWrite_QuotedValues(t *testing.T) {
	e := exporter.NewWithOptions(exporter.FormatEnv, true)
	var buf bytes.Buffer
	env := map[string]string{"MSG": "hello world"}
	if err := e.Write(&buf, env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != `MSG="hello world"`+"\n" {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestWrite_SortedOutput(t *testing.T) {
	e := exporter.New()
	var buf bytes.Buffer
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	if err := e.Write(&buf, env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "A_KEY=a" || lines[1] != "M_KEY=m" || lines[2] != "Z_KEY=z" {
		t.Errorf("output not sorted: %v", lines)
	}
}

func TestWrite_EmptyMap(t *testing.T) {
	e := exporter.New()
	var buf bytes.Buffer
	if err := e.Write(&buf, map[string]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty map")
	}
}

func TestWriteFile_CreatesFile(t *testing.T) {
	e := exporter.New()
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	env := map[string]string{"KEY": "value"}
	if err := e.WriteFile(path, env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	if string(data) != "KEY=value\n" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}
