// Package inspector provides tools for analyzing the structure and
// contents of an env map.
//
// It produces a Report containing:
//   - Total key count
//   - Count of keys with empty values
//   - Keys grouped by their prefix (portion before the first underscore)
//   - Value type hints (bool, number, url, string, empty)
//
// Example usage:
//
//	env := map[string]string{
//		"APP_HOST": "localhost",
//		"APP_PORT": "8080",
//		"DEBUG":    "true",
//	}
//	report := inspector.Inspect(env)
//	fmt.Println(report.TotalKeys)           // 3
//	fmt.Println(report.PrefixGroups["APP"]) // [APP_HOST APP_PORT]
//	fmt.Println(report.ValueHints["DEBUG"]) // bool
package inspector
