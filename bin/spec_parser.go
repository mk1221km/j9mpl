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

	classRegex := regexp.MustCompile(`(?i)(?:class|struct)\s+` + "`?" + `(\w+)`)
	methodRegex := regexp.MustCompile("^\\s*(?:\\d+\\.|\\*|-)\\s+`?(\\w+)`?\\((.*?)\\)\\s*(.*)")
	fieldRegex := regexp.MustCompile(`^\s*(?:\*|-)\s+` + "`?" + `(\w+)` + "`?" + `\s*(?:\((.*?)\))?`)
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
			inDTOs = strings.Contains(lower, "dto") || strings.Contains(lower, "data transfer") || strings.Contains(lower, "data struct")
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
							fullRet := strings.TrimSpace(matches[3])
							returns = fullRet
							// Split on ':' to separate return type from description
							if idx := strings.Index(fullRet, ":"); idx >= 0 {
								returns = strings.TrimSpace(fullRet[:idx])
								reqs = strings.TrimSpace(fullRet[idx+1:])
							}
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

func zigType(specType string) string {
	switch strings.ToLower(strings.TrimSpace(specType)) {
	case "string":
		return "[]const u8"
	case "int", "integer":
		return "i64"
	case "rexx":
		return "f64"
	case "double", "real", "float", "float64":
		return "f64"
	case "boolean", "bool":
		return "bool"
	case "usize":
		return "usize"
	case "[1024]float64", "[1024]f64":
		return "[1024]f64"
	default:
		return specType
	}
}

func GenerateZigSkeleton(spec ParsedSpec, mainClassName string) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("const std = @import(\"std\");\n\n"))

	// DTO structs
	for _, class := range spec.Classes {
		if len(class.Methods) == 0 && class.Name != mainClassName {
			sb.WriteString(fmt.Sprintf("pub const %s = struct {\n", class.Name))
			for _, field := range class.Fields {
				parts := strings.Split(field, "(")
				fName := strings.TrimSpace(parts[0])
				if fName == "" {
					continue
				}
				fType := "[]const u8"
				if len(parts) > 1 {
					fType = zigType(strings.TrimSpace(strings.Trim(parts[1], "()")))
				}
				sb.WriteString(fmt.Sprintf("\t%s: %s,\n", fName, fType))
			}
			sb.WriteString("};\n\n")
		}
	}

	// Main struct with fields inherited from DTO
	dtoFields := ""
	for _, class := range spec.Classes {
		if len(class.Methods) == 0 && class.Name == mainClassName {
			for _, field := range class.Fields {
				parts := strings.Split(field, "(")
				fName := strings.TrimSpace(parts[0])
				if fName == "" {
					continue
				}
				fType := "[]const u8"
				if len(parts) > 1 {
					fType = zigType(strings.TrimSpace(strings.Trim(parts[1], "()")))
				}
				fNameUpper := strings.ToUpper(fName[:1]) + fName[1:]
		dtoFields += fmt.Sprintf("\t%s: %s,\n", fNameUpper, fType)
			}
			break
		}
	}
	if dtoFields != "" {
		sb.WriteString(fmt.Sprintf("pub const %s = struct {\n%s\n\tconst Self = @This();\n\n", mainClassName, dtoFields))
	} else {
		sb.WriteString(fmt.Sprintf("pub const %s = struct {\n\tconst Self = @This();\n\n", mainClassName))
	}

	// Method stubs
	for _, class := range spec.Classes {
		if len(class.Methods) > 0 && class.Name == mainClassName {
			for _, m := range class.Methods {
				sb.WriteString(fmt.Sprintf("\t// SKELETON_%s\n", m.Name))
			}
		}
	}

	sb.WriteString("};\n")
	return sb.String(), nil
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

	// Determine if the spec requires database imports
	needsDB := false
	for _, inv := range spec.Invariants {
		lower := strings.ToLower(inv)
		if strings.Contains(lower, "sqlite") {
			needsDB = true
			break
		}
	}

	// Read base imports from unified_exemplars ledger
	baseImports := ""
	if needsDB && db != nil {
		db.QueryRow("SELECT default_imports FROM unified_exemplars WHERE exemplar_id='NETREXX_GRAMMAR_BASICS'").Scan(&baseImports)
	}
	if baseImports == "" {
		baseImports = ""
	}
	// For non-DB specs, use minimal imports. For DB specs, DB import is in baseImports.
	if !needsDB {
		baseImports = ""
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

	// DTO Structs (data-only classes) — skip if main class inherits its fields
	for _, class := range spec.Classes {
		if len(class.Methods) == 0 && class.Name != mainClassName {
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

	// Main struct with method stubs — inherits fields from DTO class if exists
	dtoFields := ""
	for _, class := range spec.Classes {
		if len(class.Methods) == 0 && class.Name == mainClassName {
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
				fNameUpper := strings.ToUpper(fName[:1]) + fName[1:]
				dtoFields += fmt.Sprintf("\t%s %s\n", fNameUpper, fType)
			}
			break
		}
	}
	if dtoFields != "" {
		sb.WriteString(fmt.Sprintf("\ntype %s struct {\n%s}\n", mainClassName, dtoFields))
	} else if needsDB {
		sb.WriteString(fmt.Sprintf("\ntype %s struct {\n\tdb *sql.DB\n}\n", mainClassName))
	} else {
		sb.WriteString(fmt.Sprintf("\ntype %s struct{}\n", mainClassName))
	}
	for _, class := range spec.Classes {
		if len(class.Methods) > 0 && class.Name == mainClassName {
			for _, m := range class.Methods {
				// main() lives in its own file (main.go) with separate imports
				if m.Name == "main" {
					continue
				}
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
	case "[1024]float64":
		return "[1024]float64"
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
		} else if strings.Contains(p, " ") && !strings.Contains(p, "[") {
			// Space-separated name type: "value float64"
			sub := strings.SplitN(p, " ", 2)
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

type CoverageGap struct {
	startLine  int
	endLine    int
	stmtCount  int
	sourceText string
}

func queryCoverageGaps(dbPath string, methodName string) []CoverageGap {
	var gaps []CoverageGap
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return gaps
	}
	defer db.Close()

	rows, err := db.Query("SELECT start_line, end_line, statement_count, unexecuted_source_text FROM generated_coverage_gaps WHERE method_name=? ORDER BY start_line", methodName)
	if err != nil {
		return gaps
	}
	defer rows.Close()

	for rows.Next() {
		var g CoverageGap
		var src sql.NullString
		if err := rows.Scan(&g.startLine, &g.endLine, &g.stmtCount, &src); err == nil {
			if src.Valid {
				g.sourceText = src.String
			}
			gaps = append(gaps, g)
		}
	}
	return gaps
}

func convertArgsToZig(args string) string {
	if args == "" {
		return ""
	}
	parts := strings.Split(args, ",")
	var zigParts []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// Handle "name: Type" or "name = Type" or "name Type" or just "name"
		if strings.Contains(p, ":") {
			sub := strings.SplitN(p, ":", 2)
			name := strings.TrimSpace(sub[0])
			typ := zigType(strings.TrimSpace(sub[1]))
			zigParts = append(zigParts, fmt.Sprintf("%s: %s", name, typ))
		} else if strings.Contains(p, "=") {
			sub := strings.SplitN(p, "=", 2)
			name := strings.TrimSpace(sub[0])
			typ := zigType(strings.TrimSpace(sub[1]))
			zigParts = append(zigParts, fmt.Sprintf("%s: %s", name, typ))
		} else if strings.Contains(p, " ") {
			sub := strings.SplitN(p, " ", 2)
			name := strings.TrimSpace(sub[0])
			typ := zigType(strings.TrimSpace(sub[1]))
			zigParts = append(zigParts, fmt.Sprintf("%s: %s", name, typ))
		} else {
			zigParts = append(zigParts, fmt.Sprintf("%s: void", strings.TrimSpace(p)))
		}
	}
	return strings.Join(zigParts, ", ")
}

func BuildMethodPrompt(dbPath string, mainClassName string, method SpecMethod, invariants []string, report string, classes []SpecClass) (string, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open context database: %w", err)
	}
	defer db.Close()

	// Determine if spec needs database operations
	needsDB := false
	targetLang := "go"
	for _, inv := range invariants {
		if strings.Contains(strings.ToLower(inv), "sqlite") {
			needsDB = true
		}
		if strings.Contains(strings.ToLower(inv), "zig") {
			targetLang = "zig"
		}
	}

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

	// Layer 3: Exemplar Blocks — adapts to domain (DB vs algorithmic)
	if needsDB {
		prompt.WriteString("### LAYER 3: RELATIONAL DATABASE EXEMPLARS\n")
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
	} else {
		prompt.WriteString("### LAYER 3: ALGORITHMIC/MEMORY EXEMPLARS\n")
		scope := "%Memory%"
		if targetLang == "zig" {
			scope = "%Memory.Zig%"
		}
		rows, err := db.Query("SELECT few_shot_prompt_block FROM unified_exemplars WHERE domain_scope LIKE ?", scope)
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
	}

	// Layer 4: Specific Method Target Requirements
	prompt.WriteString("### LAYER 4: TARGET METHOD REQUIREMENTS\n")
	prompt.WriteString(fmt.Sprintf("Package: %s\n", mainClassName))
	
	goArgs := convertArgsToGo(method.Args)
	goRet := method.Returns
	if goRet != "" {
		// Map NetRexx return types to Go types
		goRet = goType(strings.TrimPrefix(goRet, "returns "))
		goRet = " " + goRet
	} else if strings.HasPrefix(method.Name, "New") {
		// Constructor pattern: return pointer to self
		goRet = " *" + mainClassName
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
		// Inject exact method name inventory to prevent casing mismatches
		prompt.WriteString("AVAILABLE METHODS (case-sensitive, use EXACTLY as listed):\n")
		for _, class := range classes {
			if len(class.Methods) > 0 && class.Name == mainClassName {
				for _, m := range class.Methods {
					prompt.WriteString(fmt.Sprintf("  - %s\n", m.Name))
				}
			}
		}
		prompt.WriteString("Call these methods in lowercase (e.g., router.initRoutingTable(...), NOT router.InitRoutingTable(...)).\n")
	} else if targetLang == "zig" {
		// Zig function signature: pub fn methodName(self: *Self, args) returnType
		zigRet := goRet
		if zigRet == "" {
			zigRet = " void"
		}
		zigArgs := convertArgsToZig(method.Args)
		// Static-like methods (init constructors) don't take self
		if method.Name == "init" || method.Name == "NewRingBuffer" || method.Name == "new" {
			if zigArgs == "" {
				sig = fmt.Sprintf("pub fn %s()%s", method.Name, zigRet)
			} else {
				sig = fmt.Sprintf("pub fn %s(%s)%s", method.Name, zigArgs, zigRet)
			}
		} else {
			if zigArgs == "" {
				sig = fmt.Sprintf("pub fn %s(self: *Self)%s", method.Name, zigRet)
			} else {
				sig = fmt.Sprintf("pub fn %s(self: *Self, %s)%s", method.Name, zigArgs, zigRet)
			}
		}
		prompt.WriteString(fmt.Sprintf("Target Signature: %s\n", sig))
		prompt.WriteString("Zig CONVENTIONS:\n")
		prompt.WriteString("- Use 'pub fn' for public functions\n")
		prompt.WriteString("- Init/constructor functions take NO self parameter\n")
		prompt.WriteString("- Other methods use 'self: *Self' or 'self: *const Self'\n")
		prompt.WriteString("- Use 'void' for functions with no return value\n")
		prompt.WriteString("- Return Zig error unions for fallible operations: `!ReturnType`\n")
		prompt.WriteString("- Use `const Self = @This();` pattern\n")
		prompt.WriteString("ZIG 0.16.0 SYNTAX RULES (ABSOLUTE ENFORCEMENT):\n")
		prompt.WriteString("- For integer-to-float coercion: use @as(f64, @floatFromInt(value)) or just @as(f64, value).\n")
	prompt.WriteString("- Do NOT use @intToFloat or standalone @floatFromInt — both fail in Zig 0.16 without context.\n")
		prompt.WriteString("- For integer casts: use @intCast(target_type, value).\n")
		prompt.WriteString("- Discard unused parameters with `_ = param;`.\n")
		prompt.WriteString("- Use `while (cond) : (continue_expr) { }` for loops, not C-style for.\n")
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

	// Inject coverage gaps for this method if any exist
	coverageGaps := queryCoverageGaps(dbPath, method.Name)
	if len(coverageGaps) > 0 {
		prompt.WriteString("### COVERAGE GAPS (unexecuted code blocks from fuzzing):\n")
		for _, gap := range coverageGaps {
			prompt.WriteString(fmt.Sprintf("Lines %d-%d (%d statements) — NEVER EXECUTED:\n", gap.startLine, gap.endLine, gap.stmtCount))
			if gap.sourceText != "" {
				prompt.WriteString(gap.sourceText + "\n")
			}
		}
		prompt.WriteString("REFACTORING INVARIANT: Ensure these blocks are reachable under test, or remove dead code.\n\n")
	}

	prompt.WriteString("VALIDATION & ERROR HANDLING:\n")
	if needsDB {
		prompt.WriteString("You MUST validate all input arguments at the beginning of the function (before any database operation):\n")
		prompt.WriteString("1. If a string parameter is empty or whitespace-only, return an error: `return 0, fmt.Errorf(\"invalid input: %%s is empty\", name)`\n")
		prompt.WriteString("2. If a dbPath parameter contains path traversal (\"..\", \"/etc/\"), return `fmt.Errorf(\"path traversal blocked\")`\n")
		prompt.WriteString("3. If a string parameter contains SQL injection patterns, return `fmt.Errorf(\"sql injection blocked\")`\n")
		prompt.WriteString("4. Use Go-style error returns, not Java exceptions. Return `(result, error)` where error is nil on success.\n")
		prompt.WriteString("Ensure these checks are at the very beginning of each method body and are NOT caught by any internal try-catch blocks that handle database operations.\n\n")
		prompt.WriteString("IMPORTANT: You are generating Go code. Output ONLY the complete Go function body for '" + sig + "'. Do not include the enclosing struct definition, package, or imports. Do not wrap in markdown code blocks. Use standard Go: `func (s *Type) Name(args) (returnType, error) { ... }` with proper error handling. Use `database/sql` for queries. Always return an error as the last return value.\n")
	} else {
		prompt.WriteString("CONTRACT BOUNDARY ENFORCEMENT:\n")
		prompt.WriteString("This is a pure-memory computational struct. Do NOT implement streaming interfaces (io.Writer, io.Reader) or database operations.\n")
		prompt.WriteString("Use ONLY the struct fields defined in the skeleton. Access internal state directly via pointer receiver (r *RingBuffer).\n")
		prompt.WriteString("Output ONLY the complete Go function body for '" + sig + "'. Do not include the enclosing struct definition, package, or imports. Do not wrap in markdown code blocks.\n")
	}
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

	// Derive main class name from filename, then verify against parsed classes
	mainClassName := strings.TrimSuffix(filepath.Base(specPath), "Spec.md")
	for _, class := range spec.Classes {
		if len(class.Methods) > 0 {
			mainClassName = class.Name
			break
		}
	}

	// Detect target language from invariants
	targetLang := "go"
	for _, inv := range spec.Invariants {
		if strings.Contains(strings.ToLower(inv), "zig") {
			targetLang = "zig"
			break
		}
	}

	// Emit class skeleton immediately to disk
	var skeleton string
	var skeletonFile string
	if targetLang == "zig" {
		skeleton, err = GenerateZigSkeleton(spec, mainClassName)
		skeletonFile = filepath.Join(projectDir, "generated", mainClassName+".zig")
	} else {
		skeleton, err = GenerateGoSkeleton(dbPath, spec, mainClassName)
		skeletonFile = filepath.Join(projectDir, "generated", mainClassName+".go")
	}
	if err != nil {
		fmt.Printf("Error generating skeleton: %v\n", err)
		os.Exit(1)
	}
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
			if len(class.Methods) > 0 {
				break  // prefer the class with methods
			}
		}
	}
	if targetClass != nil && len(targetClass.Methods) == 0 {
		// Fallback: find the class with methods even if name doesn't match exactly
		for _, class := range spec.Classes {
			if len(class.Methods) > 0 {
				targetClass = &class
				mainClassName = class.Name
				break
			}
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

	// Generate type-driven fuzz seeds from DTO class fields
	seedsFile := filepath.Join(projectDir, "generated", "seeds_test.go")
	var seedsSb strings.Builder
	seedsSb.WriteString("package comfactoryrouting\n\nimport \"math\"\n\n")
	seedsSb.WriteString("// Auto-generated seed corpus via type-driven boundary induction.\n")
	seedsSb.WriteString("// Each field is tested with empty, path-traversal, SQL-injection, and boundary values.\n\n")
	seedsSb.WriteString("// SeedCorpusData returns boundary test values for each field type.\n")
	seedsSb.WriteString("// Used by the fuzz harness to construct targeted corpus entries.\n")
	seedsSb.WriteString("func SeedCorpusData() map[string][]interface{} {\n")
	seedsSb.WriteString("\treturn map[string][]interface{}{\n")
	for _, class := range spec.Classes {
		if len(class.Methods) == 0 {
			for _, field := range class.Fields {
				parts := strings.Split(field, "(")
				if len(parts) < 1 {
					continue
				}
				fName := strings.TrimSpace(parts[0])
				fType := "string"
				if len(parts) > 1 {
					fType = goType(strings.TrimSpace(strings.Trim(parts[1], "()")))
				}
				seedsSb.WriteString(fmt.Sprintf("\t\t// %s (%s)\n", fName, fType))
				seedsSb.WriteString(fmt.Sprintf("\t\t\"%s_empty\": {},\n", fName))
				switch fType {
				case "string":
					seedsSb.WriteString(fmt.Sprintf("\t\t\"%s_path_traversal\": {\"../../etc/passwd\"},\n", fName))
					seedsSb.WriteString(fmt.Sprintf("\t\t\"%s_sql_injection\": {\"'; DROP TABLE --\"},\n", fName))
				case "float64":
					seedsSb.WriteString(fmt.Sprintf("\t\t\"%s_negative\": {-1.0},\n", fName))
					seedsSb.WriteString(fmt.Sprintf("\t\t\"%s_nan\": {math.NaN()},\n", fName))
				}
			}
		}
	}
	seedsSb.WriteString("\t}\n")
	seedsSb.WriteString("}\n")
	os.WriteFile(seedsFile, []byte(seedsSb.String()), 0644)
	fmt.Printf("[INFO] Fuzz seeds written to: %s\n", seedsFile)

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
