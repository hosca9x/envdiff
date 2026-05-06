// Package sorter provides utilities for ordering the keys of an environment
// variable map according to a configurable strategy.
//
// Three strategies are available:
//
//   - Alphabetical: keys are sorted A→Z (default).
//   - Reverse: keys are sorted Z→A.
//   - GroupedByPrefix: keys are clustered by the portion of their name that
//     precedes the first underscore (e.g. all DB_* keys together, then all
//     APP_* keys together), with groups themselves ordered alphabetically and
//     keys within each group also sorted alphabetically.
//
// The Sort function is non-destructive — it never modifies the original map.
package sorter
