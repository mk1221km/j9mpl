#!/usr/bin/env tclsh
# Stress test: extreme wraparound on ring buffer binary
# Verifies correctness under heavy push cycles and modulo-boundary reads.

set binary [lindex $argv 0]
if {$binary eq ""} {
    puts stderr "Usage: [file tail [info script]] <binary>"
    exit 1
}

proc runTest {binary input expected} {
    set fd [open "|$binary" r+]
    puts $fd $input
    close $fd w
    set output [string trim [read $fd]]
    close $fd
    if {$output ne $expected} {
        puts stderr "STRESS FAIL:"
        regsub -all {\n} [string range $input 0 100] {\\n} input_short
        puts stderr "  input:    ${input_short}..."
        puts stderr "  expected: $expected"
        puts stderr "  got:      $output"
        return 0
    }
    return 1
}

set passed 0
set failed 0

# --- Stress 1: 10,000 push/avg cycles (values cycle 0..99) ---
puts "Stress 1: 10,000 push/avg..."
set input ""
for {set i 0} {$i < 10000} {incr i} {
    append input "push [expr {$i % 100}].0\n"
}
append input "avg"
# Last 1024 values: i=8976..9999, values = 8976%100..9999%100 = 76..99,0..99×10
# sum(76..99) = (76+99)*24/2 = 2100
# sum(0..99) = 4950, ×10 = 49500
# total = 51600, avg = 51600/1024 = 50.390625
if {[runTest $binary $input "50.4"]} { incr passed } { incr failed }

# --- Stress 2: readRange at 1000-push checkpoints ---
puts "Stress 2: readRange checkpoints..."
set input ""
for {set i 0} {$i < 5000} {incr i} {
    append input "push [expr {$i + 1}].0\n"
    if {$i == 999 || $i == 1999 || $i == 2999 || $i == 3999 || $i == 4999} {
        append input "readRange 1\n"
    }
}
# Checkpoint 1000: value = 1000.0
# Checkpoint 2000: buf full, overwriting, most recent = 2000.0
# Checkpoint 3000: 3000.0
# Checkpoint 4000: 4000.0
# Checkpoint 5000: 5000.0
set expected "1000.0\n2000.0\n3000.0\n4000.0\n5000.0"
if {[runTest $binary $input $expected]} { incr passed } { incr failed }

# --- Stress 3: readRange 1024 (full dump) after 10240 pushes ---
puts "Stress 3: readRange full buffer dump..."
set input ""
for {set i 0} {$i < 10240} {incr i} {
    append input "push [expr {$i + 1}].0\n"
}
append input "readRange 1024"
# After 10240 pushes (10 full overwrite cycles), most recent 1024: 9217..10240
set expected ""
for {set i 9217} {$i <= 10240} {incr i} {
    append expected "${i}.0\n"
}
set expected [string trim $expected]
if {[runTest $binary $input $expected]} { incr passed } { incr failed }

# --- Stress 4: readRange straddling modulo boundary ---
puts "Stress 4: modulo-boundary readRange..."
set input ""
for {set i 0} {$i < 1024} {incr i} {
    append input "push [expr {$i}].0\n"
}
# Push 3 more to wrap: head=3, count=1024
# buf[3..1023] = 3..1023, buf[0..2] = 1024..1026
foreach val {1024 1025 1026} {
    append input "push ${val}.0\n"
}
append input "readRange 6"
# Most recent 6: 1021, 1022, 1023, 1024, 1025, 1026
# readRange returns oldest-first: 1021..1026
set expected "1021.0\n1022.0\n1023.0\n1024.0\n1025.0\n1026.0"
if {[runTest $binary $input $expected]} { incr passed } { incr failed }

# --- Stress 5: Negative values under heavy wraparound ---
puts "Stress 5: negative values wraparound..."
set input ""
for {set i 0} {$i < 2000} {incr i} {
    append input "push [expr {-$i}].0\n"
}
append input "avg"
# Last 1024 values: -976, -977, ..., -1999
# Sum = -(976+1999)*1024/2 = -2975*512 = -1523200
# avg = -1523200/1024 = -1487.5
if {[runTest $binary $input "-1487.5"]} { incr passed } { incr failed }

# --- Stress 6: Stress 1 but using LuaJIT via the internal functions directly ---
# This exercises the internal function logic independent of the I/O layer

puts stderr "STRESS: $passed passed, $failed failed"
if {$failed > 0} { exit 1 }
exit 0
