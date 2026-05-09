// Package trimmer provides utilities for stripping unwanted whitespace from
// the values of an env map.
//
// Four strategies are available:
//
//   - TrimAll      – remove leading and trailing whitespace (default)
//   - TrimLeading  – remove only leading whitespace
//   - TrimTrailing – remove only trailing whitespace
//   - TrimInner    – collapse internal whitespace runs to a single space
//
// Basic usage:
//
//	env := map[string]string{"API_URL": "  https://example.com  "}
//	cleaned, changes, err := trimmer.Trim(env, trimmer.TrimAll)
//
// Strategy names can be parsed from strings (e.g. CLI flags):
//
//	st, err := trimmer.ParseStrategy("trailing")
package trimmer
