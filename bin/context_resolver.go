package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	// standard CGO SQLite driver. 
	// For CGO-free compilation, use: _ "modernc.org/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

// SymbolDetails holds parsed logical details from an M3 URI
type SymbolDetails struct {
	Type string // "class", "method", "constructor", "field", "unknown"
	Name string // The logical symbol name
}

// ParseSymbolURI extracts type and name from logical M3 URI.
// Example: |java+method:///com/factory/telemetry/TelemetryEngine/processData(...)| -> "method", "processData"
func ParseSymbolURI(uri string) SymbolDetails {
	cleaned := strings.Trim(uri, "|")
	
	symbolType := "unknown"
	if strings.Contains(cleaned, "java+method") {
		symbolType = "method"
	} else if strings.Contains(cleaned, "java+constructor") {
		symbolType = "constructor"
	} else if strings.Contains(cleaned, "java+class") {
		symbolType = "class"
	} else if strings.Contains(cleaned, "java+field") {
		symbolType = "field"
	}

	// Extract everything after the last slash of the path portion
	parts := strings.Split(cleaned, "///")
	if len(parts) < 2 {
		return SymbolDetails{Type: symbolType, Name: ""}
	}
	
	pathPart := parts[len(parts)-1]
	subParts := strings.Split(pathPart, "/")
	namePart := subParts[len(subParts)-1]
	
	// Strip parameter list if it's a method/constructor
	if idx := strings.Index(namePart, "("); idx != -1 {
		namePart = namePart[:idx]
	}

	return SymbolDetails{Type: symbolType, Name: namePart}
}

// ResolveContext maps a symbol URI to its original NetRexx source block
func ResolveContext(dbPath, symbolURI string) (string, error) {
	// 1. Query the SQLite context ledger for the translated Java source file path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open context database: %w", err)
	}
	defer db.Close()

	var javaPath string
	query := "SELECT file_path FROM declarations WHERE symbol_uri = ?"
	err = db.QueryRow(query, symbolURI).Scan(&javaPath)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("symbol URI not found in declarations ledger: %s", symbolURI)
	} else if err != nil {
		return "", fmt.Errorf("failed to query declarations table: %w", err)
	}

	// 2. Map the intermediate Java path to the corresponding NetRexx source path
	nrxPath := strings.Replace(javaPath, ".java", ".nrx", 1)

	// 3. Parse logical symbol information
	details := ParseSymbolURI(symbolURI)
	if details.Type == "unknown" || details.Name == "" {
		// If symbol structure is unresolvable, fallback to returning the full file content
		content, err := os.ReadFile(nrxPath)
		if err != nil {
			return "", fmt.Errorf("failed to read full source fallback file: %w", err)
		}
		return string(content), nil
	}

	// 4. Open NetRexx source file and perform straight-line regex scan
	file, err := os.Open(nrxPath)
	if err != nil {
		return "", fmt.Errorf("failed to open NetRexx source file: %w", err)
	}
	defer file.Close()

	// Compile class and method matching regexes
	classRegex := regexp.MustCompile(`(?i)^\s*class\s+(\w+)`)
	methodRegex := regexp.MustCompile(`(?i)^\s*method\s+(\w+)`)

	var extractedLines []string
	inBlock := false
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		
		if !inBlock {
			// Search for block entry boundary
			if details.Type == "class" {
				if matches := classRegex.FindStringSubmatch(line); len(matches) > 1 {
					if strings.EqualFold(matches[1], details.Name) {
						inBlock = true
						extractedLines = append(extractedLines, line)
					}
				}
			} else if details.Type == "method" || details.Type == "constructor" {
				if matches := methodRegex.FindStringSubmatch(line); len(matches) > 1 {
					if strings.EqualFold(matches[1], details.Name) {
						inBlock = true
						extractedLines = append(extractedLines, line)
					}
				}
			}
		} else {
			// Check block exit boundary (another class or method begins)
			isClass := classRegex.MatchString(line)
			isMethod := methodRegex.MatchString(line)

			if isClass || isMethod {
				// Stop capturing if we are in a method/constructor, or if we hit a new class definition
				if details.Type == "method" || details.Type == "constructor" || isClass {
					break
				}
			}
			extractedLines = append(extractedLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning NetRexx source file: %w", err)
	}

	if len(extractedLines) == 0 {
		return "", fmt.Errorf("symbol '%s' of type '%s' not found in file %s", details.Name, details.Type, nrxPath)
	}

	return strings.Join(extractedLines, "\n"), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./context_resolver <symbol_uri>")
		os.Exit(1)
	}

	symbolURI := os.Args[1]

	// Determine DB path relative to the binary location
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error resolving executable path: %v\n", err)
		os.Exit(1)
	}
	
	projectDir := filepath.Dir(filepath.Dir(execPath))
	dbPath := filepath.Join(projectDir, ".context", "project_context.db")

	resolvedCode, err := ResolveContext(dbPath, symbolURI)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(resolvedCode)
}
