// Package comparator provides utilities for comparing two env maps
// and producing a structured similarity report.
package comparator

import "github.com/yourorg/envdiff/internal/differ"

// Result holds the outcome of comparing two env maps.
type Result struct {
	// MatchingKeys are keys present in both maps with identical values.
	MatchingKeys []string
	// DriftedKeys are keys present in both maps but with different values.
	DriftedKeys []string
	// OnlyInSource are keys present only in the source map.
	OnlyInSource []string
	// OnlyInTarget are keys present only in the target map.
	OnlyInTarget []string
	// SimilarityScore is a value in [0.0, 1.0] representing how similar the maps are.
	SimilarityScore float64
}

// Compare analyses source and target env maps and returns a Result.
// It uses the differ package internally to compute the diff entries.
func Compare(source, target map[string]string) Result {
	entries := differ.Diff(source, target)

	var matching, drifted, onlySource, onlyTarget []string

	for _, e := range entries {
		switch e.Type {
		case differ.Unchanged:
			matching = append(matching, e.Key)
		case differ.Changed:
			drifted = append(drifted, e.Key)
		case differ.Added:
			onlyTarget = append(onlyTarget, e.Key)
		case differ.Removed:
			onlySource = append(onlySource, e.Key)
		}
	}

	total := len(entries)
	var score float64
	if total > 0 {
		score = float64(len(matching)) / float64(total)
	}

	return Result{
		MatchingKeys:    matching,
		DriftedKeys:     drifted,
		OnlyInSource:    onlySource,
		OnlyInTarget:    onlyTarget,
		SimilarityScore: score,
	}
}
