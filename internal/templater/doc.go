// Package templater resolves variable references embedded inside .env values.
//
// It supports two common placeholder syntaxes:
//
//	${VAR_NAME}   — curly-brace delimited reference
//	$VAR_NAME     — bare dollar-sign reference
//
// References are resolved against the same environment map being expanded,
// enabling self-referential composition of values:
//
//	BASE_URL=https://api.example.com
//	HEALTH_URL=${BASE_URL}/health
//
// # Usage
//
//	out, err := templater.Expand(env, templater.Options{Strict: true})
//
// When Strict is true, Expand returns an *ErrMissingVar error if a referenced
// variable is absent from the map. When Strict is false (the default), missing
// references are replaced with Options.Fallback (empty string by default).
package templater
