// Package exporter writes env maps to various output formats.
//
// Supported formats:
//
//   - FormatEnv    — standard KEY=VALUE lines (default)
//   - FormatExport — shell-compatible `export KEY=VALUE` lines
//   - FormatDocker — Docker --env-file compatible KEY=VALUE lines
//
// Values can optionally be double-quoted, with internal quotes escaped.
// Output keys are always written in sorted order for deterministic results.
//
// Example:
//
//	e := exporter.NewWithOptions(exporter.FormatExport, true)
//	err := e.WriteFile(".env.production", myEnvMap)
package exporter
