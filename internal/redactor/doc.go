// Package redactor provides a Redactor type that identifies and replaces
// sensitive values in env maps before they are displayed, logged, or
// transmitted.
//
// A key is considered sensitive when its lowercase form contains any of the
// built-in patterns ("secret", "password", "token", "key", etc.). Additional
// patterns can be supplied via [WithPatterns].
//
// Basic usage:
//
//	r := redactor.New()
//	safe := r.Redact(env) // returns a new map; orignal is unchanged
//
// Custom placeholder and patterns:
//
//	r := redactor.New(
//		redactor.WithPlaceholder("***"),
//		redactor.WithPatterns("internal", "priv"),
//	)
package redactor
