// Package grouper partitions a flat env map into named groups based on a
// configurable strategy.
//
// Two built-in strategies are provided:
//
//   - ByPrefix  — groups keys by the segment that precedes the first underscore
//     (e.g. DB_HOST and DB_PORT both land in the "DB" group). Keys without an
//     underscore are placed in the special "_other" group.
//
//   - ByLength  — groups keys into three size buckets:
//     "short" (≤ 8 chars), "medium" (≤ 20 chars), and "long" (> 20 chars).
//
// Example:
//
//	env := map[string]string{
//		"DB_HOST": "localhost",
//		"DB_PORT": "5432",
//		"APP_ENV": "production",
//	}
//	groups, err := grouper.GroupBy(env, grouper.ByPrefix)
//	// groups → [{Name:"APP", ...}, {Name:"DB", ...}]
package grouper
