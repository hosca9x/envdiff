package renamer

import "fmt"

// Rule defines a single key rename mapping.
type Rule struct {
	From string
	To   string
}

// Result holds the outcome of a rename operation.
type Result struct {
	Renamed  map[string]string // old key -> new key
	Skipped  []string          // keys not found in env
	Conflict []string          // target keys that already exist
}

// Rename applies the given rules to env, returning a new map with keys renamed.
// Keys not found are recorded in Result.Skipped.
// If a target key already exists in env, the rule is recorded in Result.Conflict
// and skipped unless overwrite is true.
func Rename(env map[string]string, rules []Rule, overwrite bool) (map[string]string, Result, error) {
	if env == nil {
		return nil, Result{}, fmt.Errorf("renamer: env must not be nil")
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	res := Result{
		Renamed: make(map[string]string),
	}

	for _, rule := range rules {
		if rule.From == rule.To {
			continue
		}

		val, exists := out[rule.From]
		if !exists {
			res.Skipped = append(res.Skipped, rule.From)
			continue
		}

		if _, conflict := out[rule.To]; conflict && !overwrite {
			res.Conflict = append(res.Conflict, rule.To)
			continue
		}

		delete(out, rule.From)
		out[rule.To] = val
		res.Renamed[rule.From] = rule.To
	}

	return out, res, nil
}

// MustRename is like Rename but panics on error.
func MustRename(env map[string]string, rules []Rule, overwrite bool) (map[string]string, Result) {
	out, res, err := Rename(env, rules, overwrite)
	if err != nil {
		panic(err)
	}
	return out, res
}
