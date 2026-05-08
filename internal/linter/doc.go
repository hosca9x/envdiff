// Package linter implements a rule-based linter for .env file contents.
//
// It analyses a parsed env map (map[string]string) and reports violations
// of common best-practice rules:
//
//   - Keys must be UPPER_SNAKE_CASE (no lowercase letters).
//   - Values must not exceed 512 characters.
//   - URL values (http:// or https://) should be quoted.
//   - Values must not contain leading or trailing whitespace.
//
// Usage:
//
//	result := linter.Lint(env)
//	if result.HasIssues() {
//	    for _, issue := range result.Issues {
//	        fmt.Printf("%s: %s\n", issue.Key, issue.Message)
//	    }
//	}
//
// The linter is non-destructive: it never modifies the input map.
package linter
