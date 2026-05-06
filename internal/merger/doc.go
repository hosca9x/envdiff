// Package merger combines multiple parsed .env maps into a single unified
// map using a configurable conflict resolution Strategy.
//
// Three strategies are available:
//
//   - FirstWins: the first file to define a key wins; later duplicates are ignored.
//   - LastWins: the last file to define a key wins; earlier values are overwritten.
//   - ErrorOnConflict: returns a *ConflictError if the same key appears with
//     different values across files.
//
// Example usage:
//
//	base := map[string]string{"HOST": "localhost", "PORT": "5432"}
//	override := map[string]string{"HOST": "prod.example.com"}
//
//	merged, err := merger.Merge(merger.LastWins, base, override)
//	// merged["HOST"] == "prod.example.com"
//	// merged["PORT"] == "5432"
package merger
