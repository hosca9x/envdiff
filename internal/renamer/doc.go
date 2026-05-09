// Package renamer provides key-renaming functionality for env maps.
//
// It applies a list of rename rules (From -> To) to a given environment map,
// producing a new map with the renamed keys. The original map is never mutated.
//
// Behaviour:
//
//   - If a source key (From) does not exist, it is recorded as skipped.
//   - If the target key (To) already exists and overwrite is false, the rule
//     is recorded as a conflict and skipped.
//   - If overwrite is true, an existing target key is replaced.
//   - Rules where From == To are silently ignored.
//
// Example:
//
//	rules := []renamer.Rule{
//		{From: "DB_HOST", To: "DATABASE_HOST"},
//		{From: "DB_PORT", To: "DATABASE_PORT"},
//	}
//	out, result, err := renamer.Rename(env, rules, false)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result.Renamed) // map[DB_HOST:DATABASE_HOST DB_PORT:DATABASE_PORT]
package renamer
