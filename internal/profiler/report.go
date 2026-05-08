package profiler

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// WriteReport writes a human-readable profile report to w.
func WriteReport(w io.Writer, p Profile) error {
	lines := []string{
		"=== Env Profile Report ===",
		fmt.Sprintf("Total Keys       : %d", p.TotalKeys),
		fmt.Sprintf("Empty Values     : %d", p.EmptyValues),
		fmt.Sprintf("Sensitive Keys   : %d", p.SensitiveKeys),
		fmt.Sprintf("Avg Value Length : %.2f", p.AvgValueLength),
		fmt.Sprintf("Longest Key      : %s", p.LongestKey),
		fmt.Sprintf("Shortest Key     : %s", p.ShortestKey),
	}

	if len(p.PrefixGroups) > 0 {
		lines = append(lines, "\nPrefix Groups:")
		prefixes := make([]string, 0, len(p.PrefixGroups))
		for prefix := range p.PrefixGroups {
			prefixes = append(prefixes, prefix)
		}
		sort.Strings(prefixes)
		for _, prefix := range prefixes {
			lines = append(lines, fmt.Sprintf("  %-20s %d keys", prefix, p.PrefixGroups[prefix]))
		}
	}

	if len(p.DuplicatePrefixes) > 0 {
		lines = append(lines, fmt.Sprintf("\nShared Prefixes  : %s", strings.Join(p.DuplicatePrefixes, ", ")))
	}

	_, err := fmt.Fprintln(w, strings.Join(lines, "\n"))
	return err
}
