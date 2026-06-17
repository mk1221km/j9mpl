#!/usr/bin/env tclsh
# Behavioral test harness — language-agnostic, drives any compiled binary.
# Tests verify behavior at runtime only (shift-right). No compile-time checks.
# Usage: tclsh bin/test_harness.tcl <binary_path>

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
        puts stderr "  input:    [string map {\n {\\n}} $input]"
        puts stderr "  expected: [string map {\n {\\n}} $expected]"
        puts stderr "  got:      [string map {\n {\\n}} $output]"
        return 0
    }
    return 1
}

set passed 0
set failed 0

# ============================================================
# Ring Buffer runtime verification tests
# ============================================================
if {[string match "*ringbuffer*" $binary] || [string match "*ring*" $binary]} {

    # --- avg tests ---

    # 1. Empty buffer avg
    if {[runTest $binary "avg" "0.0"]} { incr passed } { incr failed }

    # 2. Single push then avg
    if {[runTest $binary "push 5.0\navg" "5.0"]} { incr passed } { incr failed }

    # 3. Multiple pushes then avg
    if {[runTest $binary "push 1.0\npush 2.0\npush 3.0\navg" "2.0"]} { incr passed } { incr failed }

    # 4. Avg with zero values
    if {[runTest $binary "push 0.0\npush 0.0\npush 0.0\navg" "0.0"]} { incr passed } { incr failed }

    # 5. Avg with negative values
    if {[runTest $binary "push -5.0\npush -3.0\npush -1.0\navg" "-3.0"]} { incr passed } { incr failed }

    # 6. Avg with mixed positive and negative
    if {[runTest $binary "push -10.0\npush 10.0\navg" "0.0"]} { incr passed } { incr failed }

    # 7. Avg with large values
    if {[runTest $binary "push 1e10\npush 2e10\navg" "15000000000.0"]} { incr passed } { incr failed }

    # 8. Avg after exactly 1 push (boundary)
    if {[runTest $binary "push 42.0\navg" "42.0"]} { incr passed } { incr failed }

    # --- readRange tests ---

    # 9. readRange returns most recent items (newest first is last in output)
    if {[runTest $binary "push 10.0\npush 20.0\npush 30.0\nreadRange 2" "20.0\n30.0"]} { incr passed } { incr failed }

    # 10. readRange from empty buffer
    if {[runTest $binary "readRange 5" ""]} { incr passed } { incr failed }

    # 11. readRange requesting more than available
    if {[runTest $binary "push 1.0\npush 2.0\nreadRange 10" "1.0\n2.0"]} { incr passed } { incr failed }

    # 12. readRange requesting 0 items
    if {[runTest $binary "push 1.0\npush 2.0\nreadRange 0" ""]} { incr passed } { incr failed }

    # 13. readRange default (no count argument)
    if {[runTest $binary "push 7.0\npush 8.0\nreadRange" "8.0"]} { incr passed } { incr failed }

    # 14. readRange after single push
    if {[runTest $binary "push 99.0\nreadRange 1" "99.0"]} { incr passed } { incr failed }

    # --- Wraparound tests ---

    # 15. Fill buffer completely, verify avg of full buffer
    set wr_input ""
    for {set i 0} {$i < 1024} {incr i} {
        append wr_input "push [expr {$i + 1}].0\n"
    }
    append wr_input "avg"
    # Sum of 1..1024 = 1024*1025/2 = 524800, avg = 524800/1024 = 512.5
    if {[runTest $binary $wr_input "512.5"]} { incr passed } { incr failed }

    # 16. Overflow: push 1025 values, oldest (1.0) should be gone
    set ov_input ""
    for {set i 0} {$i < 1025} {incr i} {
        append ov_input "push [expr {$i + 1}].0\n"
    }
    append ov_input "avg"
    # Sum of 2..1025 = (2+1025)*1024/2 = 1027*512 = 525824, avg = 525824/1024 = 513.5
    if {[runTest $binary $ov_input "513.5"]} { incr passed } { incr failed }

    # 17. Double overflow: push 2048 values (2 full cycles), then readRange
    set do_input ""
    for {set i 0} {$i < 2048} {incr i} {
        append do_input "push [expr {$i + 1}].0\n"
    }
    append do_input "readRange 3"
    # After 2048 pushes, most recent 3 are 2046, 2047, 2048
    if {[runTest $binary $do_input "2046.0\n2047.0\n2048.0"]} { incr passed } { incr failed }

    # 18. Overflow: readRange after exactly 1024 pushes (no overwrite yet)
    set full_input ""
    for {set i 0} {$i < 1024} {incr i} {
        append full_input "push [expr {$i + 1}].0\n"
    }
    append full_input "readRange 3"
    if {[runTest $binary $full_input "1022.0\n1023.0\n1024.0"]} { incr passed } { incr failed }

    # --- utilization tests ---

    # 19. Empty utilization
    if {[runTest $binary "utilization" "0.0000"]} { incr passed } { incr failed }

    # 20. Partial utilization
    if {[runTest $binary "push 1.0\npush 2.0\npush 3.0\nutilization" "0.0029"]} { incr passed } { incr failed }

    # 21. Full utilization
    set fullutil_input ""
    for {set i 0} {$i < 1024} {incr i} {
        append fullutil_input "push 1.0\n"
    }
    append fullutil_input "utilization"
    if {[runTest $binary $fullutil_input "1.0000"]} { incr passed } { incr failed }

    # 22. Utilization after overflow (still full)
    set overflowutil_input ""
    for {set i 0} {$i < 1100} {incr i} {
        append overflowutil_input "push 1.0\n"
    }
    append overflowutil_input "utilization"
    if {[runTest $binary $overflowutil_input "1.0000"]} { incr passed } { incr failed }

    # --- Edge cases ---

    # 23. Push single value then readRange after overflow
    set single_ov_input ""
    for {set i 0} {$i < 1024} {incr i} {
        append single_ov_input "push 0.0\n"
    }
    append single_ov_input "push 999.0\nreadRange 1"
    if {[runTest $binary $single_ov_input "999.0"]} { incr passed } { incr failed }

    # 24. Alternating operations
    if {[runTest $binary "push 10.0\navg\npush 20.0\navg\nreadRange 2" "10.0\n15.0\n10.0\n20.0"]} { incr passed } { incr failed }

    # 25. Utilization after single push
    if {[runTest $binary "push 5.0\nutilization" "0.0010"]} { incr passed } { incr failed }

    # --- Summary ---
    puts stderr "\[RING\] $passed passed, $failed failed"

# ============================================================
# TransactionRouter behavioral tests (TODO)
# ============================================================
} elseif {[string match "*router*" $binary] || [string match "*txn*" $binary]} {
    puts stderr "\[ROUTER\] No tests defined yet"
    if {[runTest $binary "init\nroute tx001 alice bob 500 low\nstats" ""]} { incr passed } { incr failed }
    puts stderr "\[ROUTER\] $passed passed, $failed failed"
}

if {$failed > 0} {
    exit 1
}
exit 0
