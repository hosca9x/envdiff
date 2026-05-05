// Package formatter provides output formatting for env diff results.
package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Format represents the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Formatter writes diff results to an output writer.
type Formatter struct {
	format Format
	w      io.Writer
}

// New creates a new Formatter with the given format and writer.
func New(format Format, w io.Writer) *Formatter {
	return &Formatter{format: format, w: w}
}

// Write outputs the diff results according to the configured format.
func (f *Formatter) Write(diffs []differ.DiffEntry) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(diffs)
	default:
		return f.writeText(diffs)
	}
}

func (f *Formatter) writeText(diffs []differ.DiffEntry) error {
	if len(diffs) == 0 {
		_, err := fmt.Fprintln(f.w, "No differences found.")
		return err
	}

	sorted := make([]differ.DiffEntry, len(diffs))
	copy(sorted, diffs)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, d := range sorted {
		var line string
		switch d.Type {
		case differ.Added:
			line = fmt.Sprintf("+ %s=%s", d.Key, d.NewValue)
		case differ.Removed:
			line = fmt.Sprintf("- %s=%s", d.Key, d.OldValue)
		case differ.Changed:
			line = fmt.Sprintf("~ %s: %s -> %s", d.Key, d.OldValue, d.NewValue)
		}
		if _, err := fmt.Fprintln(f.w, line); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) writeJSON(diffs []differ.DiffEntry) error {
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, d := range diffs {
		sb.WriteString(fmt.Sprintf(
			"  {\"key\": %q, \"type\": %q, \"old\": %q, \"new\": %q}",
			d.Key, d.Type, d.OldValue, d.NewValue,
		))
		if i < len(diffs)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	_, err := fmt.Fprint(f.w, sb.String())
	return err
}
