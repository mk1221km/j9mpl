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

func GenerateNrxSkeleton(dbPath string, spec ParsedSpec, mainClassName string) (string, error) {
	var sb strings.Builder
	
	// Resolve package name from DB
	packageName := "com.factory" // default fallback
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
					packageName = strings.Join(subParts[:len(subParts)-1], ".")
				}
			}
		}
	}

	sb.WriteString(fmt.Sprintf("package %s\n", packageName))
	sb.WriteString("options binary\n")
	sb.WriteString("import java.sql.DriverManager\n")
	sb.WriteString("import java.sql.Connection\n")
	sb.WriteString("import java.sql.Statement\n")
	sb.WriteString("import java.sql.PreparedStatement\n")
	sb.WriteString("import java.sql.ResultSet\n")
	sb.WriteString("import java.sql.SQLException\n\n")

	sb.WriteString(fmt.Sprintf("class %sDummy private\n\n", mainClassName))

	// DTO Classes
	for _, class := range spec.Classes {
		if len(class.Methods) == 0 {
			sb.WriteString(fmt.Sprintf("class %s shared\n", class.Name))
			sb.WriteString("  properties public\n")
			for _, field := range class.Fields {
				parts := strings.Split(field, "(")
				fName := strings.TrimSpace(parts[0])
				fType := "Rexx"
				if len(parts) > 1 {
					fType = strings.TrimSpace(strings.Trim(parts[1], "()"))
				}
				sb.WriteString(fmt.Sprintf("    %s = %s\n", fName, fType))
			}
			sb.WriteString("\n")
		}
	}

	// Main public class
	sb.WriteString(fmt.Sprintf("class %s public\n", mainClassName))
	for _, class := range spec.Classes {
		if len(class.Methods) > 0 && class.Name == mainClassName {
			for _, m := range class.Methods {
				nrxArgs := convertArgsToNetRexx(m.Args)
				retClause := ""
				retType := "void"
				if m.Returns != "" {
					retClause = " " + m.Returns
					retType = strings.TrimSpace(strings.TrimPrefix(m.Returns, "returns"))
				}

				sb.WriteString(fmt.Sprintf("  method %s(%s) public static%s\n", m.Name, nrxArgs, retClause))
				sb.WriteString(fmt.Sprintf("    -- SKELETON_%s\n", m.Name))
				if retType == "void" || m.Name == "main" {
					sb.WriteString("    nop\n\n")
				} else {
					sb.WriteString(fmt.Sprintf("    return %s null\n\n", retType))
				}
			}
		}
	}

	return sb.String(), nil
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
	rows, err := db.Query("SELECT few_shot_prompt_block FROM unified_exemplars WHERE domain_scope = 'Database.SQLite'")
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

	// Layer 3.5: Local Machine-Verified Reference Implementations (Few-Shots)
	prompt.WriteString("### LAYER 3.5: LOCAL MACHINE-VERIFIED REFERENCE IMPLEMENTATIONS (FEW-SHOTS)\n")
	prompt.WriteString("The following NetRexx implementations have been machine-verified as correct, secure, and compliant with all invariants. Use them as reference patterns:\n\n")
	implRows, err := db.Query("SELECT exemplar_id, fact_context_predicate, few_shot_prompt_block FROM unified_exemplars WHERE domain_scope = 'Implementation.NetRexx'")
	if err == nil {
		defer implRows.Close()
		for implRows.Next() {
			var exID, predicate, snippet string
			if err := implRows.Scan(&exID, &predicate, &snippet); err == nil {
				prompt.WriteString(fmt.Sprintf("Exemplar: %s\n", exID))
				prompt.WriteString(fmt.Sprintf("Specification Context: %s\n", predicate))
				prompt.WriteString("Implementation:\n")
				prompt.WriteString(snippet)
				prompt.WriteString("\n\n")
			}
		}
	}

	// Layer 4: Specific Method Target Requirements
	prompt.WriteString("### LAYER 4: TARGET METHOD REQUIREMENTS\n")
	prompt.WriteString(fmt.Sprintf("Class: %s\n", mainClassName))
	
	nrxArgs := convertArgsToNetRexx(method.Args)
	retClause := ""
	if method.Returns != "" {
		retClause = " " + method.Returns
	}
	sig := fmt.Sprintf("method %s(%s) public static%s", method.Name, nrxArgs, retClause)
	prompt.WriteString(fmt.Sprintf("Target Method Signature: %s\n", sig))
	prompt.WriteString(fmt.Sprintf("Requirements: %s\n\n", method.Requirements))

	prompt.WriteString("Architectural Invariants & Database Schema:\n")
	for _, inv := range invariants {
		prompt.WriteString(fmt.Sprintf("- %s\n", inv))
	}
	prompt.WriteString("\n")

	prompt.WriteString("VALIDATION & EXCEPTION INVARIANTS:\n")
	prompt.WriteString("You MUST validate all input arguments and their fields at the very beginning of the method (before any database operation or do-catch block) and throw/signal the expected exception types for invalid or boundary inputs:\n")
	prompt.WriteString("1. If a String parameter (like dbPath, sender, receiver, etc.), or a record argument itself (like record), OR any String field of a record argument (e.g. record.timestamp, record.metricName, record.txId, record.sender, record.receiver, record.priority) is null, empty (\"\"), or contains only whitespace, you MUST throw/signal `java.lang.IllegalArgumentException` using the `signal` statement (e.g. `signal java.lang.IllegalArgumentException(\"Invalid input\")`).\n")
	prompt.WriteString("2. If a String parameter, OR any String field of a record argument, represents a path traversal attempt (e.g. starts with \"/etc/\", contains \"..\", or contains \"C:\\\\Windows\" / starts with \"C:\\\\\"), you MUST throw/signal `java.io.IOException` using `signal java.io.IOException(\"Path traversal blocked\")`.\n")
	prompt.WriteString("3. If a String parameter, OR any String field of a record argument, contains a SQL injection attempt (e.g. contains \"' OR '1'='1\", \"; DROP TABLE\", or \"' UNION SELECT\"), you MUST throw/signal `java.sql.SQLException` using `signal java.sql.SQLException(\"SQL Injection blocked\")`.\n")
	prompt.WriteString("4. If a parameter, OR any field of a record argument, needs to be parsed as a number (like amount, metricValue, etc.), and the string/Rexx value is null or is not a valid number (e.g. \"NaN\", \"Infinity\", \"-Infinity\", overflows Double range, or overflows Constant/Integer range for integer parameters), you MUST throw/signal `java.lang.NumberFormatException` using `signal java.lang.NumberFormatException(\"Invalid number format\")`.\n")
	prompt.WriteString("Ensure these checks are at the very beginning of each method body and are NOT caught by any internal try-catch blocks that handle database operations.\n\n")

 	prompt.WriteString("IMPORTANT INSTRUCTIONS FOR NETREXX DIALECT:\n")
	prompt.WriteString("1. Variable declarations MUST follow the NetRexx syntax: 'varName = Type initialValue' (e.g. 'dbPath = String null', 'avg = Rexx 0', 'count = int 0'). Do NOT use Java-style declarations like 'Type varName = value' or 'String dbPath = null' as they will cause syntax errors.\n")
	prompt.WriteString("2. NetRexx methods do NOT have a terminating 'end' keyword at the method level. Only inner blocks like 'do', 'loop', and 'select' should be closed with 'end'. Do NOT append a trailing 'end' at the end of the method body.\n")
	prompt.WriteString("3. Checked exceptions (like Exception, SQLException) can ONLY be caught inside a 'do ... catch' block if the body of that 'do' block calls a method that is explicitly declared to throw/signal that exception. If no such method is called, catching checked exceptions is a compile-time error. For 'main', do not catch checked exceptions, or just catch 'RuntimeException' / 'Throwable', or avoid catch blocks entirely.\n")
	prompt.WriteString("4. You MUST guard all database connection logic against null or empty/placeholder dbPath values. Since this method takes a 'dbPath = String' parameter, you MUST check if it is valid. To ensure short-circuiting and avoid NullPointerExceptions, use nested checks: `if dbPath \\== null then if dbPath \\== \"null\" then do` before connecting via JDBC, otherwise SQLite JDBC will physically create a database file named 'null' in the current working directory. The validation checks above should be run first.\n")
	prompt.WriteString("5. Mismatched block and catch syntax: (a) Never write 'finally do' on the same line. If you need a try-catch block inside a finally clause, place 'finally' on its own line and start the 'do' block on the next line. (b) Every catch clause MUST follow NetRexx syntax: 'catch ex = ExceptionType' (e.g. 'catch ex = SQLException' or 'catch ex = Exception'). Never write 'catch Exception ex' or 'catch ex' without a type.\n")
	prompt.WriteString("6. Numeric Types and Literals: Literals with decimal points (e.g. 0.0) are treated as float by default. To declare a primitive double, use `varName = double 0` or cast it like `varName = double 0.0`. Prefer primitive `double` and `int` over their boxed object wrapper classes `java.lang.Double` or `java.lang.Integer`. To check if a double is infinite or NaN, use `Double.isInfinite(val)` or `Double.isNaN(val)` instead of calling methods on a wrapper object.\n")
	prompt.WriteString("7. Backslash in String Literals: To write a literal backslash inside a NetRexx string literal, you MUST double it (e.g. write 'C:\\\\\\\\Windows' or '\\\\\\\\' in code to represent a backslash). A single backslash followed by a letter (like '\\\\W') is an invalid escape sequence and will cause translation to fail.\n")
	prompt.WriteString("8. Object Instantiation: Do NOT use the `new` keyword to instantiate classes. In NetRexx, you instantiate a class by calling its constructor directly (e.g. use `record = MetricRecord()` or `record = MetricRecord(\"arg1\", \"arg2\")` instead of `record = new MetricRecord()`).\n")
	prompt.WriteString("9. Accessing DTO fields: DTO helper classes (like `MetricRecord` or `TransactionRecord`) expose public properties directly. You MUST access them directly by name without Java-style getter/setter methods (e.g. use `record.txId` or `record.amount` instead of `record.getId()` or `record.getAmount()`).\n")
	prompt.WriteString("10. Null Check Reference Comparison: ALWAYS use the strict identity operators `== null` and `\\== null` (double equals) when checking for null references. Do NOT use `=` or `\\=` for null checking against Java null, as value comparisons against null can throw NullPointerExceptions.\n")
	prompt.WriteString("11. Non-Short-Circuiting Operators: NetRexx logical operators `|` and `&` do NOT short-circuit (they always evaluate both sides). Therefore, you MUST NOT combine null-checks and member/method calls in a single expression (e.g. do NOT write `if record == null | record.timestamp == null`). Instead, check if the object itself is null first, and only check its fields or call its methods in separate, subsequent statements.\n")
	prompt.WriteString("12. String methods vs Rexx methods: Properties or parameters declared as `String` are Java `java.lang.String` objects. You MUST use Java String methods (e.g., `val.trim().length() == 0` for empty/whitespace checks, `val.indexOf(substring) >= 0` to check for containment, `val.startsWith(prefix)` to check prefix) on them. Do NOT use Rexx-specific methods like `.pos()`, `.datatype()`, `.strip()`, `.left()`, `.stripspace()`, or `.upper()` on `String` types. If you want to use Rexx methods, you must first cast/convert them to Rexx (e.g., `Rexx(val)`).\n")
	prompt.WriteString("13. Loops: In NetRexx, loops MUST use the `loop` keyword, not `do` (e.g. use `loop pattern over patterns` instead of `do for pattern over patterns` or `do pattern over patterns`). Closing a loop block must be done with `end`.\n")
	prompt.WriteString("14. Unique exception names in catch blocks: In NetRexx, exception variables in catch blocks have method scope. If a method contains multiple catch blocks, you MUST use different variable names (e.g., `catch exNum = NumberFormatException`, `catch exSql = SQLException`) to avoid compiler errors about type mismatch or duplicate declarations.\n")
	prompt.WriteString("Output ONLY the complete NetRexx method block starting with '" + sig + "'. Do not include the enclosing class, package, or imports. Do not wrap in markdown code blocks.\n")

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
	skeleton, err := GenerateNrxSkeleton(dbPath, spec, mainClassName)
	if err != nil {
		fmt.Printf("Error generating skeleton: %v\n", err)
		os.Exit(1)
	}

	skeletonFile := filepath.Join(projectDir, "generated", mainClassName+".nrx")
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
