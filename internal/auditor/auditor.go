// Package auditor provides change auditing for .env file diffs,
// recording who changed what and when with structured audit log entries.
package auditor

import (
	"fmt"
	"time"

	"github.com/user/envdiff/internal/differ"
)

// Entry represents a single audited change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Actor     string    `json:"actor"`
	Key       string    `json:"key"`
	ChangeType string   `json:"change_type"`
	OldValue  string    `json:"old_value,omitempty"`
	NewValue  string    `json:"new_value,omitempty"`
	Environment string  `json:"environment,omitempty"`
}

// Auditor records diff entries into a structured audit log.
type Auditor struct {
	actor       string
	environment string
	now         func() time.Time
}

// New returns an Auditor for the given actor and environment label.
func New(actor, environment string) *Auditor {
	return &Auditor{
		actor:       actor,
		environment: environment,
		now:         time.Now,
	}
}

// Audit converts a slice of differ.Entry values into audit log entries.
// Sensitive values should be masked before calling Audit.
func (a *Auditor) Audit(diffs []differ.Entry) []Entry {
	entries := make([]Entry, 0, len(diffs))
	for _, d := range diffs {
		e := Entry{
			Timestamp:   a.now(),
			Actor:       a.actor,
			Key:         d.Key,
			ChangeType:  changeLabel(d.Type),
			Environment: a.environment,
		}
		switch d.Type {
		case differ.Added:
			e.NewValue = d.NewValue
		case differ.Removed:
			e.OldValue = d.OldValue
		case differ.Changed:
			e.OldValue = d.OldValue
			e.NewValue = d.NewValue
		}
		entries = append(entries, e)
	}
	return entries
}

// Summary returns a human-readable summary line for a single audit Entry.
func Summary(e Entry) string {
	ts := e.Timestamp.UTC().Format(time.RFC3339)
	switch e.ChangeType {
	case "added":
		return fmt.Sprintf("[%s] %s added %q in %s (value: %q)", ts, e.Actor, e.Key, e.Environment, e.NewValue)
	case "removed":
		return fmt.Sprintf("[%s] %s removed %q in %s (was: %q)", ts, e.Actor, e.Key, e.Environment, e.OldValue)
	case "changed":
		return fmt.Sprintf("[%s] %s changed %q in %s (%q -> %q)", ts, e.Actor, e.Key, e.Environment, e.OldValue, e.NewValue)
	}
	return fmt.Sprintf("[%s] %s touched %q in %s", ts, e.Actor, e.Key, e.Environment)
}

// Summaries returns a human-readable summary line for each audit Entry.
// It is a convenience wrapper around Summary for processing a batch of entries.
func Summaries(entries []Entry) []string {
	lines := make([]string, len(entries))
	for i, e := range entries {
		lines[i] = Summary(e)
	}
	return lines
}

func changeLabel(t differ.ChangeType) string {
	switch t {
	case differ.Added:
		return "added"
	case differ.Removed:
		return "removed"
	case differ.Changed:
		return "changed"
	default:
		return "unchanged"
	}
}
