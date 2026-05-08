// Package transformer provides utilities for applying value transformations
// to env maps.
//
// A Transformer chains one or more TransformFunc functions, each of which
// receives a key-value pair and returns a (possibly modified) value or an error.
//
// Built-in transforms:
//
//	- TrimSpace  — strips leading/trailing whitespace from values
//	- ToUpper    — uppercases all values
//	- ToLower    — lowercases all values
//	- ReplaceValue — replaces a substring in every value
//
// Example:
//
//	tr := transformer.New(
//		transformer.TrimSpace(),
//		transformer.ReplaceValue("localhost", "prod-db"),
//	)
//	result, err := tr.Apply(env)
package transformer
