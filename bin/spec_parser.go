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

	// Build prompt from SQLite ledger (pure data-driven, no hardcoded constants)
	prompt, err := BuildPromptFromLedger(dbPath, spec.Title, spec.Invariants, spec.Classes, report)
	if err != nil {
		fmt.Printf("Error building prompt from ledger: %v\n", err)
		os.Exit(1)
	}

	// Save to synthesis prompt file
	promptFile := filepath.Join(projectDir, ".context", "synthesis_prompt.txt")
	err = os.WriteFile(promptFile, []byte(prompt), 0644)
	if err != nil {
		fmt.Printf("Error writing synthesis prompt file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[INFO] Synthesis prompt written to: %s\n", promptFile)
}

// BuildPromptFromLedger constructs the synthesis prompt from SQLite ledger data.
// Layer 1: grammar exemplar from language_substrates (static, KV cache optimized).
// Layer 2: report from symbol ledger (semi-static schema data).
// Layer 3: dynamic specification requirements (volatile suffix).
func BuildPromptFromLedger(dbPath string, title string, invariants []string, classes []SpecClass, report string) (string, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open context database: %w", err)
	}
	defer db.Close()

	var prompt strings.Builder

	// Layer 1: Static grammar exemplar from language_substrates table
	var grammarPrefix, structuralExemplar string
	err = db.QueryRow("SELECT grammar_prefix, structural_exemplar FROM language_substrates WHERE language_id = 'netrexx'").Scan(&grammarPrefix, &structuralExemplar)
	if err != nil {
		// No grammar found — fall back to minimal instructions
		prompt.WriteString("### LAYER 1: TARGET LANGUAGE\n")
		prompt.WriteString("Generate NetRexx source code. Use nominal imports (no wildcards). Modifiers follow class/method names.\n\n")
	} else {
		prompt.WriteString("### LAYER 1: FIXED TARGET GRAMMAR EXEMPLAR\n")
		prompt.WriteString(grammarPrefix)
		prompt.WriteString("\n\n")
		prompt.WriteString(structuralExemplar)
		prompt.WriteString("\n\n")
	}

	// Layer 2: Relational schema data (low volatility)
	prompt.WriteString("### LAYER 2: ACTIVE SYMBOL LEDGER AND SCHEMA TUPLES\n")
	prompt.WriteString(report)
	prompt.WriteString("\n\n")

	// Layer 3: Dynamic specification + instructions
	prompt.WriteString("### LAYER 3: TARGET REQUIREMENTS\n")
	prompt.WriteString(fmt.Sprintf("Target: %s\n\n", title))
	for _, inv := range invariants {
		prompt.WriteString(fmt.Sprintf("- %s\n", inv))
	}
	prompt.WriteString("\n")
	for _, class := range classes {
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
	prompt.WriteString("IMPORTANT: You MUST guard all database connection logic against null or empty/placeholder dbPath values. Every method taking a 'dbPath = String' parameter MUST check `if dbPath \\= null & dbPath \\= \"null\" then do` before connecting via JDBC, otherwise SQLite JDBC will physically create a database file named 'null' in the current working directory.\n")
	prompt.WriteString("Output ONLY the complete NetRexx source code matching the layer 1 structural pattern and layer 3 requirements. No explanations, no markdown block wrapping.\n")

	return prompt.String(), nil
}
