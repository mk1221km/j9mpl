#!/usr/bin/env tclsh
# Minimal behavioral TDD loop: generate → compile → test → feedback
# No skeleton generators, no type guardrails, no compiler parsing

if {[llength $argv] < 1} {
    puts "Usage: [file tail [info script]] <spec_path>"
    exit 1
}

set specPath [file normalize [lindex $argv 0]]
set projectDir [file dirname [file dirname [file normalize [info script]]]]
set specName [file rootname [file tail $specPath]]
set srcFile [file join $projectDir "generated" "${specName}.go"]
set binFile [file join $projectDir "generated" $specName]

# Build the prompt
set promptFile [file join $projectDir ".context" "prompt_${specName}.txt"]
set fd [open $promptFile r]
set prompt [read $fd]
close $fd

# Add behavioral test instruction
append prompt "\n\nWrite a complete Go program that reads commands from stdin and writes results to stdout.\n"
append prompt "Commands: push <value>, avg, readRange <n>, utilization\n"
append prompt "Output format: one value per line, formatted to 1 decimal place.\n"
append prompt "Output ONLY the code. No explanations, no markdown fences.\n"

set maxAttempts 5

for {set attempt 1} {$attempt <= $maxAttempts} {incr attempt} {
    puts "Attempt $attempt/$maxAttempts"

    # Generate code
    puts "  -> Dispatching to model..."
    set status [catch {exec bin/llm --print $prompt} code]
    if {$status != 0} {
        puts stderr "  [ERROR] LLM failed: $code"
        exit 1
    }
    puts "  -> Received [string length $code] chars"

    # Write source file
    set fd [open $srcFile w]
    puts -nonewline $fd $code
    close $fd

    # Compile
    puts "  -> Compiling..."
    set compileStatus [catch {exec go build -o $binFile $srcFile} compileErr]
    if {$compileStatus != 0} {
        puts "  -> Compile failed. Feeding error back..."
        append prompt "\n\nCOMPILATION ERROR:\n$compileErr\n\nFix the error and output the complete corrected Go program.\n"
        continue
    }
    puts "  -> Compile OK"

    # Run behavioral test
    puts "  -> Running behavioral tests..."
    set testStatus [catch {exec tclsh [file join $projectDir "bin" "test_harness.tcl"] $binFile} testOut]
    
    if {$testStatus == 0} {
        puts "  -> ALL TESTS PASSED"
        puts "=========================================================="
        puts " [SUCCESS] ${specName} passes behavioral validation"
        puts "=========================================================="
        exit 0
    }

    puts "  -> Test failed. Feeding output back..."
    append prompt "\n\nTEST FAILURE:\n$testOut\n\nFix the program to pass all tests and output the complete corrected Go program.\n"
}

puts "=========================================================="
puts " [FAIL] ${specName} did not converge after $maxAttempts attempts"
puts "=========================================================="
exit 1
