// Package profiler provides statistical analysis of env variable maps.
//
// It computes a Profile containing:
//   - Total key count
//   - Empty value count
//   - Sensitive key count (based on common naming patterns)
//   - Prefix group distribution
//   - Average value length
//   - Longest and shortest key names
//   - Prefixes shared by multiple keys
//
// Usage:
//
//	env := map[string]string{
//		"APP_HOST":    "localhost",
//		"APP_PORT":    "8080",
//		"DB_PASSWORD": "secret",
//	}
//	p := profiler.Analyze(env)
//	profiler.WriteReport(os.Stdout, p)
package profiler
