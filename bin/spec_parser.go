package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type SpecMethod struct {
	Name         string
	Args         string
	Returns      string
	Requirements string
}

type SpecClass struct {
	Name    string
	Methods []SpecMethod
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

	classRegex := regexp.MustCompile(`(?i)(?:class|struct)\s+(\w+)`)
	methodRegex := regexp.MustCompile(`^\s*(?:\d+\.|\*|-)\s+(\w+)\((.*?)\)(?:\s*(returns\s+\w+))?\s*:\s*(.*)`)
	fieldRegex := regexp.MustCompile(`^\s*(?:\*|-)\s+(\w+)\s*(?:\((.*?)\))?`)
	titleRegex := regexp.MustCompile(`^#\s+(.*)`)

	var currentClass *SpecClass
	inDTOs := false
	inInterfaces := false
	inInvariants := false

	for scanner.Scan() {
		line := scanner.Text()
		// Strip backticks before processing to handle different signature wrapping styles
		line = strings.ReplaceAll(line, "`", "")
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
				// Check for methods (only in interfaces/methods sections)
				if inInterfaces && methodRegex.MatchString(trimmed) {
					if matches := methodRegex.FindStringSubmatch(trimmed); len(matches) > 0 {
						name := matches[1]
						args := matches[2]
						returns := ""
						reqs := ""
						if len(matches) > 3 {
							returns = strings.TrimSpace(matches[3])
						}
						if len(matches) > 4 {
							reqs = strings.TrimSpace(matches[4])
						}
						currentClass.Methods = append(currentClass.Methods, SpecMethod{
							Name:         name,
							Args:         args,
							Returns:      returns,
							Requirements: reqs,
						})
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
				report.WriteString(fmt.Sprintf("  - [MISSING] Method %s.%s\n", class.Name, method.Name))
			}
		} else if err != nil {
			return "", err
		} else {
			report.WriteString(fmt.Sprintf("[FOUND] Class %s exists in codebase (File: %s, URI: %s)\n", class.Name, filePath, classURI))
			
			for _, method := range class.Methods {
				var methodURI string
				methodQuery := "SELECT symbol_uri FROM declarations WHERE symbol_uri LIKE ? LIMIT 1"
				mErr := db.QueryRow(methodQuery, "%method%/"+class.Name+"/"+method.Name+"%").Scan(&methodURI)

				if mErr == sql.ErrNoRows {
					report.WriteString(fmt.Sprintf("  - [MISSING] Method %s.%s\n", class.Name, method.Name))
				} else if mErr != nil {
					return "", mErr
				} else {
					report.WriteString(fmt.Sprintf("  - [FOUND] Method %s.%s (URI: %s)\n", class.Name, method.Name, methodURI))
				}
			}
		}
		report.WriteString("\n")
	}

	return report.String(), nil
}

func convertArgsToNetRexx(args string) string {
	if strings.TrimSpace(args) == "" {
		return ""
	}
	parts := strings.Split(args, ",")
	var nrxParts []string
	for _, p := range parts {
		subParts := strings.Split(p, ":")
		if len(subParts) == 2 {
			name := strings.TrimSpace(subParts[0])
			typ := strings.TrimSpace(subParts[1])
			nrxParts = append(nrxParts, fmt.Sprintf("%s = %s", name, typ))
		} else {
			nrxParts = append(nrxParts, strings.TrimSpace(p))
		}
	}
	return strings.Join(nrxParts, ", ")
}

func GenerateGoSkeleton(dbPath string, spec ParsedSpec, mainClassName string) (string, error) {
	var sb strings.Builder
	
	// Package name derived from DB or default
	packageName := mainClassName
	db, err := sql.Open("sqlite3", dbPath)
	if err == nil {
		defer db.Close()
		var symbolURI string
		classQuery := "SELECT symbol_uri FROM declarations WHERE symbol_uri LIKE ? LIMIT 1"
		err = db.QueryRow(classQuery, "%class%/"+mainClassName+"%").Scan(&symbolURI)
		if err == nil {
			cleaned := strings.Trim(symbolURI, "|")
			parts := strings.Split(cleaned, "///")
			if len(parts) >= 2 {
				pathPart := parts[len(parts)-1]
				subParts := strings.Split(pathPart, "/")
				if len(subParts) > 1 {
					pkg := strings.ToLower(strings.Join(subParts[:len(subParts)-1], ""))
					if pkg != "" {
						packageName = pkg
					}
				}
			}
		}
	}

	sb.WriteString(fmt.Sprintf("package %s\n\n", strings.ToLower(strings.ReplaceAll(packageName, ".", ""))))

	// Collect all field types to resolve dynamic imports
	fieldTypes := make(map[string]bool)
	for _, class := range spec.Classes {
		if len(class.Methods) == 0 {
			for _, field := range class.Fields {
				parts := strings.Split(field, "(")
				if len(parts) > 1 {
					t := goType(strings.TrimSpace(strings.Trim(parts[1], "()")))
					fieldTypes[t] = true
				}
			}
		}
	}

	// Read base imports from unified_exemplars ledger
	baseImports := ""
	if db != nil {
		db.QueryRow("SELECT default_imports FROM unified_exemplars WHERE exemplar_id='NETREXX_GRAMMAR_BASICS'").Scan(&baseImports)
	}
	if baseImports == "" {
		baseImports = "import (\"database/sql\")"
	}

	// Query structural_import_mappings for dynamic package requirements
	var dynamicPkgs []string
	if db != nil {
		for t := range fieldTypes {
			var pkg string
			err := db.QueryRow("SELECT required_package FROM structural_import_mappings WHERE primitive_type=? AND target_language='go'", t).Scan(&pkg)
			if err == nil && !strings.Contains(baseImports, pkg) {
				dynamicPkgs = append(dynamicPkgs, pkg)
			}
		}
	}

	// Build final import block
	if len(dynamicPkgs) > 0 {
		// Extract the import block prefix (everything before the first import entry)
		if strings.HasPrefix(baseImports, "import (") {
			// Insert dynamic packages before the closing paren
			insert := "\n\t" + strings.Join(dynamicPkgs, "\n\t") + "\n"
			baseImports = baseImports[:len(baseImports)-1] + insert + ")"
		} else if strings.HasPrefix(baseImports, "import \"") {
			// Single-line import, expand to multi-line
			pkg := strings.TrimPrefix(baseImports, "import ")
			baseImports = "import (\n\t" + pkg + "\n"
			for _, dp := range dynamicPkgs {
				baseImports += "\t" + dp + "\n"
			}
			baseImports += ")"
		}
	}
	sb.WriteString(baseImports)
	sb.WriteString("\n")

	// DTO Structs (data-only classes)
	for _, class := range spec.Classes {
		if len(class.Methods) == 0 {
			sb.WriteString(fmt.Sprintf("\ntype %s struct {\n", class.Name))
			for _, field := range class.Fields {
				parts := strings.Split(field, "(")
				fName := strings.TrimSpace(parts[0])
				if fName == "" {
					continue
				}
				fType := "string"
				if len(parts) > 1 {
					fType = goType(strings.TrimSpace(strings.Trim(parts[1], "()")))
				}
				// Uppercase first letter for Go convention (exported fields)
				fNameUpper := strings.ToUpper(fName[:1]) + fName[1:]
				sb.WriteString(fmt.Sprintf("\t%s %s\n", fNameUpper, fType))
			}
			sb.WriteString("}\n")
		}
	}

	// Main struct with method stubs
	sb.WriteString(fmt.Sprintf("\ntype %s struct{}\n", mainClassName))
	for _, class := range spec.Classes {
		if len(class.Methods) > 0 && class.Name == mainClassName {
			for _, m := range class.Methods {
				sb.WriteString(fmt.Sprintf("\n// SKELETON_%s\n", m.Name))
			}
		}
	}

	return sb.String(), nil
}

func goType(netrexxType string) string {
	switch strings.ToLower(strings.TrimSpace(netrexxType)) {
	case "string":
		return "string"
	case "int", "integer":
		return "int"
	case "rexx":
		return "float64"
	case "double", "real", "float":
		return "float64"
	case "boolean", "bool":
		return "bool"
	case "connection":
		return "*sql.DB"
	case "preparedstatement":
		return "*sql.Stmt"
	case "resultset":
		return "*sql.Rows"
	default:
		return netrexxType
	}
}

func convertArgsToGo(args string) string {
	if args == "" {
		return ""
	}
	// Handle spec format like "dbPath: String" or "dbPath = String"
	parts := strings.Split(args, ",")
	var goParts []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		name, typ := "", ""
		if strings.Contains(p, "=") {
			sub := strings.SplitN(p, "=", 2)
			name = strings.TrimSpace(sub[0])
			typ = strings.TrimSpace(sub[1])
		} else if strings.Contains(p, ":") {
			sub := strings.SplitN(p, ":", 2)
			name = strings.TrimSpace(sub[0])
			typ = strings.TrimSpace(sub[1])
		} else {
			// No type: assume string
			goParts = append(goParts, fmt.Sprintf("%s string", strings.TrimSpace(p)))
			continue
		}
		goParts = append(goParts, fmt.Sprintf("%s %s", name, goType(typ)))
	}
	return strings.Join(goParts, ", ")
}

func BuildMethodPrompt(dbPath string, mainClassName string, method SpecMethod, invariants []string, report string, classes []SpecClass) (string, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open context database: %w", err)
	}
	defer db.Close()

	var prompt strings.Builder

	// Layer 1: Static grammar exemplar from unified_exemplars table
	var grammarPrefix, structuralExemplar string
	err = db.QueryRow("SELECT fact_context_predicate, few_shot_prompt_block FROM unified_exemplars WHERE exemplar_id = 'NETREXX_GRAMMAR_BASICS'").Scan(&grammarPrefix, &structuralExemplar)
	if err == nil {
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

	// Layer 2.5: Helper Classes and Data Transfer Objects (DTOs)
	prompt.WriteString("### LAYER 2.5: HELPER CLASSES AND DATA TRANSFER OBJECTS (DTOs)\n")
	for _, class := range classes {
		if len(class.Methods) == 0 {
			prompt.WriteString(fmt.Sprintf("Class %s represents data transfer structures and exposes public properties:\n", class.Name))
			for _, field := range class.Fields {
				prompt.WriteString(fmt.Sprintf("- %s\n", field))
			}
			prompt.WriteString("\n")
		}
	}
	prompt.WriteString("\n")

	// Layer 3: Exemplar Blocks (from unified_exemplars table for SQLite/JDBC references)
	prompt.WriteString("### LAYER 3: JDBC/SQLITE REFERENCE EXEMPLARS\n")
	rows, err := db.Query("SELECT few_shot_prompt_block FROM unified_exemplars WHERE domain_scope LIKE '%SQLite%'")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var snippet string
			if err := rows.Scan(&snippet); err == nil {
				prompt.WriteString(snippet)
				prompt.WriteString("\n\n")
			}
		}
	}

	// Layer 4: Specific Method Target Requirements
	prompt.WriteString("### LAYER 4: TARGET METHOD REQUIREMENTS\n")
	prompt.WriteString(fmt.Sprintf("Package: %s\n", mainClassName))
	
	goArgs := convertArgsToGo(method.Args)
	goRet := method.Returns
	if goRet != "" {
		goRet = " " + goRet
	}

	// In Go, main() is a standalone package-level function, not a method
	var sig string
	if method.Name == "main" {
		sig = "func main() error"
		prompt.WriteString(fmt.Sprintf("Target Signature: %s\n", sig))
		prompt.WriteString("NOTE: main() is a standalone function, not a method on TransactionRouter.\n")
		prompt.WriteString("Create the database connection inside main() using sql.Open().\n")
		prompt.WriteString("Call helper functions via a local TransactionRouter instance.\n")
		prompt.WriteString("IMPORTANT: Use PascalCase field names when accessing struct fields: record.TxId (NOT record.txId).\n")
	} else {
		sig = fmt.Sprintf("func (s *%s) %s(%s)%s", mainClassName, method.Name, goArgs, goRet)
		if goRet == "" {
			sig = fmt.Sprintf("func (s *%s) %s(%s)", mainClassName, method.Name, goArgs)
		}
		prompt.WriteString(fmt.Sprintf("Target Signature: %s\n", sig))
	}
	prompt.WriteString(fmt.Sprintf("Requirements: %s\n\n", method.Requirements))

	prompt.WriteString("Architectural Invariants & Database Schema:\n")
	for _, inv := range invariants {
		prompt.WriteString(fmt.Sprintf("- %s\n", inv))
	}
	prompt.WriteString("\n")

	prompt.WriteString("VALIDATION & ERROR HANDLING:\n")
	prompt.WriteString("You MUST validate all input arguments at the beginning of the function (before any database operation):\n")
	prompt.WriteString("1. If a string parameter is empty or whitespace-only, return an error: `return 0, fmt.Errorf(\"invalid input: %%s is empty\", name)`\n")
	prompt.WriteString("2. If a dbPath parameter contains path traversal (\"..\", \"/etc/\"), return `fmt.Errorf(\"path traversal blocked\")`\n")
	prompt.WriteString("3. If a string parameter contains SQL injection patterns, return `fmt.Errorf(\"sql injection blocked\")`\n")
	prompt.WriteString("4. Use Go-style error returns, not Java exceptions. Return `(result, error)` where error is nil on success.\n")
	prompt.WriteString("Ensure these checks are at the very beginning of each method body and are NOT caught by any internal try-catch blocks that handle database operations.\n\n")

	prompt.WriteString("IMPORTANT: You are generating Go code. Output ONLY the complete Go function body for '" + sig + "'. Do not include the enclosing struct definition, package, or imports. Do not wrap in markdown code blocks. Use standard Go: `func (s *Type) Name(args) (returnType, error) { ... }` with proper error handling. Use `database/sql` for queries. Always return an error as the last return value.\n")
	prompt.WriteString("NOTE: Struct fields use PascalCase — access them as record.TxId (NOT record.txId). Use Sender, Receiver, Priority, Amount with capital letters.\n")

	return prompt.String(), nil
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

	mainClassName := strings.TrimSuffix(filepath.Base(specPath), "Spec.md")

	// Emit class skeleton immediately to disk
	skeleton, err := GenerateGoSkeleton(dbPath, spec, mainClassName)
	if err != nil {
		fmt.Printf("Error generating skeleton: %v\n", err)
		os.Exit(1)
	}

	skeletonFile := filepath.Join(projectDir, "generated", mainClassName+".go")
	skipWrite := false
	if existingContent, readErr := os.ReadFile(skeletonFile); readErr == nil {
		if !strings.Contains(string(existingContent), "SKELETON_") {
			skipWrite = true
		}
	}
	if skipWrite {
		fmt.Printf("[INFO] %s already synthesized. Skipping skeleton overwrite.\n", skeletonFile)
	} else {
		err = os.WriteFile(skeletonFile, []byte(skeleton), 0644)
		if err != nil {
			fmt.Printf("Error writing skeleton file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("[INFO] Structural skeleton written to: %s\n", skeletonFile)
	}

	// Setup context directory
	contextDir := filepath.Join(projectDir, ".context")
	os.MkdirAll(contextDir, 0755)

	// Extract fuzzer boundaries from SQLite to .context/fuzzer_boundaries.json
	err = ExtractFuzzerBoundaries(dbPath, filepath.Join(contextDir, "fuzzer_boundaries.json"))
	if err != nil {
		fmt.Printf("Error extracting fuzzer boundaries: %v\n", err)
		os.Exit(1)
	}

	// Write methods list
	var methodNames []string
	var targetClass *SpecClass
	for _, class := range spec.Classes {
		if class.Name == mainClassName {
			targetClass = &class
			break
		}
	}

	if targetClass == nil {
		fmt.Printf("Error: Main class %s not found in parsed specification classes.\n", mainClassName)
		os.Exit(1)
	}

	for _, m := range targetClass.Methods {
		methodNames = append(methodNames, m.Name)

		// Generate method skeleton block
		nrxArgs := convertArgsToNetRexx(m.Args)
		retClause := ""
		retType := "void"
		if m.Returns != "" {
			retClause = " " + m.Returns
			retType = strings.TrimSpace(strings.TrimPrefix(m.Returns, "returns"))
		}

		var mSkeleton strings.Builder
		mSkeleton.WriteString(fmt.Sprintf("  method %s(%s) public static%s\n", m.Name, nrxArgs, retClause))
		mSkeleton.WriteString(fmt.Sprintf("    -- SKELETON_%s\n", m.Name))
		if retType == "void" || m.Name == "main" {
			mSkeleton.WriteString("    nop\n")
		} else {
			mSkeleton.WriteString(fmt.Sprintf("    return %s null\n", retType))
		}

		mSkeletonFile := filepath.Join(contextDir, fmt.Sprintf("skeleton_%s.txt", m.Name))
		err = os.WriteFile(mSkeletonFile, []byte(mSkeleton.String()), 0644)
		if err != nil {
			fmt.Printf("Error writing method skeleton file for %s: %v\n", m.Name, err)
			os.Exit(1)
		}

		// Generate method prompt
		mPrompt, err := BuildMethodPrompt(dbPath, mainClassName, m, spec.Invariants, report, spec.Classes)
		if err != nil {
			fmt.Printf("Error building method prompt for %s: %v\n", m.Name, err)
			os.Exit(1)
		}

		mPromptFile := filepath.Join(contextDir, fmt.Sprintf("prompt_%s.txt", m.Name))
		err = os.WriteFile(mPromptFile, []byte(mPrompt), 0644)
		if err != nil {
			fmt.Printf("Error writing method prompt file for %s: %v\n", m.Name, err)
			os.Exit(1)
		}
	}

	methodsFile := filepath.Join(contextDir, "methods.txt")
	err = os.WriteFile(methodsFile, []byte(strings.Join(methodNames, "\n")), 0644)
	if err != nil {
		fmt.Printf("Error writing methods list file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[INFO] Methods list written to: %s\n", methodsFile)

	// Build fallback synthesis prompt just in case
	fallbackPrompt, err := BuildMethodPrompt(dbPath, mainClassName, SpecMethod{
		Name:         "all",
		Args:         "",
		Returns:      "",
		Requirements: "Generate all methods.",
	}, spec.Invariants, report, spec.Classes)
	if err == nil {
		os.WriteFile(filepath.Join(contextDir, "synthesis_prompt.txt"), []byte(fallbackPrompt), 0644)
	}
}

type CounterItem struct {
	Value    string `json:"value"`
	Expected string `json:"expected"`
}

type BoundaryPayload struct {
	Valid   []string      `json:"valid"`
	Counter []CounterItem `json:"counter"`
}

type BoundaryItem struct {
	Domain  string          `json:"domain"`
	Payload BoundaryPayload `json:"payload"`
}

func ExtractFuzzerBoundaries(dbPath string, outputPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT exemplar_id, fact_context_predicate, input_state_payload FROM unified_exemplars WHERE domain_scope = 'Fuzzer.Boundary'")
	if err != nil {
		return err
	}
	defer rows.Close()

	boundaries := make(map[string][]BoundaryItem)

	for rows.Next() {
		var exemplarId, targetPrimitive, inputStatePayload string
		if err := rows.Scan(&exemplarId, &targetPrimitive, &inputStatePayload); err != nil {
			return err
		}

		var payload BoundaryPayload
		if err := json.Unmarshal([]byte(inputStatePayload), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal fuzzer payload for %s: %w", exemplarId, err)
		}

		item := BoundaryItem{
			Domain:  exemplarId,
			Payload: payload,
		}

		boundaries[targetPrimitive] = append(boundaries[targetPrimitive], item)
	}

	jsonBytes, err := json.MarshalIndent(boundaries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, jsonBytes, 0644)
}
