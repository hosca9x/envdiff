package differ

// OnlyChanged returns a new slice containing only entries that represent
// a change (Added, Removed, or Changed).
func OnlyChanged(entries []Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if e.IsChanged() {
			out = append(out, e)
		}
	}
	return out
}

// ByType returns a new slice containing only entries matching the given
// ChangeType values.
func ByType(entries []Entry, types ...ChangeType) []Entry {
	set := make(map[ChangeType]struct{}, len(types))
	for _, t := range types {
		set[t] = struct{}{}
	}
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if _, ok := set[e.ChangeType]; ok {
			out = append(out, e)
		}
	}
	return out
}

// Keys returns the keys from the given entries as a string slice.
func Keys(entries []Entry) []string {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	return keys
}
