package parser

import (
	"strings"
	"testing"
)

func TestParse_BasicKeyValue(t *testing.T) {
	input := `APP_NAME=envdiff
DEBUG=true
PORT=8080
`
	em, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(em.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(em.Entries))
	}
	if e, ok := em.Get("APP_NAME"); !ok || e.Value != "envdiff" {
		t.Errorf("APP_NAME: got %q, want %q", e.Value, "envdiff")
	}
}

func TestParse_QuotedValues(t *testing.T) {
	input := `DB_URL="postgres://localhost/mydb"
SECRET='super secret'
`
	em, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e, _ := em.Get("DB_URL"); e.Value != "postgres://localhost/mydb" {
		t.Errorf("DB_URL: got %q", e.Value)
	}
	if e, _ := em.Get("SECRET"); e.Value != "super secret" {
		t.Errorf("SECRET: got %q", e.Value)
	}
}

func TestParse_CommentsAndBlanks(t *testing.T) {
	input := `# This is a comment

APP_ENV=production # inline comment
`
	em, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(em.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(em.Entries))
	}
	e, _ := em.Get("APP_ENV")
	if e.Value != "production" {
		t.Errorf("APP_ENV value: got %q", e.Value)
	}
	if e.Comment != "inline comment" {
		t.Errorf("APP_ENV comment: got %q", e.Comment)
	}
}

func TestParse_DuplicateKey(t *testing.T) {
	input := `FOO=bar
FOO=baz
`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestParse_InvalidLine(t *testing.T) {
	input := `NODEQUALS
`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParse_EmptyValue(t *testing.T) {
	input := `EMPTY=
`
	em, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e, ok := em.Get("EMPTY"); !ok || e.Value != "" {
		t.Errorf("EMPTY: got %q", e.Value)
	}
}
