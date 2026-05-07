package templater

// Pipeline applies a sequence of environment maps in order, expanding each
// layer's values using the merged context of all previously resolved layers.
// Later layers may reference variables defined in earlier ones.
//
// The returned map is a flat merge of all layers after expansion; in case of
// key conflicts the last layer wins.
func Pipeline(layers []map[string]string, opts Options) (map[string]string, error) {
	merged := make(map[string]string)

	for _, layer := range layers {
		// Build the expansion context: already-resolved keys + current layer.
		ctx := make(map[string]string, len(merged)+len(layer))
		for k, v := range merged {
			ctx[k] = v
		}
		for k, v := range layer {
			ctx[k] = v
		}

		expanded, err := Expand(ctx, opts)
		if err != nil {
			return nil, err
		}

		// Promote only the keys from this layer into merged.
		for k := range layer {
			merged[k] = expanded[k]
		}
	}

	return merged, nil
}
