// Package validator inspects parsed .env maps for structural and semantic
// issues before they are diffed, reconciled, or written to disk.
//
// Checks performed
//
// Error-level checks (block further processing when strict mode is on):
//
//	- Key format: keys must match [A-Za-z_][A-Za-z0-9_]* to be portable
//	  across shells and CI systems.
//
// Warning-level checks (reported but non-blocking):
//
//	- Empty values: a key present with no value is often a misconfiguration.
//	- Unresolved placeholders: values containing ${ or % may indicate that
//	  variable substitution was expected but did not occur.
//
// Usage
//
//	issues := validator.Validate(env)
//	if validator.HasErrors(issues) {
//	    // handle blocking errors
//	}
//	for _, issue := range issues {
//	    fmt.Println(issue)
//	}
package validator
