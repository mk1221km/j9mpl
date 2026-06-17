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
    // Input validation
    dbPath = strings.TrimSpace(dbPath)
    if len(dbPath) == 0 {
        return fmt.Errorf("invalid input: dbPath is empty")
    }
    if strings.Contains(dbPath, "..") || strings.Contains(dbPath, "/etc/") {
        return fmt.Errorf("path traversal blocked")
    }
    sqlKeywords := []string{"SELECT", "UPDATE", "DELETE", "INSERT", "DROP", "ALTER", "CREATE", "--", ";", "OR", "AND"}
    for _, keyword := range sqlKeywords {
        if strings.Contains(strings.ToUpper(dbPath), keyword) {
            return fmt.Errorf("sql injection blocked")
        }
    }

    // Open database
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    defer db.Close()

    // Create routing_rules table
    createRulesTable := `
        CREATE TABLE IF NOT EXISTS routing_rules (
            min_amount REAL NOT NULL,
            priority TEXT NOT NULL,
            channel TEXT NOT NULL,
            PRIMARY KEY (min_amount, priority)
        );`
    if _, err := db.Exec(createRulesTable); err != nil {
        return fmt.Errorf("failed to create routing_rules table: %w", err)
    }

    // Create transaction_log table
    createLogTable := `
        CREATE TABLE IF NOT EXISTS transaction_log (
            tx_id TEXT PRIMARY KEY,
            sender TEXT NOT NULL,
            receiver TEXT NOT NULL,
            amount REAL NOT NULL,
            channel TEXT NOT NULL,
            status TEXT NOT NULL
        );`
    if _, err := db.Exec(createLogTable); err != nil {
        return fmt.Errorf("failed to create transaction_log table: %w", err)
    }

    // Default routing rules
    type defaultRule struct {
        minAmount float64
        priority  string
        channel   string
    }
    defaultRules := []defaultRule{
        {minAmount: 0, priority: "normal", channel: "default_channel"},
        {minAmount: 1000, priority: "high", channel: "premium_channel"},
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
	// Validate all input arguments at the beginning
	if strings.TrimSpace(dbPath) == "" {
		return "", fmt.Errorf("invalid input: dbPath is empty")
	}
	if strings.Contains(dbPath, "..") || strings.Contains(dbPath, "/etc/") {
		return "", fmt.Errorf("path traversal blocked")
	}
	if strings.Contains(dbPath, "'") || strings.Contains(dbPath, ";") || strings.Contains(dbPath, "--") {
		return "", fmt.Errorf("sql injection blocked")
	}
	if strings.TrimSpace(record.TxId) == "" {
		return "", fmt.Errorf("invalid input: TxId is empty")
	}
	if strings.TrimSpace(record.Sender) == "" {
		return "", fmt.Errorf("invalid input: Sender is empty")
	}
	if strings.TrimSpace(record.Receiver) == "" {
		return "", fmt.Errorf("invalid input: Receiver is empty")
	}
	if strings.TrimSpace(record.Priority) == "" {
		return "", fmt.Errorf("invalid input: Priority is empty")
	}
	// Check SQL injection on string fields of record
	for _, s := range []string{record.TxId, record.Sender, record.Receiver, record.Priority} {
		if strings.Contains(s, "'") || strings.Contains(s, ";") || strings.Contains(s, "--") {
			return "", fmt.Errorf("sql injection blocked")
		}
	}

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	// Query routing rule based on amount and priority
	var channel string
	row := db.QueryRow("SELECT channel FROM routing_rules WHERE min_amount = ? AND priority = ?", record.Amount, record.Priority)
	if err := row.Scan(&channel); err != nil {
		return "", fmt.Errorf("query routing rule: %w", err)
	}

	// Insert transaction log
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
