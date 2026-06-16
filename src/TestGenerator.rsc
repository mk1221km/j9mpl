module TestGenerator

import IO;
import String;
import List;
import lang::json::IO;

// Parse parameter list into individual types
list[str] parseParams(str paramsStr) {
    if (trim(paramsStr) == "") {
        return [];
    }
    return split(",", paramsStr);
}

// Generate boundary values based on type
str getBoundaryList(str paramType, str methodName, int paramIndex, str recordType) {
    if (contains(paramType, "String")) {
        if (paramIndex == 0 && (methodName == "initDatabase" || methodName == "logMetric" || methodName == "getAverageMetric" || methodName == "initRoutingTable" || methodName == "routeTransaction" || methodName == "getTransactionCount")) {
            return "dbPathBounds";
        }
        return "stringBounds";
    } else if (contains(paramType, "Rexx")) {
        return "rexxBounds";
    } else if (contains(paramType, "int") || contains(paramType, "double") || contains(paramType, "float") || contains(paramType, "long")) {
        return "doubleBounds";
    } else if (recordType != "" && contains(paramType, recordType)) {
        return "recordBounds";
    }
    return "[\"null\"]";
}

data BoundaryPayload = boundaryPayload(list[str] valid, list[map[str, str]] counter);
data BoundaryItem = boundaryItem(str domain, BoundaryPayload payload);
alias BoundariesMap = map[str, list[BoundaryItem]];
alias FuzzTuple = tuple[str val, int isCounter, str expected];

str escapeRexxStr(str s) {
    return replaceAll(replaceAll(s, "\\", "\\\\"), "\'", "\'\'");
}

list[FuzzTuple] getStringBounds(BoundariesMap bm) {
    list[FuzzTuple] vals = [<"normal_string_test", 0, "null">];
    if ("String" in bm) {
        for (item <- bm["String"]) {
            for (v <- item.payload.valid) {
                vals += <v, 0, "null">;
            }
            for (c <- item.payload.counter) {
                vals += <c["value"], 1, "\'" + c["expected"] + "\'">;
            }
        }
    }
    return dup(vals);
}

list[FuzzTuple] getDoubleBounds(BoundariesMap bm) {
    list[FuzzTuple] vals = [<"0", 0, "null">, <"1", 0, "null">, <"-1", 0, "null">];
    if ("double" in bm) {
        for (item <- bm["double"]) {
            for (v <- item.payload.valid) {
                vals += <v, 0, "null">;
            }
            for (c <- item.payload.counter) {
                str s = c["value"];
                if (s != "NaN" && s != "Infinity" && s != "-Infinity") {
                    vals += <s, 1, "\'" + c["expected"] + "\'">;
                }
            }
        }
    }
    return dup(vals);
}

list[FuzzTuple] getIntBounds(BoundariesMap bm) {
    list[FuzzTuple] vals = [<"0", 0, "null">, <"1", 0, "null">, <"-1", 0, "null">];
    if ("int" in bm) {
        for (item <- bm["int"]) {
            for (v <- item.payload.valid) {
                vals += <v, 0, "null">;
            }
            for (c <- item.payload.counter) {
                vals += <c["value"], 1, "\'" + c["expected"] + "\'">;
            }
        }
    }
    return dup(vals);
}

str formatFuzzInput(FuzzTuple t) {
    str vRepr = "null";
    if (t.val != "null" && t.val != "NULL") {
        vRepr = "\'" + escapeRexxStr(t.val) + "\'";
    }
    return "FuzzInput(<vRepr>, <t.isCounter>, <t.expected>)";
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
                // Skip constructor, main, and validate helper methods for boundary exhaustion
                if (methodName != className && methodName != "main" && !startsWith(methodName, "validate")) {
                    methods += <methodName, parseParams(params)>;
                }
            }
        }
    }
    
    if (packageName == "") {
        packageName = "com.factory"; // fallback
    }
    
    // Find custom record types in method parameters
    str recordType = "";
    for (m <- methods) {
        for (p <- m.params) {
            str tp = trim(p);
            if (tp != "" && !contains(tp, "String") && !contains(tp, "Rexx") && 
                !contains(tp, "int") && !contains(tp, "double") && !contains(tp, "float") && !contains(tp, "long") && !contains(tp, "[]")) {
                if (contains(tp, ".")) {
                    list[str] parts = split(".", tp);
                    recordType = parts[size(parts)-1];
                } else {
                    recordType = tp;
                }
            }
        }
    }
    if (recordType == "") {
        recordType = "MetricRecord"; // fallback
    }
    
    // Find fields of the recordType dynamically
    list[tuple[str name, str typeName]] recordFields = [];
    for (line <- declLines) {
        if (/java\+field:\/\/\/[a-zA-Z0-9_\/]+\/<recName:\w+>\/<fieldName:\w+>/ := line) {
            if (recName == recordType) {
                str fType = "String";
                if (fieldName == "amount" || fieldName == "metricValue" || fieldName == "voltage") {
                    fType = "Rexx";
                }
                if (!startsWith(fieldName, "$")) {
                    recordFields += <fieldName, fType>;
                }
            }
        }
    }
    if (size(recordFields) == 0) {
        recordFields = [<"timestamp", "String">, <"metricName", "String">, <"metricValue", "Rexx">];
    }
    
    // Generate and write type manifest JSON
    list[str] methodJsonList = [];
    for (m <- methods) {
        list[str] paramJsonList = [];
        for (i <- [0..size(m.params)]) {
            paramJsonList += "{\"index\": <i>, \"type\": \"<trim(m.params[i])>\"}";
        }
        str pJson = intercalate(", ", paramJsonList);
        methodJsonList += "{\"name\": \"<m.name>\", \"parameters\": [<pJson>]}";
    }
    str mJson = intercalate(", ", methodJsonList);

    list[str] fieldJsonList = [];
    for (f <- recordFields) {
        fieldJsonList += "{\"name\": \"<f.name>\", \"type\": \"<f.typeName>\"}";
    }
    str fJson = intercalate(", ", fieldJsonList);

    str jsonContent = "{\n  \"class\": \"<className>\",\n  \"methods\": [<mJson>],\n  \"records\": [\n    {\n      \"name\": \"<recordType>\",\n      \"fields\": [<fJson>]\n    }\n  ]\n}";
    
    loc manifestLoc = declsFile;
    manifestLoc.path = replaceFirst(declsFile.path, "declarations.csv", "type_manifest.json");
    writeFile(manifestLoc, jsonContent);

    list[str] code = [];
    code += "package <packageName>";
    code += "options binary";
    code += "import java.sql.SQLException";
    code += "import <packageName>.<recordType>";
    code += "import <packageName>.<className>";
    code += "";
    code += "class <className>Test public";
    code += "  method main(args = String[]) public static";
    code += "    say \"=== [Phase III] Starting Boundary Input Exhaustion Test for <className> ===\"";
    loc boundariesLoc = declsFile;
    boundariesLoc.path = replaceFirst(declsFile.path, "declarations.csv", "fuzzer_boundaries.json");
    BoundariesMap bm = ();
    try {
        bm = readJSON(#BoundariesMap, boundariesLoc);
    } catch ex: {
        println("[WARNING] Failed to read fuzzer boundaries JSON file: <boundariesLoc>. Using defaults.");
    }

    list[FuzzTuple] strVals = getStringBounds(bm);
    list[str] strReprs = [];
    for (t <- strVals) {
        strReprs += formatFuzzInput(t);
    }
    
    list[FuzzTuple] doubleVals = getDoubleBounds(bm);
    list[str] doubleReprs = [];
    for (t <- doubleVals) {
        doubleReprs += formatFuzzInput(t);
    }

    list[FuzzTuple] intVals = getIntBounds(bm);
    list[str] rexxReprs = [
        "FuzzInput(\'normal\', 0, null)",
        "FuzzInput(\'\', 0, null)"
    ];
    for (t <- intVals) {
        rexxReprs += formatFuzzInput(t);
    }

    list[str] dbPathReprs = [
        "FuzzInput(\'generated/<toLowerCase(className)>_test.db\', 0, null)",
        "FuzzInput(\':memory:\', 0, null)",
        "FuzzInput(null, 1, \'java.lang.IllegalArgumentException\')"
    ];

    code += "    ";
    code += "    -- Boundary payloads";
    code += "    stringBounds = [<intercalate(", ", strReprs)>]";
    code += "    dbPathBounds = [<intercalate(", ", dbPathReprs)>]";
    code += "    doubleBounds = [<intercalate(", ", doubleReprs)>]";
    code += "    rexxBounds = [<intercalate(", ", rexxReprs)>]";
    code += "    ";
    code += "    -- Build <recordType> boundary instances";
    code += "    recordBounds = java.util.ArrayList()";
    code += "    recordBounds.add(RecordFuzzInput(null, 1, \'java.lang.IllegalArgumentException\'))";
    code += "    loop tsVal over stringBounds";
    code += "      loop nameVal over stringBounds";
    code += "        loop valVal over doubleBounds";
    code += "          rec = <recordType>()";
    
    int strVarIdx = 1;
    for (f <- recordFields) {
        if (f.typeName == "String") {
            str loopVar = (strVarIdx % 2 == 1) ? "tsVal" : "nameVal";
            code += "          if <loopVar>.val \\== null then rec.<f.name> = String <loopVar>.val";
            strVarIdx += 1;
        } else if (f.typeName == "Rexx") {
            code += "          if valVal.val \\== null then rec.<f.name> = Rexx(valVal.val)";
        }
    }
    
    code += "          recIsCounter = tsVal.isCounter | nameVal.isCounter | valVal.isCounter";
    code += "          recExpected = String null";
    code += "          if tsVal.isCounter \\= 0 then do";
    code += "            recExpected = tsVal.expected";
    code += "          end";
    code += "          else if nameVal.isCounter \\= 0 then do";
    code += "            recExpected = nameVal.expected";
    code += "          end";
    code += "          else if valVal.isCounter \\= 0 then do";
    code += "            recExpected = valVal.expected";
    code += "          end";
    code += "          counterCount = tsVal.isCounter + nameVal.isCounter + valVal.isCounter";
    code += "          if counterCount == 0 | counterCount == 1 then do";
    code += "            recordBounds.add(RecordFuzzInput(rec, recIsCounter, recExpected))";
    code += "          end";
    code += "        end";
    code += "      end";
    code += "    end";
    code += "    if stringBounds \\== null & dbPathBounds \\== null & doubleBounds \\== null & rexxBounds \\== null & recordBounds \\== null then say \"ok\"";
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
            str boundaryList = getBoundaryList(paramType, m.name, i, recordType);
            
            str indentStr = left("", indent, " ");
            code += "<indentStr>loop <m.name>_val<i+1> over <boundaryList>";
            if (boundaryList == "recordBounds") {
                code += "<indentStr>  <m.name>_val<i+1>_input = RecordFuzzInput <m.name>_val<i+1>";
            } else {
                code += "<indentStr>  <m.name>_val<i+1>_input = FuzzInput <m.name>_val<i+1>";
            }
            code += "<indentStr>  <varName> = <paramType> null";
            
            loopVars += varName;
            indent += 2;
        }
        
        // Inside the innermost loop: make the call wrapped in do-catch-finally
        str innerIndent = left("", indent, " ");
        
        code += "<innerIndent><m.name>_counterCount = int 0";
        for (i <- [0..size(m.params)]) {
            code += "<innerIndent>if <m.name>_val<i+1>_input.isCounter \\= 0 then <m.name>_counterCount = <m.name>_counterCount + 1";
        }
        
        code += "<innerIndent>if <m.name>_counterCount \<= 1 then do";
        
        code += "<innerIndent>  <m.name>_expectedEx = String null";
        for (i <- [0..size(m.params)]) {
            code += "<innerIndent>  if <m.name>_val<i+1>_input.isCounter \\= 0 then do";
            code += "<innerIndent>    <m.name>_expectedEx = <m.name>_val<i+1>_input.expected";
            code += "<innerIndent>  end";
        }
        
        code += "<innerIndent>  do";
        code += "<innerIndent>    if 1 == 0 then <className>Test.dummySignal()";
        
        // Perform parameter assignments inside the do block
        for (i <- [0..size(m.params)]) {
            str paramType = m.params[i];
            str varName = "<m.name>_p<i+1>";
            if (contains(paramType, "String")) {
                code += "<innerIndent>    if <m.name>_val<i+1>_input.val \\== null then <varName> = <m.name>_val<i+1>_input.val";
            } else if (recordType != "" && contains(paramType, recordType)) {
                code += "<innerIndent>    <varName> = <m.name>_val<i+1>_input.rec";
            } else if (contains(paramType, "Rexx")) {
                code += "<innerIndent>    if <m.name>_val<i+1>_input.val \\== null then <varName> = Rexx(<m.name>_val<i+1>_input.val)";
            } else if (contains(paramType, "int")) {
                code += "<innerIndent>    if <m.name>_val<i+1>_input.val \\== null then <varName> = Integer.parseInt(<m.name>_val<i+1>_input.val)";
            } else if (contains(paramType, "double")) {
                code += "<innerIndent>    if <m.name>_val<i+1>_input.val \\== null then <varName> = Double.parseDouble(<m.name>_val<i+1>_input.val)";
            } else if (contains(paramType, "float")) {
                code += "<innerIndent>    if <m.name>_val<i+1>_input.val \\== null then <varName> = Float.parseFloat(<m.name>_val<i+1>_input.val)";
            } else if (contains(paramType, "long")) {
                code += "<innerIndent>    if <m.name>_val<i+1>_input.val \\== null then <varName> = Long.parseLong(<m.name>_val<i+1>_input.val)";
            }
        }
        
        str callArgs = intercalate(", ", loopVars);
        code += "<innerIndent>    <className>.<m.name>(<callArgs>)";
        code += "<innerIndent>    if <m.name>_expectedEx \\== null then do";
        code += "<innerIndent>      say \"Assertion Failure in <m.name>: counter-example bypassed validation (no exception thrown). Expected: \" || <m.name>_expectedEx";
        code += "<innerIndent>      java.lang.System.exit(1)";
        code += "<innerIndent>    end";
        code += "<innerIndent>  catch caughtEx = Throwable";
        code += "<innerIndent>    thrownClass = caughtEx.getClass().getName()";
        code += "<innerIndent>    if <m.name>_expectedEx == null then do";
        code += "<innerIndent>      say \"Assertion Failure in <m.name>: happy path regression (unexpected exception: \" || thrownClass || \": \" || caughtEx.getMessage() || \")\"";
        code += "<innerIndent>      caughtEx.printStackTrace()";
        code += "<innerIndent>      java.lang.System.exit(1)";
        code += "<innerIndent>    end";
        code += "<innerIndent>    else do";
        code += "<innerIndent>      matched = int 0";
        code += "<innerIndent>      do";
        code += "<innerIndent>        expectedClass = java.lang.Class.forName(<m.name>_expectedEx)";
        code += "<innerIndent>        if expectedClass.isInstance(caughtEx) then matched = 1";
        code += "<innerIndent>      catch ClassNotFoundException";
        code += "<innerIndent>        if thrownClass \\== <m.name>_expectedEx then matched = 1";
        code += "<innerIndent>      end";
        code += "<innerIndent>      if matched == 0 then do";
        code += "<innerIndent>        say \"Assertion Failure in <m.name>: caught \" || thrownClass || \" (\" || caughtEx.getMessage() || \") but expected \" || <m.name>_expectedEx";
        code += "<innerIndent>        java.lang.System.exit(1)";
        code += "<innerIndent>      end";
        code += "<innerIndent>    end";
        code += "<innerIndent>  end";
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
    code += "";
    
    code += "  method dummySignal() private static signals java.lang.Throwable";
    code += "    nop";
    code += "";
    
    // Add the helper classes
    code += "class FuzzInput";
    code += "  properties public";
    code += "    val = String";
    code += "    isCounter = int";
    code += "    expected = String";
    code += "  method FuzzInput(aVal = String, aIsCounter = int, aExpected = String)";
    code += "    this.val = aVal";
    code += "    this.isCounter = aIsCounter";
    code += "    this.expected = aExpected";
    code += "";
    code += "class RecordFuzzInput";
    code += "  properties public";
    code += "    rec = <recordType>";
    code += "    isCounter = int";
    code += "    expected = String";
    code += "  method RecordFuzzInput(aRec = <recordType>, aIsCounter = int, aExpected = String)";
    code += "    this.rec = aRec";
    code += "    this.isCounter = aIsCounter";
    code += "    this.expected = aExpected";
    code += "";
    
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
