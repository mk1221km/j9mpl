package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// SymbolDetails holds parsed logical details from an M3 URI
type SymbolDetails struct {
	Type string // "class", "method", "constructor", "field", "unknown"
	Name string // The logical symbol name
}

// ParseSymbolURI extracts type and name from logical M3 URI.
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

	parts := strings.Split(cleaned, "///")
	if len(parts) < 2 {
		return SymbolDetails{Type: symbolType, Name: ""}
	}
	
	pathPart := parts[len(parts)-1]
	subParts := strings.Split(pathPart, "/")
	namePart := subParts[len(subParts)-1]
	
	if idx := strings.Index(namePart, "("); idx != -1 {
		namePart = namePart[:idx]
	}

	return SymbolDetails{Type: symbolType, Name: namePart}
}

// ResolveContext extracts NetRexx block from DB using Symbol URI
func ResolveContext(dbPath, symbolURI string) (string, error) {
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

	nrxPath := strings.Replace(javaPath, ".java", ".nrx", 1)
	details := ParseSymbolURI(symbolURI)
	if details.Type == "unknown" || details.Name == "" {
		content, err := os.ReadFile(nrxPath)
		if err != nil {
			return "", fmt.Errorf("failed to read full source fallback file: %w", err)
		}
		return string(content), nil
	}

	return extractBlockFromNrx(nrxPath, details.Type, details.Name)
}

func extractBlockFromNrx(nrxPath, detailsType, detailsName string) (string, error) {
	file, err := os.Open(nrxPath)
	if err != nil {
		return "", fmt.Errorf("failed to open NetRexx source file: %w", err)
	}
	defer file.Close()

	classRegex := regexp.MustCompile(`(?i)^\s*class\s+(\w+)`)
	methodRegex := regexp.MustCompile(`(?i)^\s*method\s+(\w+)`)

	var extractedLines []string
	inBlock := false
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		
		if !inBlock {
			if detailsType == "class" {
				if matches := classRegex.FindStringSubmatch(line); len(matches) > 1 {
					if strings.EqualFold(matches[1], detailsName) {
						inBlock = true
						extractedLines = append(extractedLines, line)
					}
				}
			} else if detailsType == "method" || detailsType == "constructor" {
				if matches := methodRegex.FindStringSubmatch(line); len(matches) > 1 {
					if strings.EqualFold(matches[1], detailsName) {
						inBlock = true
						extractedLines = append(extractedLines, line)
					}
				}
			}
		} else {
			isClass := classRegex.MatchString(line)
			isMethod := methodRegex.MatchString(line)

			if isClass || isMethod {
				if detailsType == "method" || detailsType == "constructor" || isClass {
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
		return "", fmt.Errorf("symbol '%s' of type '%s' not found in file %s", detailsName, detailsType, nrxPath)
	}

	return strings.Join(extractedLines, "\n"), nil
}

func findEnclosingBlock(nrxPath string, errLine int) (blockType, blockName string, startLine, endLine int) {
	file, err := os.Open(nrxPath)
	if err != nil {
		return "", "", 0, 0
	}
	defer file.Close()

	classRegex := regexp.MustCompile(`(?i)^\s*class\s+(\w+)`)
	methodRegex := regexp.MustCompile(`(?i)^\s*method\s+(\w+)`)

	var lines []string
	scanner := bufio.NewScanner(file)
	lineIdx := 0

	type blockBoundary struct {
		name      string
		kind      string // "class" or "method"
		startLine int
	}
	var boundaries []blockBoundary

	for scanner.Scan() {
		lineIdx++
		line := scanner.Text()
		lines = append(lines, line)

		if matches := classRegex.FindStringSubmatch(line); len(matches) > 1 {
			boundaries = append(boundaries, blockBoundary{name: matches[1], kind: "class", startLine: lineIdx})
		} else if matches := methodRegex.FindStringSubmatch(line); len(matches) > 1 {
			boundaries = append(boundaries, blockBoundary{name: matches[1], kind: "method", startLine: lineIdx})
		}
	}

	var activeBoundary blockBoundary
	nextStart := lineIdx + 1

	for i, b := range boundaries {
		if errLine >= b.startLine {
			activeBoundary = b
			if i+1 < len(boundaries) {
				nextStart = boundaries[i+1].startLine
			} else {
				nextStart = lineIdx + 1
			}
		}
	}

	if activeBoundary.name == "" {
		return "file", "", 1, lineIdx
	}

	return activeBoundary.kind, activeBoundary.name, activeBoundary.startLine, nextStart - 1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage (ECJ mode): ./self_correct <java_file_path> [ecj_compiler_args...]")
		fmt.Println("Usage (NetRexx mode): ./self_correct <nrx_file_path>")
		os.Exit(1)
	}

	targetPath := os.Args[1]

	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error resolving executable path: %v\n", err)
		os.Exit(1)
	}
	
	projectDir := filepath.Dir(filepath.Dir(execPath))
	dbPath := filepath.Join(projectDir, ".context", "project_context.db")
	ecjPath := filepath.Join(projectDir, "bin", "ecj")
	nrcPath := filepath.Join(projectDir, "bin", "nrc")

	// Check if we are starting from a NetRexx file directly
	if strings.HasSuffix(targetPath, ".nrx") {
		nrxPathClean := targetPath
		javaPathClean := strings.Replace(nrxPathClean, ".nrx", ".java", 1)

		fmt.Printf("[INFO] Running in NetRexx Translation + ECJ Compilation mode for %s\n", nrxPathClean)
		fmt.Printf("[DEBUG] projectDir=%s nrcPath=%s ecjPath=%s javaPathClean=%s\n", projectDir, nrcPath, ecjPath, javaPathClean)

		// 1. Run nrc translation
		nrcCmd := exec.Command(nrcPath, "-nocompile", "-keepasjava", "-replace", "-format", nrxPathClean)
		var nrcOut bytes.Buffer
		nrcCmd.Stdout = &nrcOut
		nrcCmd.Stderr = &nrcOut

		nrcErr := nrcCmd.Run()
		nrcOutputStr := strings.ReplaceAll(nrcOut.String(), "\r\n", "\n")
		ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
		nrcOutputStr = ansiRegex.ReplaceAllString(nrcOutputStr, "")

		// Regular expression to parse NetRexx errors:
		// 64 +++       pstmt.setString(1, invalidVarName)
		//    +++                          ^^^^^^^^^^^^^^
		//    +++ Error: Unknown variable
		netrexxErrorRegex := regexp.MustCompile(`(?m)^\s*(\d+)\s*\+\+\+\s*(.*?)\n\s*\+\+\+\s*([-^]+)\s*\n\s*\+\+\+\s*((?:Error|Warning):.*)`)
		nrcMatches := netrexxErrorRegex.FindStringSubmatch(nrcOutputStr)

		if nrcErr != nil || nrcMatches != nil {
			if nrcMatches == nil {
				fmt.Printf("[ERROR] NetRexx translation failed, but errors could not be parsed. Output:\n%s\n", nrcOutputStr)
				os.Exit(1)
			}

			errLine, _ := strconv.Atoi(nrcMatches[1])
			errSource := nrcMatches[2]
			errMarker := nrcMatches[3]
			errMessage := nrcMatches[4]

			fmt.Printf("[INFO] Intercepted NetRexx Translator Error at %s (line %d)\n", nrxPathClean, errLine)

			// Find enclosing method/class in .nrx file
			bKind, bName, startLine, endLine := findEnclosingBlock(nrxPathClean, errLine)
			symbolContext := fmt.Sprintf("%s %s (lines %d-%d)", bKind, bName, startLine, endLine)

			// Extract source block
			nrxLines, err := readLines(nrxPathClean)
			if err != nil {
				fmt.Printf("[ERROR] Failed to read NetRexx source: %v\n", err)
				os.Exit(1)
			}
			
			var blockLines []string
			for i := startLine - 1; i < endLine && i < len(nrxLines); i++ {
				blockLines = append(blockLines, nrxLines[i])
			}
			sourceContext := strings.Join(blockLines, "\n")

			// Format diagnostic
			errMsg := fmt.Sprintf("line %d: %s\nContext:\n%s\n%s", errLine, errMessage, errSource, errMarker)

			// Package prompt
			generateAndWritePrompt(projectDir, symbolContext, errMsg, sourceContext)
			os.Exit(1)
		}

		// If translation succeeded, we run ECJ on the translated Java file
		fmt.Println("[INFO] Translation succeeded. Executing ECJ compiler validation...")

		// Verify the .java file exists where expected; search fallback if not
		origJavaPath := javaPathClean
		if _, err := os.Stat(javaPathClean); os.IsNotExist(err) {
			found := false
			// Search common output locations
			searchDirs := []string{
				".",                                          // CWD
				projectDir,                                   // project root
				filepath.Dir(nrxPathClean),                   // .nrx source directory
				filepath.Join(projectDir, "generated"),       // generated/
				filepath.Join(projectDir, "bin"),             // bin/
			}
			for _, dir := range searchDirs {
				altPath := filepath.Join(dir, filepath.Base(javaPathClean))
				if _, err2 := os.Stat(altPath); err2 == nil {
					javaPathClean = altPath
					found = true
					break
				}
			}
			// Walk the project to find it as last resort
			if !found {
				filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
					if err != nil || found {
						return filepath.SkipDir
					}
					if info.Name() == filepath.Base(javaPathClean) {
						javaPathClean = path
						found = true
					}
					return nil
				})
			}
		}
		if origJavaPath != javaPathClean {
			fmt.Printf("[DEBUG] Java path resolved via fallback: %s -> %s\n", origJavaPath, javaPathClean)
		}
		
		// Run ECJ with classpaths
		cmdArgs := []string{"-17", "-proceedOnError", "-d", filepath.Join(projectDir, "bin")}
		// Add any extra arguments passed after targetPath
		if len(os.Args) > 2 {
			cmdArgs = append(cmdArgs, os.Args[2:]...)
		}
		cmdArgs = append(cmdArgs, javaPathClean)
		
		ecjCmd := exec.Command(ecjPath, cmdArgs...)
		var ecjOut, ecjErr bytes.Buffer
		ecjCmd.Stdout = &ecjOut
		ecjCmd.Stderr = &ecjErr

		ecjRunErr := ecjCmd.Run()
		ecjOutputStr := strings.ReplaceAll(ecjOut.String()+"\n"+ecjErr.String(), "\r\n", "\n")
		ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
		ecjOutputStr = ansiRegex.ReplaceAllString(ecjOutputStr, "")

		if ecjRunErr == nil {
			fmt.Println("[INFO] Compilation succeeded with zero errors.")
			os.Exit(0)
		}

		// Intercept and parse ECJ error
		processECJFailure(projectDir, dbPath, javaPathClean, ecjOutputStr)
		os.Exit(1)

	} else {
		// Java mode (original logic)
		javaPathClean := targetPath
		
		cmdArgs := append([]string{"-17", "-proceedOnError", "-d", filepath.Join(projectDir, "bin")}, os.Args[2:]...)
		cmd := exec.Command(ecjPath, cmdArgs...)
		
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf

		err = cmd.Run()
		if err == nil {
			fmt.Println("[INFO] Compilation succeeded with zero errors.")
			os.Exit(0)
		}

		outputStr := strings.ReplaceAll(outBuf.String()+"\n"+errBuf.String(), "\r\n", "\n")
		ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
		outputStr = ansiRegex.ReplaceAllString(outputStr, "")
		processECJFailure(projectDir, dbPath, javaPathClean, outputStr)
		os.Exit(1)
	}
}

func processECJFailure(projectDir, dbPath, javaPathClean, outputStr string) {
	errorRegex := regexp.MustCompile(`(?s)\d+\.\s+ERROR\s+in\s+([^\s\(]+)\s+\(at\s+line\s+(\d+)\)\s*\n(.*?)\n\s*([-^]+)\s*\n\s*(.*?)\n-`)
	matches := errorRegex.FindAllStringSubmatch(outputStr, -1)
	if len(matches) == 0 {
		errorRegex = regexp.MustCompile(`(?m)^.*?ERROR in (.*?) \(at line (\d+)\)\n(.*?)\n`)
		matches = errorRegex.FindAllStringSubmatch(outputStr, -1)
	}

	if len(matches) == 0 {
		fmt.Printf("[ERROR] ECJ execution failed, but diagnostics could not be parsed. Output:\n%s\n", outputStr)
		os.Exit(1)
	}

	match := matches[0]
	filePath := match[1]
	lineNo, _ := strconv.Atoi(match[2])
	
	var errMsg string
	if len(match) >= 6 {
		errMsg = fmt.Sprintf("line %d: %s\nContext:\n%s\n%s", lineNo, match[5], match[3], match[4])
	} else {
		errMsg = fmt.Sprintf("line %d: %s", lineNo, match[3])
	}

	fmt.Printf("[INFO] Intercepted Compiler Error at %s (line %d)\n", filePath, lineNo)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("[ERROR] Failed to open context ledger: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	var symbolURI string
	query := `
		SELECT symbol_uri FROM declarations 
		WHERE file_path = ? AND ? >= start_line AND ? <= end_line
		ORDER BY (end_line - start_line) ASC 
		LIMIT 1`
		
	err = db.QueryRow(query, filePath, lineNo, lineNo).Scan(&symbolURI)
	if err == sql.ErrNoRows {
		relPath := strings.TrimPrefix(filePath, projectDir+"/")
		query = `
			SELECT symbol_uri FROM declarations 
			WHERE (file_path LIKE ? OR file_path = ?) AND ? >= start_line AND ? <= end_line
			ORDER BY (end_line - start_line) ASC 
			LIMIT 1`
		err = db.QueryRow(query, "%"+relPath, relPath, lineNo, lineNo).Scan(&symbolURI)
	}

	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("[ERROR] Failed to query declarations ledger: %v\n", err)
		os.Exit(1)
	}

	var sourceContext string
	if symbolURI != "" {
		fmt.Printf("[INFO] SQLite resolved logic symbol: %s\n", symbolURI)
		if resolved, err := ResolveContext(dbPath, symbolURI); err == nil {
			sourceContext = resolved
		}
	}

	if sourceContext == "" {
		fmt.Println("[WARNING] Could not resolve specific method block. Falling back to file-level context.")
		nrxPath := strings.Replace(filePath, ".java", ".nrx", 1)
		content, err := os.ReadFile(nrxPath)
		if err == nil {
			sourceContext = string(content)
		} else {
			sourceContext = fmt.Sprintf("[Error: Could not locate source file %s]", nrxPath)
		}
	}

	symbolContext := symbolURI
	if symbolContext == "" {
		symbolContext = "unknown method/file context"
	}

	generateAndWritePrompt(projectDir, symbolContext, errMsg, sourceContext)
}

func generateAndWritePrompt(projectDir, symbolContext, errMsg, sourceContext string) {
	var promptBuilder strings.Builder
	promptBuilder.WriteString("[SYSTEM INVARIANT VIOLATION]\n")
	promptBuilder.WriteString("The previous NetRexx output failed compiler validation.\n\n")
	
	promptBuilder.WriteString(fmt.Sprintf("Target Symbol Context: %s\n", symbolContext))
	promptBuilder.WriteString(fmt.Sprintf("Compiler Diagnostic:\n%s\n\n", errMsg))
	promptBuilder.WriteString("Original NetRexx Code Context:\n")
	promptBuilder.WriteString("```rexx\n")
	promptBuilder.WriteString(sourceContext)
	promptBuilder.WriteString("\n```\n\n")
	promptBuilder.WriteString("Correct the scope, type allocation, or syntax inside the target block.\n")
	promptBuilder.WriteString("IMPORTANT INSTRUCTIONS FOR NETREXX DIALECT:\n")
	promptBuilder.WriteString("1. Variable declarations MUST follow the NetRexx syntax: 'varName = Type initialValue' (e.g. 'dbPath = String null', 'avg = Rexx 0', 'count = int 0'). Do NOT use Java-style declarations like 'Type varName = value' or 'String dbPath = null' as they will cause syntax errors.\n")
	promptBuilder.WriteString("2. NetRexx methods do NOT have a terminating 'end' keyword at the method level. Only inner blocks like 'do', 'loop', and 'select' should be closed with 'end'. Do NOT append a trailing 'end' at the end of the method body.\n")
	promptBuilder.WriteString("3. Checked exceptions (like Exception, SQLException) can ONLY be caught inside a 'do ... catch' block if the body of that 'do' block calls a method that is explicitly declared to throw/signal that exception. If no such method is called, catching checked exceptions is a compile-time error. For 'main', do not catch checked exceptions, or just catch 'RuntimeException' / 'Throwable', or avoid catch blocks entirely.\n")
	promptBuilder.WriteString("4. If the block contains database operations using a 'dbPath = String' parameter, you MUST ensure they are guarded with `if dbPath \\= null & dbPath \\= \"null\" then do` before connecting via JDBC, to prevent creating database files in the current working directory named 'null'.\n")
	promptBuilder.WriteString("5. Mismatched block and catch syntax: (a) Never write 'finally do' on the same line. If you need a try-catch block inside a finally clause, place 'finally' on its own line and start the 'do' block on the next line. (b) Every catch clause MUST follow NetRexx syntax: 'catch ex = ExceptionType' (e.g. 'catch ex = SQLException' or 'catch ex = Exception'). Never write 'catch Exception ex' or 'catch ex' without a type.\n")
	promptBuilder.WriteString("Output ONLY the revised, complete NetRexx source block. Do not include markdown code block formatting or explanations outside the block.\n")

	fmt.Println("\n=================== GENERATED SELF-CORRECTION PROMPT ===================")
	fmt.Println(promptBuilder.String())
	fmt.Println("========================================================================")

	// Save prompt to a temporary file for API dispatcher consumption
	promptFile := filepath.Join(projectDir, ".context", "self_correct_prompt.txt")
	os.WriteFile(promptFile, []byte(promptBuilder.String()), 0644)
	fmt.Printf("[INFO] Prompt written to: %s\n", promptFile)

	originalBlockFile := filepath.Join(projectDir, ".context", "errant_block_original.txt")
	os.WriteFile(originalBlockFile, []byte(sourceContext), 0644)
	fmt.Printf("[INFO] Original block saved to: %s\n", originalBlockFile)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
