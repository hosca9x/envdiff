package encoder

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Format represents the encoding format to apply to env values.
type Format int

const (
	// FormatBase64 encodes values using standard base64 encoding.
	FormatBase64 Format = iota
	// FormatHex encodes values as lowercase hexadecimal.
	FormatHex
	// FormatRaw returns values unchanged.
	FormatRaw
)

// Encoder encodes or decodes env map values.
type Encoder struct {
	format Format
	keys   map[string]struct{} // if non-empty, only encode these keys
}

// New returns an Encoder using the given format.
func New(f Format) *Encoder {
	return &Encoder{format: f, keys: make(map[string]struct{})}
}

// WithKeys restricts encoding to the specified keys.
func (e *Encoder) WithKeys(keys ...string) *Encoder {
	for _, k := range keys {
		e.keys[k] = struct{}{}
	}
	return e
}

// Encode returns a new map with values encoded according to the configured format.
// The original map is never mutated.
func (e *Encoder) Encode(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if len(e.keys) > 0 {
			if _, ok := e.keys[k]; !ok {
				out[k] = v
				continue
			}
		}
		encoded, err := e.encodeValue(v)
		if err != nil {
			return nil, fmt.Errorf("encoder: key %q: %w", k, err)
		}
		out[k] = encoded
	}
	return out, nil
}

// Decode returns a new map with values decoded according to the configured format.
// The original map is never mutated.
func (e *Encoder) Decode(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if len(e.keys) > 0 {
			if _, ok := e.keys[k]; !ok {
				out[k] = v
				continue
			}
		}
		decoded, err := e.decodeValue(v)
		if err != nil {
			return nil, fmt.Errorf("encoder: key %q: %w", k, err)
		}
		out[k] = decoded
	}
	return out, nil
}

func (e *Encoder) encodeValue(v string) (string, error) {
	switch e.format {
	case FormatBase64:
		return base64.StdEncoding.EncodeToString([]byte(v)), nil
	case FormatHex:
		var sb strings.Builder
		for _, b := range []byte(v) {
			fmt.Fprintf(&sb, "%02x", b)
		}
		return sb.String(), nil
	default:
		return v, nil
	}
}

func (e *Encoder) decodeValue(v string) (string, error) {
	switch e.format {
	case FormatBase64:
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return "", fmt.Errorf("invalid base64: %w", err)
		}
		return string(b), nil
	case FormatHex:
		if len(v)%2 != 0 {
			return "", fmt.Errorf("invalid hex string: odd length")
		}
		bytes := make([]byte, len(v)/2)
		for i := 0; i < len(v); i += 2 {
			var b byte
			_, err := fmt.Sscanf(v[i:i+2], "%02x", &b)
			if err != nil {
				return "", fmt.Errorf("invalid hex byte at %d: %w", i, err)
			}
			bytes[i/2] = b
		}
		return string(bytes), nil
	default:
		return v, nil
	}
}
