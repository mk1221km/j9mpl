package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type SpecClass struct {
	Name    string
	Methods []string
	Fields  []string
}

type ParsedSpec struct {
	Title      string
	Classes    []SpecClass
	Invariants []string
}

// ParseMarkdownSpec parses the structural spec markdown file
func ParseMarkdownSpec(filePath string) (ParsedSpec, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return ParsedSpec{}, err
	}
	defer file.Close()

	var spec ParsedSpec
	scanner := bufio.NewScanner(file)

	classRegex := regexp.MustCompile(`(?i)(?:class|struct)\s+\x60?(\w+)\x60?`)
	methodRegex := regexp.MustCompile(`^\s*(?:\d+\.|\*|-)\s+\x60?(\w+)\((.*?)\)`)
	fieldRegex := regexp.MustCompile(`^\s*(?:\*|-)\s+\x60?(\w+)\x60?\s*(?:\((.*?)\))?`)
	titleRegex := regexp.MustCompile(`^#\s+(.*)`)

	var currentClass *SpecClass
	inDTOs := false
	inInterfaces := false
	inInvariants := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Title extraction
		if spec.Title == "" {
			if matches := titleRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
				spec.Title = matches[1]
				continue
			}
		}

		// Section tracking
		if strings.HasPrefix(trimmed, "##") {
			lower := strings.ToLower(trimmed)
			inDTOs = strings.Contains(lower, "dto") || strings.Contains(lower, "data transfer")
			inInterfaces = strings.Contains(lower, "interface") || strings.Contains(lower, "method") || strings.Contains(lower, "api")
			inInvariants = strings.Contains(lower, "invariant") || strings.Contains(lower, "layout") || strings.Contains(lower, "rule")
			currentClass = nil
			continue
		}

		if inInvariants {
			if strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "- ") {
				inv := strings.TrimPrefix(strings.TrimPrefix(trimmed, "* "), "- ")
				spec.Invariants = append(spec.Invariants, inv)
			}
		}

		if inDTOs || inInterfaces {
			// Check for new class definition
			if matches := classRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
				className := matches[1]
				spec.Classes = append(spec.Classes, SpecClass{Name: className})
				currentClass = &spec.Classes[len(spec.Classes)-1]
				continue
			}

			if currentClass != nil {
				// Check for methods (only in interfaces/methods sections or if contains parentheses)
				if inInterfaces && methodRegex.MatchString(trimmed) {
					if matches := methodRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
						methodSig := fmt.Sprintf("%s(%s)", matches[1], matches[2])
						currentClass.Methods = append(currentClass.Methods, methodSig)
					}
				} else if inDTOs && fieldRegex.MatchString(trimmed) {
					// Check for fields (only in DTOs section)
					if matches := fieldRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
						fieldName := matches[1]
						fieldType := "Rexx"
						if len(matches) > 2 && matches[2] != "" {
							fieldType = matches[2]
						}
						currentClass.Fields = append(currentClass.Fields, fmt.Sprintf("%s (%s)", fieldName, fieldType))
					}
				}
			}
		}
	}

	return spec, scanner.Err()
}

// CheckCodeModel checks if classes and methods already exist in the SQLite context database
func CheckCodeModel(dbPath string, spec ParsedSpec) (string, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var report strings.Builder
	report.WriteString("=== Codebase Existence Ledger Analysis ===\n")

	for _, class := range spec.Classes {
		var classURI, filePath string
		classQuery := "SELECT symbol_uri, file_path FROM declarations WHERE symbol_uri LIKE ? LIMIT 1"
		err := db.QueryRow(classQuery, "%class%/"+class.Name+"%").Scan(&classURI, &filePath)

		if err == sql.ErrNoRows {
			report.WriteString(fmt.Sprintf("[MISSING] Class %s not found in context ledger.\n", class.Name))
			for _, method := range class.Methods {
				nameOnly := strings.Split(method, "(")[0]
				report.WriteString(fmt.Sprintf("  - [MISSING] Method %s.%s\n", class.Name, nameOnly))
			}
		} else if err != nil {
			return "", err
		} else {
			report.WriteString(fmt.Sprintf("[FOUND] Class %s exists in codebase (File: %s, URI: %s)\n", class.Name, filePath, classURI))
			
			for _, method := range class.Methods {
				nameOnly := strings.Split(method, "(")[0]
				var methodURI string
				methodQuery := "SELECT symbol_uri FROM declarations WHERE symbol_uri LIKE ? LIMIT 1"
				mErr := db.QueryRow(methodQuery, "%method%/"+class.Name+"/"+nameOnly+"%").Scan(&methodURI)

				if mErr == sql.ErrNoRows {
					report.WriteString(fmt.Sprintf("  - [MISSING] Method %s.%s\n", class.Name, nameOnly))
				} else if mErr != nil {
					return "", mErr
				} else {
					report.WriteString(fmt.Sprintf("  - [FOUND] Method %s.%s (URI: %s)\n", class.Name, nameOnly, methodURI))
				}
			}
		}
		report.WriteString("\n")
	}

	return report.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./spec_parser <spec_markdown_file>")
		os.Exit(1)
	}

	specPath := os.Args[1]

	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error resolving executable path: %v\n", err)
		os.Exit(1)
	}

	projectDir := filepath.Dir(filepath.Dir(execPath))
	dbPath := filepath.Join(projectDir, ".context", "project_context.db")

	spec, err := ParseMarkdownSpec(specPath)
	if err != nil {
		fmt.Printf("Error parsing specification: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[INFO] Successfully parsed specification: %s\n", spec.Title)

	report, err := CheckCodeModel(dbPath, spec)
	if err != nil {
		fmt.Printf("Error checking context ledger: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(report)

	// Format synthesis prompt payload
	var prompt strings.Builder
	prompt.WriteString("=== SPECIFICATION INGEST SYNTHESIS PROMPT ===\n")
	prompt.WriteString(fmt.Sprintf("Target Specification: %s\n\n", spec.Title))
	
	prompt.WriteString("Architectural Invariants & Constraints:\n")
	for _, inv := range spec.Invariants {
		prompt.WriteString(fmt.Sprintf("- %s\n", inv))
	}
	prompt.WriteString("\n")

	prompt.WriteString("Target Class Specifications:\n")
	for _, class := range spec.Classes {
		prompt.WriteString(fmt.Sprintf("Class: %s\n", class.Name))
		if len(class.Fields) > 0 {
			prompt.WriteString("  Fields:\n")
			for _, f := range class.Fields {
				prompt.WriteString(fmt.Sprintf("    - %s\n", f))
			}
		}
		if len(class.Methods) > 0 {
			prompt.WriteString("  Methods:\n")
			for _, m := range class.Methods {
				prompt.WriteString(fmt.Sprintf("    - %s\n", m))
			}
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("Ledger Coverage Analysis:\n")
	prompt.WriteString(report)

	prompt.WriteString("\nGenerate the complete NetRexx source file meeting this specification.\n")
	prompt.WriteString("IMPORTANT: You MUST guard all database connection logic against null or empty/placeholder dbPath values. Every method taking a 'dbPath = String' parameter MUST check `if dbPath \\= null & dbPath \\= \"null\" then do` before connecting via JDBC, otherwise SQLite JDBC will physically create a database file named 'null' in the current working directory. Ensure this check wraps all JDBC code in those methods.\n")
	prompt.WriteString("Reuse and align structural types where marked as [FOUND]. Synthesize all blocks marked as [MISSING].\n")
	prompt.WriteString("Output ONLY the complete revised NetRexx source code. No explanations, no markdown block wrapping.\n")

	// Save to synthesis prompt file
	promptFile := filepath.Join(projectDir, ".context", "synthesis_prompt.txt")
	err = os.WriteFile(promptFile, []byte(prompt.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing synthesis prompt file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[INFO] Synthesis prompt written to: %s\n", promptFile)
}
