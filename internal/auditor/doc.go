// Package auditor provides structured audit logging for .env file changes.
//
// It consumes a slice of [differ.Entry] values produced by the differ package
// and converts them into [Entry] records that capture who made a change,
// when it happened, which environment was affected, and what the before/after
// values were.
//
// Typical usage:
//
//	a := auditor.New("deploy-bot", "production")
//	log := a.Audit(diffs)
//	for _, entry := range log {
//		fmt.Println(auditor.Summary(entry))
//	}
//
// Values should be masked or redacted (via the masker or redactor packages)
// before being passed to Audit to avoid storing secrets in the audit log.
package auditor
