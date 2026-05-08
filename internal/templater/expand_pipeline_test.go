package templater

import (
	"testing"
)

func TestPipeline_SingleLayer(t *testing.T) {
	layers := []map[string]string{
		{"HOST": "localhost", "DSN": "postgres://${HOST}/mydb"},
	}
	out, err := Pipeline(layers, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost/mydb" {
		t.Errorf("got %q", out["DSN"])
	}
}

func TestPipeline_CrossLayerReference(t *testing.T) {
	layers := []map[string]string{
		{"BASE": "https://api.example.com"},
		{"ENDPOINT": "${BASE}/v1/users"},
	}
	out, err := Pipeline(layers, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ENDPOINT"] != "https://api.example.com/v1/users" {
		t.Errorf("got %q", out["ENDPOINT"])
	}
	if out["BASE"] != "https://api.example.com" {
		t.Errorf("BASE should be preserved, got %q", out["BASE"])
	}
}

func TestPipeline_LaterLayerWinsOnConflict(t *testing.T) {
	layers := []map[string]string{
		{"PORT": "3000"},
		{"PORT": "8080"},
	}
	out, err := Pipeline(layers, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected 8080, got %q", out["PORT"])
	}
}

func TestPipeline_StrictMissingRef(t *testing.T) {
	layers := []map[string]string{
		{"URL": "${UNDEFINED_VAR}/path"},
	}
	_, err := Pipeline(layers, Options{Strict: true})
	if err == nil {
		t.Fatal("expected error for missing reference")
	}
}

func TestPipeline_EmptyLayers(t *testing.T) {
	out, err := Pipeline([]map[string]string{}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty result, got %v", out)
	}
}

func TestPipeline_ChainedReferences(t *testing.T) {
	// Verify that variables resolved in earlier layers can be referenced
	// transitively by variables defined in later layers.
	layers := []map[string]string{
		{"SCHEME": "https", "HOST": "example.com"},
		{"BASE": "${SCHEME}://${HOST}"},
		{"URL": "${BASE}/api"},
	}
	out, err := Pipeline(layers, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "https://example.com/api" {
		t.Errorf("expected \"https://example.com/api\", got %q", out["URL"])
	}
	if out["BASE"] != "https://example.com" {
		t.Errorf("expected \"https://example.com\", got %q", out["BASE"])
	}
}
