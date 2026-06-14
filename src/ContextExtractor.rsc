module ContextExtractor

import lang::java::m3::Core;
import IO;
import String;
import Relation;
import Set;
import List;

// Function to escape single quotes for SQLite
str escapeSQL(str val) {
    return replaceAll(val, "\'", "\'\'");
}

// Extract structural context from the target directory and generate SQL commands
void extractContext(loc projectDir, loc outputSqlFile) {
    println("Generating M3 model for directory: <projectDir>...");
    M3 model = createM3FromDirectory(projectDir);
    
    list[str] sqlLines = [];
    
    // Wrap inserts in a transaction for microsecond-speed bulk processing
    sqlLines += "BEGIN TRANSACTION;";
    sqlLines += "DELETE FROM declarations;";
    sqlLines += "DELETE FROM containment;";
    sqlLines += "DELETE FROM symbol_uses;";
    
    println("Extracting declarations...");
    // 1. Extract declarations mapping: rel[loc name, loc src]
    for (<name, src> <- model@declarations) {
        str symbolUri = "<name>";
        str filePath = src.path? ? src.path : "";
        int startLine = src.begin? ? src.begin.line : 0;
        int endLine = src.end? ? src.end.line : 0;
        
        sqlLines += "INSERT OR REPLACE INTO declarations (symbol_uri, file_path, start_line, end_line) VALUES (\'<escapeSQL(symbolUri)>\', \'<escapeSQL(filePath)>\', <startLine>, <endLine>);";
    }
    
    println("Extracting containment relations...");
    // 2. Extract containment hierarchy: rel[loc parent, loc child]
    for (<parent, child> <- model@containment) {
        str parentUri = "<parent>";
        str childUri = "<child>";
        
        sqlLines += "INSERT OR REPLACE INTO containment (parent_uri, child_uri) VALUES (\'<escapeSQL(parentUri)>\', \'<escapeSQL(childUri)>\');";
    }
    
    println("Extracting symbol references and mapping callers...");
    // 3. Extract symbol usages and map them to their enclosing logical callers
    // First, map physical source ranges back to their logical declarations
    map[loc src, loc name] physicalToLogical = (src: name | <name, src> <- model@declarations);
    
    // Extract usages: rel[loc use, loc name]
    for (<use, name> <- model@uses) {
        str calleeUri = "<name>";
        str locationSpan = "unknown";
        if (use.begin? && use.end?) {
            locationSpan = "<use.begin.line>:<use.begin.column>-<use.end.line>:<use.end.column>";
        }
        
        // Locate caller by checking which declaration's physical span contains the physical reference
        str callerUri = "unknown";
        if (use.path? && use.offset? && use.length?) {
            for (declSrc <- physicalToLogical) {
                if (declSrc.path? && declSrc.offset? && declSrc.length? &&
                    use.path == declSrc.path && 
                    use.offset >= declSrc.offset && 
                    (use.offset + use.length) <= (declSrc.offset + declSrc.length)) {
                    
                    loc callerLoc = physicalToLogical[declSrc];
                    // Prioritize method or constructor caller context over general class scope
                    if (callerLoc.scheme == "java+method" || callerLoc.scheme == "java+constructor") {
                        callerUri = "<callerLoc>";
                        break;
                    } else if (callerUri == "unknown") {
                        callerUri = "<callerLoc>";
                    }
                }
            }
        }
        
        sqlLines += "INSERT INTO symbol_uses (caller_uri, callee_uri, location_span) VALUES (\'<escapeSQL(callerUri)>\', \'<escapeSQL(calleeUri)>\', \'<escapeSQL(locationSpan)>\');";
    }
    
    sqlLines += "COMMIT;";
    
    // Save to target file
    writeFile(outputSqlFile, intercalate("\n", sqlLines));
    println("Successfully generated SQL script at: <outputSqlFile>");
}

// Main execution entry point
void main(list[str] args) {
    if (size(args) < 2) {
        println("Usage: java -jar rascal-shell-stable.jar ContextExtractor.rsc <projectDir> <outputSqlFile>");
        return;
    }
    
    str pDir = args[0];
    str oFile = args[1];
    
    // Strip leading slash if present to prevent double slash in URI compilation
    if (startsWith(pDir, "/")) {
        pDir = substring(pDir, 1);
    }
    if (startsWith(oFile, "/")) {
        oFile = substring(oFile, 1);
    }
    
    loc projectDirLoc = |file:///| + pDir;
    loc outputSqlFileLoc = |file:///| + oFile;
    
    extractContext(projectDirLoc, outputSqlFileLoc);
}
