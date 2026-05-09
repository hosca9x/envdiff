// Package normalizer provides utilities for normalising the keys of an env map
// to a consistent format before diffing, exporting, or reconciling.
//
// Three strategies are supported:
//
//   - UpperSnake – converts camelCase, kebab-case, and dot-separated keys to
//     UPPER_SNAKE_CASE, which is the conventional format for environment
//     variables (e.g. "dbHost" → "DB_HOST").
//
//   - LowerSnake – same transformation as UpperSnake but the result is
//     lower-cased (e.g. "DB-HOST" → "db_host").
//
//   - TrimPrefix – removes a fixed string prefix from every key, useful when
//     merging env files that namespace keys differently
//     (e.g. "APP_HOST" → "HOST" with prefix "APP_").
//
// The original map is never modified; Normalize always returns a fresh copy.
package normalizer
