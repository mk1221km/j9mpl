package comfactoryrouting

import "testing"

// Auto-generated seed corpus via type-driven boundary induction.
// Each field is tested with empty, path-traversal, SQL-injection, and boundary values.

func init() {
	// Seeds are loaded via transaction_router_test.go's f.Add() calls.
	// This file exists to ensure seeds_test.go is a valid Go file.
}

// SeedCorpusData returns boundary test values for each field type.
// Used by the fuzz harness to construct targeted corpus entries.
func SeedCorpusData() map[string][]interface{} {
	return map[string][]interface{}{
		// txId (string)
		"txId_empty": {},
		"txId_path_traversal": {"../../etc/passwd"},
		"txId_sql_injection": {"'; DROP TABLE --"},
		// sender (string)
		"sender_empty": {},
		"sender_path_traversal": {"../../etc/passwd"},
		"sender_sql_injection": {"'; DROP TABLE --"},
		// receiver (string)
		"receiver_empty": {},
		"receiver_path_traversal": {"../../etc/passwd"},
		"receiver_sql_injection": {"'; DROP TABLE --"},
		// amount (float64)
		"amount_empty": {},
		"amount_negative": {-1.0},
		"amount_nan": {math.NaN()},
		// priority (string)
		"priority_empty": {},
		"priority_path_traversal": {"../../etc/passwd"},
		"priority_sql_injection": {"'; DROP TABLE --"},
	}
}
