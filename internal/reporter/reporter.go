// Package reporter provides summary reporting for env diff operations.
package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/yourorg/envdiff/internal/differ"
)

// Summary holds aggregated statistics about a diff operation.
type Summary struct {
	Added   int
	Removed int
	Changed int
	Total   int
}

// Reporter generates human-readable summary reports from diff results.
type Reporter struct {
	w io.Writer
}

// New creates a new Reporter that writes to w.
func New(w io.Writer) *Reporter {
	return &Reporter{w: w}
}

// Summarize computes a Summary from a slice of DiffEntries.
func Summarize(entries []differ.DiffEntry) Summary {
	s := Summary{Total: len(entries)}
	for _, e := range entries {
		switch e.Status {
		case differ.Added:
			s.Added++
		case differ.Removed:
			s.Removed++
		case differ.Changed:
			s.Changed++
		}
	}
	return s
}

// WriteSummary writes a concise summary of diff results to the reporter's writer.
func (r *Reporter) WriteSummary(entries []differ.DiffEntry) error {
	s := Summarize(entries)
	_, err := fmt.Fprintf(r.w,
		"Summary: %d total | +%d added | -%d removed | ~%d changed\n",
		s.Total, s.Added, s.Removed, s.Changed,
	)
	return err
}

// WriteKeyList writes an alphabetically sorted list of keys grouped by status.
func (r *Reporter) WriteKeyList(entries []differ.DiffEntry) error {
	groups := map[differ.Status][]string{}
	for _, e := range entries {
		groups[e.Status] = append(groups[e.Status], e.Key)
	}

	order := []differ.Status{differ.Added, differ.Removed, differ.Changed}
	labels := map[differ.Status]string{
		differ.Added:   "ADDED",
		differ.Removed: "REMOVED",
		differ.Changed: "CHANGED",
	}

	for _, status := range order {
		keys := groups[status]
		if len(keys) == 0 {
			continue
		}
		sort.Strings(keys)
		for _, k := range keys {
			if _, err := fmt.Fprintf(r.w, "[%s] %s\n", labels[status], k); err != nil {
				return err
			}
		}
	}
	return nil
}
