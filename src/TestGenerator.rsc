module TestGenerator

import IO;
import String;
import List;

// Parse parameter list into individual types
list[str] parseParams(str paramsStr) {
    if (trim(paramsStr) == "") {
        return [];
    }
    return split(",", paramsStr);
}

// Generate boundary values based on type
str getBoundaryList(str paramType, str methodName, int paramIndex) {
    if (contains(paramType, "String")) {
        if (paramIndex == 0 && (methodName == "initDatabase" || methodName == "logMetric" || methodName == "getAverageMetric")) {
            return "dbPathBounds";
        }
        return "stringBounds";
    } else if (contains(paramType, "Rexx")) {
        return "rexxBounds";
    } else if (contains(paramType, "int") || contains(paramType, "double") || contains(paramType, "float") || contains(paramType, "long")) {
        return "doubleBounds";
    } else if (contains(paramType, "MetricRecord")) {
        return "recordBounds";
    }
    return "[\"null\"]";
}

void generateTest(str className, loc declsFile, loc testFile) {
    println("Generating NetRexx property-based test harness for class: <className>...");
    
    list[str] declLines = [];
    try {
        declLines = readFileLines(declsFile);
    } catch ex: {
        println("[ERROR] Failed to read declarations CSV file: <declsFile>");
        return;
    }
    
    // Structure to hold method details
    // We only care about methods belonging to the target class
    list[tuple[str name, list[str] params]] methods = [];
    str packageName = "";
    
    for (line <- declLines) {
        if (/java\+method:\/\/\/<classPath:[a-zA-Z0-9_\/]+>\/<clsName:\w+>\/<methodName:\w+>\(<params:[^)]*>\)/ := line) {
            if (clsName == className) {
                packageName = replaceAll(classPath, "/", ".");
                // Skip constructor or main method for boundary exhaustion
                if (methodName != className && methodName != "main") {
                    methods += <methodName, parseParams(params)>;
                }
            }
        }
    }
    
    if (packageName == "") {
        packageName = "com.factory"; // fallback
    }
    
    list[str] code = [];
    code += "package <packageName>";
    code += "options binary";
    code += "import java.sql.SQLException";
    code += "";
    code += "class <className>Test public";
    code += "  method main(args = String[]) public static";
    code += "    say \"=== [Phase III] Starting Boundary Input Exhaustion Test for <className> ===\"";
    code += "    ";
    code += "    -- Boundary payloads";
    code += "    stringBounds = [\"\", \"normal_string_test\", \"\'; DROP TABLE system_metrics; --\", \"null\"]";
    code += "    dbPathBounds = [\"generated/metrics_test.db\", \":memory:\", \"null\"]";
    code += "    doubleBounds = [Rexx(0), Rexx(1), Rexx(-1), Rexx(999999999), Rexx(1.79e+308), Rexx(-1.79e+308)]";
    code += "    rexxBounds = [Rexx(0), Rexx(1), Rexx(-1), Rexx(\"normal\"), Rexx(\"\")]";
    code += "    ";
    code += "    -- Build MetricRecord boundary instances";
    code += "    recordBounds = java.util.ArrayList()";
    code += "    recordBounds.add(null)";
    code += "    loop tsVal over stringBounds";
    code += "      loop nameVal over stringBounds";
    code += "        loop valVal over doubleBounds";
    code += "          rec = MetricRecord()";
    code += "          if tsVal \\= \"null\" then rec.timestamp = String tsVal";
    code += "          if nameVal \\= \"null\" then rec.metricName = String nameVal";
    code += "          rec.metricValue = valVal";
    code += "          recordBounds.add(rec)";
    code += "        end";
    code += "      end";
    code += "    end";
    code += "";
    
    // For each method, generate a nested testing loop
    for (m <- methods) {
        code += "    say \"Testing method <m.name>...\"";
        
        // Generate nested loops for each parameter
        int indent = 4;
        list[str] loopVars = [];
        
        for (i <- [0..size(m.params)]) {
            str paramType = m.params[i];
            str varName = "<m.name>_p<i+1>";
            str boundaryList = getBoundaryList(paramType, m.name, i);
            
            str indentStr = left("", indent, " ");
            code += "<indentStr>loop val<i+1> over <boundaryList>";
            code += "<indentStr>  <varName> = <paramType> null";
            
            // Casting/null check logic
            if (contains(paramType, "String")) {
                code += "<indentStr>  if val<i+1> \\= \"null\" then <varName> = String val<i+1>";
            } else if (contains(paramType, "MetricRecord")) {
                code += "<indentStr>  <varName> = MetricRecord val<i+1>";
            } else {
                code += "<indentStr>  <varName> = val<i+1>";
            }
            
            loopVars += varName;
            indent += 2;
        }
        
        // Inside the innermost loop: make the call wrapped in do-catch-finally
        str innerIndent = left("", indent, " ");
        code += "<innerIndent>do";
        
        str callArgs = intercalate(", ", loopVars);
        code += "<innerIndent>  <className>.<m.name>(<callArgs>)";
        code += "<innerIndent>catch ex = RuntimeException";
        code += "<innerIndent>  -- Silent capture of expected exceptions";
        code += "<innerIndent>  nop";
        code += "<innerIndent>end";
        
        // Close loops
        for (i <- [0..size(m.params)]) {
            indent -= 2;
            str closeIndent = left("", indent, " ");
            code += "<closeIndent>end";
        }
        code += "    say \"  Method <m.name> boundary exhaustion completed.\"";
        code += "";
    }
    
    code += "    say \"=== [Phase III] Boundary Input Exhaustion Test Completed successfully! ===\"";
    
    writeFile(testFile, intercalate("\n", code));
    println("Successfully generated test script at: <testFile>");
}

void main(list[str] args) {
    if (size(args) < 3) {
        println("Usage: rascal TestGenerator <className> <declarationsCsvPath> <outputTestNrxPath>");
        return;
    }
    
    str className = args[0];
    str declsCsv = args[1];
    str outputNrx = args[2];
    
    // Remove leading slash if needed
    if (startsWith(declsCsv, "/")) declsCsv = substring(declsCsv, 1);
    if (startsWith(outputNrx, "/")) outputNrx = substring(outputNrx, 1);
    
    loc declsFileLoc = |file:///| + declsCsv;
    loc testFileLoc = |file:///| + outputNrx;
    
    generateTest(className, declsFileLoc, testFileLoc);
}
