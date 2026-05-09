// Package scorer provides a quality scoring mechanism for .env files,
// rating them based on key formatting, value completeness, secret hygiene,
// and structural consistency.
package scorer

import (
	"strings"
)

// Score holds the result of evaluating an env map.
type Score struct {
	// Total is the overall quality score out of 100.
	Total int
	// Breakdown maps each criterion name to its individual score.
	Breakdown map[string]int
	// Penalties lists human-readable reasons for deductions.
	Penalties []string
}

const (
	maxScore          = 100
	emptyValuePenalty = 5
	lowercaseKeyPenalty = 3
	plainSecretPenalty  = 10
)

var sensitivePatterns = []string{"SECRET", "PASSWORD", "TOKEN", "PRIVATE", "API_KEY", "PASS"}

// Evaluate scores the given env map and returns a Score.
func Evaluate(env map[string]string) Score {
	if len(env) == 0 {
		return Score{Total: 0, Breakdown: map[string]int{}, Penalties: []string{"empty env map"}}
	}

	breakdown := map[string]int{
		"key_format":    25,
		"value_complete": 25,
		"secret_hygiene": 25,
		"consistency":   25,
	}
	var penalties []string

	// Check key formatting (uppercase + underscore only)
	for k := range env {
		if k != strings.ToUpper(k) {
			breakdown["key_format"] -= lowercaseKeyPenalty
			penalties = append(penalties, "lowercase key: "+k)
		}
	}
	if breakdown["key_format"] < 0 {
		breakdown["key_format"] = 0
	}

	// Check value completeness
	for k, v := range env {
		if strings.TrimSpace(v) == "" {
			breakdown["value_complete"] -= emptyValuePenalty
			penalties = append(penalties, "empty value for key: "+k)
		}
	}
	if breakdown["value_complete"] < 0 {
		breakdown["value_complete"] = 0
	}

	// Check secret hygiene (sensitive keys should not have short/trivial values)
	for k, v := range env {
		if isSensitive(k) && len(v) < 8 {
			breakdown["secret_hygiene"] -= plainSecretPenalty
			penalties = append(penalties, "weak secret value for key: "+k)
		}
	}
	if breakdown["secret_hygiene"] < 0 {
		breakdown["secret_hygiene"] = 0
	}

	total := 0
	for _, v := range breakdown {
		total += v
	}
	if total > maxScore {
		total = maxScore
	}

	return Score{
		Total:     total,
		Breakdown: breakdown,
		Penalties: penalties,
	}
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
