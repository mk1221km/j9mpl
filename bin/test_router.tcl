#!/usr/bin/env tclsh
# Transaction Router — runtime verification tests.
# Protocol: space-separated fields.
#   route <txId> <type> <amount>       →  "routed <txId> <channel>"
#   rate                               →  "0.0"
#   count <status>                     →  integer count

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

# ============================================================
# Transaction Router runtime tests
# ============================================================

# Routing rules:
#   type=high, amount <= 1000  → HIGH_PRIORITY_WIRE
#   type=low/other, amount <= 1000 → CLEAR_ACH
#   any type, amount > 1000    → HIGH_PRIORITY_WIRE

# 1. Small low-priority → CLEAR_ACH
if {[runTest $binary "route tx001 alice 500" "routed tx001 CLEAR_ACH"]} { incr passed } { incr failed }

# 2. Small high-priority → HIGH_PRIORITY_WIRE
if {[runTest $binary "route tx002 high 500" "routed tx002 HIGH_PRIORITY_WIRE"]} { incr passed } { incr failed }

# 3. Large low-priority → HIGH_PRIORITY_WIRE (amount overrides)
if {[runTest $binary "route tx003 alice 5000" "routed tx003 HIGH_PRIORITY_WIRE"]} { incr passed } { incr failed }

# 4. Large high-priority → HIGH_PRIORITY_WIRE
if {[runTest $binary "route tx004 high 5000" "routed tx004 HIGH_PRIORITY_WIRE"]} { incr passed } { incr failed }

# 5. Boundary: amount exactly 1000, low → CLEAR_ACH
if {[runTest $binary "route tx005 alice 1000" "routed tx005 CLEAR_ACH"]} { incr passed } { incr failed }

# 6. Boundary: amount exactly 1000, high → HIGH_PRIORITY_WIRE
if {[runTest $binary "route tx006 high 1000" "routed tx006 HIGH_PRIORITY_WIRE"]} { incr passed } { incr failed }

# 7. Unknown type defaults to low → CLEAR_ACH
if {[runTest $binary "route tx007 urgent 500" "routed tx007 CLEAR_ACH"]} { incr passed } { incr failed }

# 8. ACH type → low → CLEAR_ACH (chaos injector compatibility)
if {[runTest $binary "route tx008 ACH 500" "routed tx008 CLEAR_ACH"]} { incr passed } { incr failed }

# 9. Multiple routes + rate
if {[runTest $binary "route tx001 alice 500\nroute tx002 high 2000\nrate" "routed tx001 CLEAR_ACH\nrouted tx002 HIGH_PRIORITY_WIRE\n0.0"]} { incr passed } { incr failed }

# 10. Transaction count
if {[runTest $binary "route tx001 alice 500\nroute tx002 high 100\ncount routed" "routed tx001 CLEAR_ACH\nrouted tx002 HIGH_PRIORITY_WIRE\n2"]} { incr passed } { incr failed }

# 11. Count of non-existent status
if {[runTest $binary "count rejected" "0"]} { incr passed } { incr failed }

# 12. Empty input
if {[runTest $binary "" ""]} { incr passed } { incr failed }

# 13. Unknown command (graceful ignore)
if {[runTest $binary "route tx001 alice 500\nblargh\nrate" "routed tx001 CLEAR_ACH\n0.0"]} { incr passed } { incr failed }

# 14. Malformed route (missing fields) → INVALID
if {[runTest $binary "route tx001" "INVALID"]} { incr passed } { incr failed }

# 15. Non-numeric amount → INVALID
if {[runTest $binary "route tx001 alice abc" "INVALID"]} { incr passed } { incr failed }

puts stderr "\[ROUTER\] $passed passed, $failed failed"
if {$failed > 0} { exit 1 }
exit 0
