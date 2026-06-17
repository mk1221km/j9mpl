#!/usr/bin/env tclsh
# Chaos Injector — forces stateless binaries through environmental anomalies.
# No per-line acknowledgments. Validate via query after corruption.
# Same pipe protocol as base tests: batch input, read output.

set metrics_target "generated/metrics/metrics.py"
set router_target "generated/router/router.py"

if {![file executable $metrics_target]} {
    puts stderr "ERROR: $metrics_target not found or not executable"
    exit 1
}
if {![file executable $router_target]} {
    puts stderr "ERROR: $router_target not found or not executable"
    exit 1
}

set total_passed 0
set total_failed 0

proc batch_test {target input expected_lines} {
    set fd [open "|$target" r+]
    puts $fd $input
    close $fd w
    set output [split [string trim [read $fd]] "\n"]
    close $fd

    set ok 1
    set i 0
    foreach exp $expected_lines {
        set actual [lindex $output $i]
        if {$actual ne $exp} {
            puts "    line $i: expected '$exp', got '$actual'"
            set ok 0
        }
        incr i
    }
    return $ok
}

proc chaos_assert {test_name result} {
    global total_passed total_failed
    if {$result} {
        puts "  PASS: $test_name"
        incr total_passed
    } else {
        puts "  FAIL: $test_name"
        incr total_failed
    }
}

# ============================================================
# METRICS LOGGER CHAOS TESTS
# ============================================================
puts "=== Metrics Logger Chaos Battery ==="

# Chaos 1: Corrupt line mid-stream — avg after corruption should be correct
puts "Chaos 1: Interleaved corrupt lines..."
set result [batch_test $metrics_target "log cpu 50.0\nlog\nlog cpu 60.0\navg cpu" "55.0"]
chaos_assert "avg correct after corrupt log" $result

# Chaos 2: IEEE 754 non-finite values silently dropped
puts "Chaos 2: IEEE 754 NaN/Inf injection..."
set result [batch_test $metrics_target "log mem NaN\nlog mem Inf\nlog mem -Inf\nlog mem 10.0\navg mem" "10.0"]
chaos_assert "NaN/Inf silently dropped, avg unaffected" $result

# Chaos 3: Long metric name
puts "Chaos 3: Long metric name..."
set long_name [string repeat "x" 4096]
set result [batch_test $metrics_target "log $long_name 42.0\nlog $long_name 84.0\navg $long_name" "63.0"]
chaos_assert "long name no crash, avg correct" $result

# Chaos 4: Large stream (10K metrics)
puts "Chaos 4: 10,000 metric stream..."
set input ""
for {set i 0} {$i < 10000} {incr i} {
    append input "log chaos [expr {$i % 100}]\n"
}
append input "count chaos\navg chaos"
chaos_assert "10K count+avg correct" [batch_test $metrics_target $input [list "10000" "49.5"]]

# Chaos 5: Alternating valid/invalid values
puts "Chaos 5: NaN filters from aggregation..."
chaos_assert "NaN excluded, valid logs preserved" [batch_test $metrics_target "log a 10.0\nlog b NaN\nlog a 20.0\nlog b 30.0\navg a\ncount b" [list "15.0" "1"]]

# ============================================================
# TRANSACTION ROUTER CHAOS TESTS
# ============================================================
puts "=== Transaction Router Chaos Battery ==="

# Chaos 6: Epsilon boundary around 1000
puts "Chaos 6: Epsilon boundary \$1000..."
set result [batch_test $router_target "route t001 low 999.999\nroute t002 low 1000.0\nroute t003 low 1000.001" [list \
    "routed t001 CLEAR_ACH" \
    "routed t002 CLEAR_ACH" \
    "routed t003 HIGH_PRIORITY_WIRE" \
]]
chaos_assert "999.999->ACH, 1000->ACH, 1000.001->WIRE" $result

# Chaos 7: Priority casing permutations
puts "Chaos 7: Priority casing..."
set result [batch_test $router_target "route t001 HIGH 500\nroute t002 High 500\nroute t003 high 500\nroute t004 hIgH 500" [list \
    "routed t001 HIGH_PRIORITY_WIRE" \
    "routed t002 HIGH_PRIORITY_WIRE" \
    "routed t003 HIGH_PRIORITY_WIRE" \
    "routed t004 HIGH_PRIORITY_WIRE" \
]]
chaos_assert "case-insensitive HIGH" $result

# Chaos 8: Truncated/malformed route lines
puts "Chaos 8: Truncated route lines..."
set result [batch_test $router_target "route t001\nroute t002 low\nroute t003 low abc" [list "INVALID" "INVALID" "INVALID"]]
chaos_assert "truncated routes rejected" $result

# Chaos 9: Extreme amounts
puts "Chaos 9: Extreme amounts..."
set result [batch_test $router_target "route t001 low 1e10\nroute t002 low 0.0000001" [list \
    "routed t001 HIGH_PRIORITY_WIRE" \
    "routed t002 CLEAR_ACH" \
]]
chaos_assert "extreme amounts routed correctly" $result

# Chaos 10: ACH type routing
puts "Chaos 10: ACH type routing..."
set result [batch_test $router_target "route t001 ACH 999.999\nroute t002 ACH 1000.001" [list \
    "routed t001 CLEAR_ACH" \
    "routed t002 HIGH_PRIORITY_WIRE" \
]]
chaos_assert "ACH type epsilon boundary" $result

# Chaos 11: Empty input
puts "Chaos 11: Empty input..."
set result [batch_test $router_target "" {}]
chaos_assert "empty input" $result

# Chaos 12: Unknown commands interleaved
puts "Chaos 12: Unknown commands..."
set result [batch_test $router_target "route t001 low 500\nblargh\nroute t002 high 500\nrate" [list \
    "routed t001 CLEAR_ACH" \
    "routed t002 HIGH_PRIORITY_WIRE" \
    "0.0" \
]]
chaos_assert "unknown commands ignored" $result

# ============================================================
puts "========================================"
puts "Chaos injection complete."
puts "  Passed: $total_passed"
puts "  Failed: $total_failed"
puts "========================================"
if {$total_failed > 0} { exit 1 }
exit 0
