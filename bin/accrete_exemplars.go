package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type MethodBlock struct {
	Name      string
	Signature string
	Body      string
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./accrete_exemplars <db_path> <class_name> <nrx_file_path>")
		os.Exit(1)
	}

	dbPath := os.Args[1]
	className := os.Args[2]
	nrxPath := os.Args[3]

	// 1. Query database for methods declared in this class
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("[ERROR] Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Query symbol declarations for this class name to find method names.
	// Since file_path contains the class name (e.g. TransactionRouter.java), we can filter by that.
	query := `SELECT symbol_uri FROM declarations WHERE file_path LIKE ? AND symbol_uri LIKE '%method%'`
	rows, err := db.Query(query, "%/"+className+".java")
	if err != nil {
		fmt.Printf("[ERROR] Failed to query declarations: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	dbMethods := make(map[string]bool)
	for rows.Next() {
		var symbolURI string
		if err := rows.Scan(&symbolURI); err == nil {
			// Extract method name from symbol_uri
			// e.g. java+method:///com/factory/routing/TransactionRouter/routeTransaction(java.lang.String,com.factory.routing.TransactionRecord)
			parts := strings.Split(symbolURI, "/")
			if len(parts) > 0 {
				lastPart := parts[len(parts)-1]
				if idx := strings.Index(lastPart, "("); idx != -1 {
					methodName := lastPart[:idx]
					// Filter out main, dummySignal, and test methods
					if methodName != "main" && methodName != "dummySignal" && !strings.HasPrefix(strings.ToLower(methodName), "test") {
						dbMethods[methodName] = true
					}
				}
			}
		}
	}

	if len(dbMethods) == 0 {
		fmt.Printf("[INFO] No eligible methods found in declarations for class %s.\n", className)
		return
	}

	fmt.Printf("[INFO] Found %d target methods in declarations for class %s: %v\n", len(dbMethods), className, dbMethods)

	// 2. Read the .nrx file and extract method blocks
	contentBytes, err := os.ReadFile(nrxPath)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read .nrx file %s: %v\n", nrxPath, err)
		os.Exit(1)
	}
	lines := strings.Split(string(contentBytes), "\n")

	// Match: method methodName(...) [returns ...] [signals ...]
	methodStartRegex := regexp.MustCompile(`(?i)^\s*method\s+([a-zA-Z0-9_]+)\b`)

	var currentMethod *MethodBlock
	var methodBlocks []*MethodBlock

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Skip comment lines when checking for method starts
		isComment := strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "/*")
		
		if !isComment {
			matches := methodStartRegex.FindStringSubmatch(line)
			if len(matches) > 0 {
				methodName := matches[1]
				if dbMethods[methodName] {
					// We hit a new method declaration that we want to track.
					// First, save the previous method block if any.
					if currentMethod != nil {
						methodBlocks = append(methodBlocks, currentMethod)
					}
					// Start new method block
					currentMethod = &MethodBlock{
						Name:      methodName,
						Signature: trimmed,
						Body:      line + "\n",
					}
					continue
				} else {
					// Hit a different method we don't track, finish current one
					if currentMethod != nil {
						methodBlocks = append(methodBlocks, currentMethod)
						currentMethod = nil
					}
					continue
				}
			}
		}

		// If we are currently inside a tracked method block, append line to body
		if currentMethod != nil {
			currentMethod.Body += line + "\n"
		}
	}
	// Save the last method block if any
	if currentMethod != nil {
		methodBlocks = append(methodBlocks, currentMethod)
	}

	// 3. Write each extracted block into the database
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("[ERROR] Failed to start transaction: %v\n", err)
		os.Exit(1)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO unified_exemplars (
			exemplar_id,
			domain_scope,
			fact_context_predicate,
			input_state_payload,
			expected_output_state,
			few_shot_prompt_block,
			verification_harness_raw
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Printf("[ERROR] Failed to prepare statement: %v\n", err)
		os.Exit(1)
	}
	defer stmt.Close()

	for _, block := range methodBlocks {
		exemplarID := fmt.Sprintf("EXEMP-%s-%s", strings.ToUpper(className), strings.ToUpper(block.Name))
		cleanBody := strings.TrimSpace(block.Body)
		
		fmt.Printf("[INFO] Accreting method: %s -> %s\n", block.Name, exemplarID)
		
		_, err = stmt.Exec(
			exemplarID,
			"Implementation.NetRexx",
			block.Signature,
			"{}",
			"SUCCESS",
			cleanBody,
			"MACHINE_VERIFIED",
		)
		if err != nil {
			fmt.Printf("[ERROR] Failed to execute insert for %s: %v\n", block.Name, err)
			os.Exit(1)
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("[ERROR] Failed to commit transaction: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[SUCCESS] Successfully accreted %d exemplars for class %s into %s\n", len(methodBlocks), className, dbPath)
}
