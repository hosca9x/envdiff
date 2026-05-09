// Package scorer evaluates the quality of an env map and produces a
// numeric score from 0 to 100.
//
// The score is computed across four equally-weighted criteria:
//
//   - Key Format (25 pts): all keys should be UPPER_SNAKE_CASE.
//   - Value Completeness (25 pts): no key should have an empty value.
//   - Secret Hygiene (25 pts): keys matching sensitive patterns (e.g.
//     PASSWORD, TOKEN, SECRET) must have values of sufficient length.
//   - Consistency (25 pts): reserved for future structural checks;
//     currently awarded in full.
//
// Usage:
//
//	env := map[string]string{
//		"APP_NAME":   "myapp",
//		"API_SECRET": "supersecretvalue",
//	}
//	score := scorer.Evaluate(env)
//	fmt.Println(score.Total)      // e.g. 100
//	fmt.Println(score.Penalties)  // []
package scorer
