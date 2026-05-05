// Package reconciler provides functionality to reconcile differences
// between two .env files by applying changes from a source environment
// to a target environment.
package reconciler

import (
	"fmt"

	"github.com/user/envdiff/internal/differ"
)

// Strategy defines how conflicts are resolved during reconciliation.
type Strategy int

const (
	// StrategySourceWins overwrites target values with source values on conflict.
	StrategySourceWins Strategy = iota
	// StrategyTargetWins keeps target values on conflict.
	StrategyTargetWins
	// StrategyAddOnly only adds missing keys, never modifies existing ones.
	StrategyAddOnly
)

// Result holds the outcome of a reconciliation operation.
type Result struct {
	// Merged is the final reconciled map.
	Merged map[string]string
	// Applied is the list of diff entries that were applied.
	Applied []differ.Entry
	// Skipped is the list of diff entries that were skipped.
	Skipped []differ.Entry
}

// Reconcile merges source into target according to the given strategy.
// It returns a Result describing what was applied and what was skipped.
func Reconcile(target, source map[string]string, strategy Strategy) (*Result, error) {
	if target == nil {
		return nil, fmt.Errorf("reconciler: target map must not be nil")
	}
	if source == nil {
		return nil, fmt.Errorf("reconciler: source map must not be nil")
	}

	diffs := differ.Diff(target, source)

	// Copy target so we do not mutate the original.
	merged := make(map[string]string, len(target))
	for k, v := range target {
		merged[k] = v
	}

	var applied, skipped []differ.Entry

	for _, entry := range diffs {
		switch entry.Type {
		case differ.Added:
			// Key exists in source but not in target — always add.
			merged[entry.Key] = entry.NewValue
			applied = append(applied, entry)

		case differ.Removed:
			// Key exists in target but not in source — skip removal by default.
			skipped = append(skipped, entry)

		case differ.Changed:
			switch strategy {
			case StrategySourceWins:
				merged[entry.Key] = entry.NewValue
				applied = append(applied, entry)
			case StrategyTargetWins, StrategyAddOnly:
				skipped = append(skipped, entry)
			}
		}
	}

	return &Result{
		Merged:  merged,
		Applied: applied,
		Skipped: skipped,
	}, nil
}
