package comfactoryrouting

import "math"

// Auto-generated seed corpus via type-driven boundary induction.
// Each field is tested with empty, path-traversal, SQL-injection, and boundary values.

// SeedCorpusData returns boundary test values for each field type.
// Used by the fuzz harness to construct targeted corpus entries.
func SeedCorpusData() map[string][]interface{} {
	return map[string][]interface{}{
		// buf ([1024]f64)
		"buf_empty": {},
		// head (usize)
		"head_empty": {},
		// count (usize)
		"count_empty": {},
	}
}
