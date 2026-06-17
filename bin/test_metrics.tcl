#!/usr/bin/env tclsh
# Metrics Logger — runtime verification tests.
# Protocol: space-separated fields.
#   log <name> <value> [timestamp]  →  (no output)
#   avg <name>                      →  numeric value
#   count <name>                    →  integer count

set binary [lindex $argv 0]
if {$binary eq ""} {
    puts stderr "Usage: [file tail [info script]] <binary_path>"
    exit 1
}
if {![file executable $binary]} {
    puts stderr "\[TEST\] Binary not found or not executable: $binary"
    exit 1
}

proc runTest {binary input expected} {
    set fd [open "|$binary" r+]
    puts $fd $input
    close $fd w
    set output [string trim [read $fd]]
    close $fd
    if {$output ne $expected} {
        puts stderr "\[TEST FAIL\]"
        regsub -all {\n} $input {\\n} in_esc
        puts stderr "  input:    $in_esc"
        puts stderr "  expected: [string map {\n {\\n}} $expected]"
        puts stderr "  got:      [string map {\n {\\n}} $output]"
        return 0
    }
    return 1
}

set passed 0
set failed 0

# 1. Log single value, then avg
if {[runTest $binary "log cpu 95.5\navg cpu" "95.5"]} { incr passed } { incr failed }

# 2. Log multiple values, avg returns mean
if {[runTest $binary "log cpu 10.0\nlog cpu 20.0\nlog cpu 30.0\navg cpu" "20.0"]} { incr passed } { incr failed }

# 3. Avg of empty metric
if {[runTest $binary "avg memory" "0.0"]} { incr passed } { incr failed }

# 4. Count of logged values
if {[runTest $binary "log cpu 1.0\nlog cpu 2.0\nlog cpu 3.0\ncount cpu" "3"]} { incr passed } { incr failed }

# 5. Count of empty metric
if {[runTest $binary "count disk" "0"]} { incr passed } { incr failed }

# 6. Multiple independent metrics
if {[runTest $binary "log cpu 50.0\nlog memory 100.0\nlog cpu 150.0\navg cpu\navg memory" "100.0\n100.0"]} { incr passed } { incr failed }

# 7. Single value: count then avg
if {[runTest $binary "log temp 37.5\ncount temp\navg temp" "1\n37.5"]} { incr passed } { incr failed }

# 8. Zero value
if {[runTest $binary "log cpu 0.0\navg cpu" "0.0"]} { incr passed } { incr failed }

# 9. Negative values
if {[runTest $binary "log signal -5.0\nlog signal -3.0\navg signal" "-4.0"]} { incr passed } { incr failed }

# 10. Large values
if {[runTest $binary "log large 10000000000.0\nlog large 20000000000.0\navg large" "15000000000.0"]} { incr passed } { incr failed }

# 11. Empty input
if {[runTest $binary "" ""]} { incr passed } { incr failed }

# 12. Malformed log (missing fields) silently ignored, subsequent correct log works
if {[runTest $binary "log 95.5\nlog cpu 10.0\navg cpu" "10.0"]} { incr passed } { incr failed }

# 13. Non-numeric value silently ignored
if {[runTest $binary "log cpu abc\nlog cpu 10.0\navg cpu" "10.0"]} { incr passed } { incr failed }

# 14. NaN value silently ignored
if {[runTest $binary "log cpu NaN\nlog cpu 10.0\navg cpu" "10.0"]} { incr passed } { incr failed }

# 15. Inf value silently ignored
if {[runTest $binary "log cpu Inf\nlog cpu 10.0\navg cpu" "10.0"]} { incr passed } { incr failed }

puts stderr "\[METRICS\] $passed passed, $failed failed"
if {$failed > 0} { exit 1 }
exit 0
