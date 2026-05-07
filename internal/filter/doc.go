// Package filter provides key-based filtering for env maps.
//
// It supports three filtering strategies that can be combined:
//
//   - WithPrefixes: keep (or exclude) keys that start with any given prefix.
//   - WithKeys: keep (or exclude) an explicit set of keys.
//   - Exclude: invert the filter so matched keys are removed.
//
// Example usage:
//
//	// Keep only APP_ and DB_ keys
//	subset := filter.Filter(env, filter.WithPrefixes("APP_", "DB_"))
//
//	// Remove all SECRET_ keys
//	redacted := filter.Filter(env, filter.WithPrefixes("SECRET_"), filter.Exclude())
//
// Filter never mutates the input map; it always returns a new map.
package filter
