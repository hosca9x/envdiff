// Package encoder provides utilities for encoding and decoding values
// within an env map using configurable formats such as base64 or hex.
//
// This is useful when storing encoded secrets in .env files or when
// preparing env values for transmission across systems that require
// safe ASCII representations.
//
// # Formats
//
// FormatBase64 encodes values using standard base64 encoding.
// FormatHex encodes values as lowercase hexadecimal strings.
// FormatRaw passes values through unchanged (identity transform).
//
// # Selective Encoding
//
// By default all keys in the map are encoded. Use WithKeys to restrict
// encoding to a specific subset of keys, leaving all others untouched.
//
// # Immutability
//
// Encode and Decode always return a new map; the original is never mutated.
//
// # Example
//
//	enc := encoder.New(encoder.FormatBase64).WithKeys("DB_PASSWORD", "API_KEY")
//	encoded, err := enc.Encode(env)
package encoder
