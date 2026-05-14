// Package deduplicator resolves duplicate keys that arise when merging
// multiple .env maps together.
//
// When loading env files from several sources (e.g. a base file, an
// environment-specific override, and a local developer file), the same key
// may appear more than once. The deduplicator provides three strategies for
// handling such conflicts:
//
//   - FirstWins  – the value from the earliest map is kept.
//   - LastWins   – the value from the latest map wins.
//   - ErrorOnDuplicate – processing stops and an error is returned as soon
//     as a duplicate key is detected.
//
// The Result type carries both the final merged map and a Duplicates index
// that records every value seen for each conflicting key, making it easy for
// callers to audit or report on collisions.
//
// Example:
//
//	result, err := deduplicator.Deduplicate(
//		deduplicator.LastWins,
//		baseEnv,
//		overrideEnv,
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result.Env)
package deduplicator
