#!/usr/bin/env tclsh
# Parse go test -coverprofile output, ingest zero-execution blocks into SQLite.
# Usage: tclsh bin/parse_coverage.tcl <project_dir> <coverage_file>

if {[llength $argv] < 2} {
    puts stderr "Usage: [file tail [info script]] <project_dir> <coverage_file>"
    exit 1
}

set projectDir [file normalize [lindex $argv 0]]
set coverageFile [file normalize [lindex $argv 1]]
set dbPath [file join $projectDir ".context" "project_context.db"]

if {![file exists $coverageFile]} {
    puts stderr "[ERROR] Coverage file not found: $coverageFile"
    exit 1
}

package require sqlite3
sqlite3 db $dbPath

set fp [open $coverageFile r]
set lines [split [read $fp] "\n"]
close $fp

set ingestCount 0
set sourceDir [file join $projectDir "generated"]

foreach line $lines {
    set line [string trim $line]
    if {$line eq "" || [string match "mode:*" $line]} {
        continue
    }
    # Format: filename.go:startLine.startCol,endLine.endCol numStatements execCount
    if {[regexp {^([^:]+):([0-9]+)\.[0-9]+,([0-9]+)\.[0-9]+\s+([0-9]+)\s+([0-9]+)$} $line -> filePath startLine endLine stmtCount execCount]} {
        if {$execCount == 0} {
            # Extract source lines for the uncovered block
            set sourceFile [file join $sourceDir [file tail $filePath]]
            set sourceBlock ""
            if {[file exists $sourceFile]} {
                set sfp [open $sourceFile r]
                set srcLines [split [read $sfp] "\n"]
                close $sfp
                set sourceBlock [join [lrange $srcLines [expr {$startLine - 1}] [expr {$endLine - 1}]] "\n"]
            }

            # Determine method name from source context
            set methodName "unknown"
            if {[file exists $sourceFile]} {
                set sfp [open $sourceFile r]
                set srcLines [split [read $sfp] "\n"]
                close $sfp
                for {set i [expr {$startLine - 1}]} {$i >= 0} {incr i -1} {
                    set srcLine [string trim [lindex $srcLines $i]]
                    if {[regexp {^func\s+(\([^)]+\)\s+)?([A-Za-z0-9_]+)} $srcLine -> _ mname]} {
                        set methodName $mname
                        break
                    }
                }
            }

            set now [clock seconds]
            db eval {
                INSERT OR REPLACE INTO generated_coverage_gaps
                (file_path, method_name, start_line, end_line, statement_count, unexecuted_source_text, identified_at)
                VALUES ($filePath, $methodName, $startLine, $endLine, $stmtCount, $sourceBlock, $now)
            }
            incr ingestCount
        }
    }
}

db close
puts "\[INFO\] Coverage analysis complete: $ingestCount gaps ingested."
exit 0
