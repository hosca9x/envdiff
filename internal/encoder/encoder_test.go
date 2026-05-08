package encoder

import (
	"testing"
)

func TestEncode_Base64_AllKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASS": "secret",
		"APP_ENV": "production",
	}
	enc := New(FormatBase64)
	out, err := enc.Encode(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "c2VjcmV0" {
		t.Errorf("expected c2VjcmV0, got %s", out["DB_PASS"])
	}
	if out["APP_ENV"] != "cHJvZHVjdGlvbg==" {
		t.Errorf("expected cHJvZHVjdGlvbg==, got %s", out["APP_ENV"])
	}
}

func TestDecode_Base64_RoundTrip(t *testing.T) {
	original := map[string]string{
		"TOKEN": "my-secret-token",
		"HOST":  "localhost",
	}
	enc := New(FormatBase64)
	encoded, err := enc.Encode(original)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	decoded, err := enc.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	for k, v := range original {
		if decoded[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, decoded[k])
		}
	}
}

func TestEncode_Hex_AllKeys(t *testing.T) {
	env := map[string]string{"KEY": "hi"}
	enc := New(FormatHex)
	out, err := enc.Encode(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 'h'=0x68, 'i'=0x69
	if out["KEY"] != "6869" {
		t.Errorf("expected 6869, got %s", out["KEY"])
	}
}

func TestDecode_Hex_RoundTrip(t *testing.T) {
	original := map[string]string{"KEY": "hello"}
	enc := New(FormatHex)
	encoded, _ := enc.Encode(original)
	decoded, err := enc.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if decoded["KEY"] != "hello" {
		t.Errorf("expected hello, got %s", decoded["KEY"])
	}
}

func TestEncode_WithKeys_OnlyTargetedKeys(t *testing.T) {
	env := map[string]string{
		"SECRET": "pass",
		"PUBLIC": "open",
	}
	enc := New(FormatBase64).WithKeys("SECRET")
	out, err := enc.Encode(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET"] == "pass" {
		t.Error("expected SECRET to be encoded")
	}
	if out["PUBLIC"] != "open" {
		t.Errorf("expected PUBLIC to remain unchanged, got %s", out["PUBLIC"])
	}
}

func TestEncode_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	enc := New(FormatBase64)
	_, err := enc.Encode(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "value" {
		t.Error("original map was mutated")
	}
}

func TestDecode_InvalidBase64_ReturnsError(t *testing.T) {
	env := map[string]string{"KEY": "!!!not-base64!!!"}
	enc := New(FormatBase64)
	_, err := enc.Decode(env)
	if err == nil {
		t.Error("expected error for invalid base64, got nil")
	}
}

func TestEncode_Raw_Passthrough(t *testing.T) {
	env := map[string]string{"KEY": "unchanged"}
	enc := New(FormatRaw)
	out, err := enc.Encode(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "unchanged" {
		t.Errorf("expected unchanged, got %s", out["KEY"])
	}
}
