package comfactoryrouting

import ("database/sql"; "fmt"; "strings"; _ "github.com/mattn/go-sqlite3")

type TransactionRecord struct {
	TxId string
	Sender string
	Receiver string
	Amount float64
	Priority string
}

type TransactionRouter struct{}

func (s *TransactionRouter) initRoutingTable(dbPath string) error {
    // Validate input
    if len(strings.TrimSpace(dbPath)) == 0 {
        return fmt.Errorf("invalid input: dbPath is empty")
    }
    if strings.Contains(dbPath, "..") || strings.Contains(dbPath, "/etc/") {
        return fmt.Errorf("path traversal blocked")
    }
    keywords := []string{"DROP", "SELECT", "INSERT", "DELETE", "UPDATE", "CREATE", "ALTER", "UNION", "--", ";"}
    for _, keyword := range keywords {
        if strings.Contains(strings.ToUpper(dbPath), keyword) {
            return fmt.Errorf("sql injection blocked")
        }
    }

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    defer db.Close()

    // Create routing_rules table
    createRulesTable := `CREATE TABLE IF NOT EXISTS routing_rules (
        min_amount REAL,
        priority TEXT,
        channel TEXT,
        PRIMARY KEY (min_amount, priority)
    )`
    if _, err := db.Exec(createRulesTable); err != nil {
        return fmt.Errorf("failed to create routing_rules table: %w", err)
    }

    // Create transaction_log table
    createLogTable := `CREATE TABLE IF NOT EXISTS transaction_log (
        tx_id TEXT PRIMARY KEY,
        sender TEXT,
        receiver TEXT,
        amount REAL,
        channel TEXT,
        status TEXT
    )`
    if _, err := db.Exec(createLogTable); err != nil {
        return fmt.Errorf("failed to create transaction_log table: %w", err)
    }

    // Insert default routing rules
    type defaultRule struct {
        minAmount float64
        priority  string
        channel   string
    }
    defaultRules := []defaultRule{
        {0.0, "low", "standard"},
        {1000.0, "high", "express"},
    }
    insertDefaultRules := `INSERT OR IGNORE INTO routing_rules (min_amount, priority, channel) VALUES (?, ?, ?)`
    for _, rule := range defaultRules {
        if _, err := db.Exec(insertDefaultRules, rule.minAmount, rule.priority, rule.channel); err != nil {
            return fmt.Errorf("failed to insert default routing rule: %w", err)
        }
    }

    return nil
}

func (s *TransactionRouter) routeTransaction(dbPath string, record TransactionRecord) (string, error) {
	if strings.TrimSpace(dbPath) == "" {
		return "", fmt.Errorf("invalid input: dbPath is empty")
	}
	if strings.TrimSpace(record.Receiver) == "" {
		return "", fmt.Errorf("invalid input: Receiver is empty")
	}
	if strings.TrimSpace(record.Priority) == "" {
		return "", fmt.Errorf("invalid input: Priority is empty")
	}
	if strings.Contains(dbPath, "..") || strings.Contains(dbPath, "/etc/") {
		return "", fmt.Errorf("path traversal blocked")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	var channel string
	err = db.QueryRow("SELECT channel FROM routing_rules WHERE priority=? AND min_amount <= ? ORDER BY min_amount DESC LIMIT 1", record.Priority, record.Amount).Scan(&channel)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no matching routing rule")
		}
		return "", fmt.Errorf("query routing rule: %w", err)
	}

	_, err = db.Exec(
		"INSERT INTO transaction_log(tx_id, sender, receiver, amount, channel, status) VALUES (?, ?, ?, ?, ?, ?)",
		record.TxId, record.Sender, record.Receiver, record.Amount, channel, "PROCESSED",
	)
	if err != nil {
		return "", fmt.Errorf("insert transaction log: %w", err)
	}

	return channel, nil
}

func (s *TransactionRouter) getTransactionCount(dbPath string, status string) (float64, error) {
	// Validate inputs
	if strings.TrimSpace(status) == "" {
		return 0, fmt.Errorf("invalid input: status is empty")
	}
	if strings.TrimSpace(dbPath) == "" {
		return 0, fmt.Errorf("invalid input: dbPath is empty")
	}
	if strings.Contains(dbPath, "..") || strings.Contains(dbPath, "/etc/") {
		return 0, fmt.Errorf("path traversal blocked")
	}
	injectionPatterns := []string{"'", ";", "--", "/*", "*/"}
	for _, p := range injectionPatterns {
		if strings.Contains(status, p) || strings.Contains(dbPath, p) {
			return 0, fmt.Errorf("sql injection blocked")
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return 0, fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	var count float64
	row := db.QueryRow("SELECT COUNT(*) FROM transaction_log WHERE status = ?", status)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("query transaction count: %w", err)
	}
	return count, nil
}
