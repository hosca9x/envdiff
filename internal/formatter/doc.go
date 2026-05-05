// Package formatter provides output formatting for envdiff results.
//
// It supports multiple output formats for displaying differences between
// .env files. The available formats are:
//
//   - text: human-readable plain text output with +/-/~ prefixes
//   - json: machine-readable JSON array output
//
// Example usage:
//
//	f := formatter.New(formatter.FormatText, os.Stdout)
//	if err := f.Write(diffs); err != nil {
//		log.Fatal(err)
//	}
package formatter
