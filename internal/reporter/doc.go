// Package reporter provides summary and listing reports for the results
// of env file diff operations.
//
// It consumes []differ.DiffEntry values produced by the differ package and
// formats them into human-readable output suitable for CLI display or logs.
//
// # Usage
//
//	r := reporter.New(os.Stdout)
//
//	// Print a one-line statistical summary
//	r.WriteSummary(entries)
//
//	// Print keys grouped and sorted by change status
//	r.WriteKeyList(entries)
//
// Summarize can also be used independently to obtain a Summary struct
// without any I/O side effects:
//
//	s := reporter.Summarize(entries)
//	fmt.Println(s.Added, s.Removed, s.Changed)
package reporter
