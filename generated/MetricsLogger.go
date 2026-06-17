package comfactorymetrics

import ("database/sql"; "fmt"; "strings"; _ "github.com/mattn/go-sqlite3")

type MetricRecord struct {
	Timestamp string
	MetricName string
	MetricValue float64
}

type MetricsLogger struct {
	db *sql.DB
}

func (s *MetricsLogger) initDatabase(dbPath string) error {
	// Input validation
	if strings.TrimSpace(dbPath) == "" {
		return fmt.Errorf("invalid input: dbPath is empty")
	}
	if strings.Contains(dbPath, "..") || strings.Contains(dbPath, "/etc/") {
		return fmt.Errorf("path traversal blocked")
	}
	if strings.ContainsAny(dbPath, "'\";--") {
		return fmt.Errorf("sql injection blocked")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS system_metrics (
		timestamp TEXT,
		name TEXT,
		value REAL
	)`)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	return nil
}

func (s *MetricsLogger) logMetric(dbPath string, record MetricRecord) (int64, error) {
    // Validate dbPath
    if strings.TrimSpace(dbPath) == "" {
        return 0, fmt.Errorf("invalid input: dbPath is empty")
    }
    if strings.Contains(dbPath, "..") || strings.Contains(dbPath, "/etc/") {
        return 0, fmt.Errorf("path traversal blocked")
    }
    if containsSQLInjectionPatterns(dbPath) {
        return 0, fmt.Errorf("sql injection blocked")
    }

    // Validate record fields
    if strings.TrimSpace(record.Timestamp) == "" {
        return 0, fmt.Errorf("invalid input: Timestamp is empty")
    }
    if strings.TrimSpace(record.MetricName) == "" {
        return 0, fmt.Errorf("invalid input: MetricName is empty")
    }
    if containsSQLInjectionPatterns(record.Timestamp) {
        return 0, fmt.Errorf("sql injection blocked")
    }
    if containsSQLInjectionPatterns(record.MetricName) {
        return 0, fmt.Errorf("sql injection blocked")
    }

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return 0, fmt.Errorf("open database: %w", err)
    }
    defer db.Close()

    result, err := db.Exec(
        "INSERT INTO system_metrics (timestamp, name, value) VALUES (?, ?, ?)",
        record.Timestamp, record.MetricName, record.MetricValue,
    )
    if err != nil {
        return 0, fmt.Errorf("insert metric: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return 0, fmt.Errorf("get rows affected: %w", err)
    }

    return rowsAffected, nil
}

// containsSQLInjectionPatterns checks for common SQL injection patterns.
func containsSQLInjectionPatterns(s string) bool {
    patterns := []string{"'", ";", "--", "/*", "*/", "1=1", " OR ", " AND "}
    for _, p := range patterns {
        if strings.Contains(strings.ToUpper(s), strings.ToUpper(p)) {
            return true
        }
    }
    return false
}

func (s *MetricsLogger) getAverageMetric(dbPath string, name string) (float64, error) {
    // Validate name
    if strings.TrimSpace(name) == "" {
        return 0, fmt.Errorf("invalid input: name is empty")
    }

    // Validate dbPath against path traversal
    if strings.Contains(dbPath, "..") || strings.Contains(dbPath, "/etc/") {
        return 0, fmt.Errorf("path traversal blocked")
    }

    // SQL injection pattern check (applied to both name and dbPath)
    injectionPatterns := []string{";", "'", "--", "/*", "*/"}
    for _, pattern := range injectionPatterns {
        if strings.Contains(name, pattern) {
            return 0, fmt.Errorf("sql injection blocked")
        }
        if strings.Contains(dbPath, pattern) {
            return 0, fmt.Errorf("sql injection blocked")
        }
    }

    // Open database
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return 0, fmt.Errorf("open db: %w", err)
    }
    defer db.Close()

    // Query average
    var avg sql.NullFloat64
    query := "SELECT AVG(value) FROM system_metrics WHERE name = ?"
    if err := db.QueryRow(query, name).Scan(&avg); err != nil {
        return 0, fmt.Errorf("query average metric: %w", err)
    }

    if avg.Valid {
        return avg.Float64, nil
    }
    // If no rows matched, return 0 (or could treat as error; requirements not explicit, but returning 0 is reasonable)
    return 0, nil
}
