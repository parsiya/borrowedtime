package shared

import "sort"

// map (mainly map[string]string) utils.

// SortedKeys returns a sorted slice of keys for a map[string]string.
func SortedKeys(mp map[string]string) []string {
	// Create a slice equivalent to the size of mp.
	keys := make([]string, len(mp))
	// Might be a bit faster than append.
	i := 0
	for key := range mp {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}
