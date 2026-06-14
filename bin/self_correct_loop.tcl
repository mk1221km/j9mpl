#!/usr/bin/env tclsh
# =================================================================
# Autonomous Self-Correction Loop for NetRexx / ECJ Toolchain (Tcl)
# =================================================================

if {[llength $argv] < 1} {
    puts "Usage: [file tail [info script]] <target_nrx_file>"
    exit 1
}

set nrxFile [file normalize [lindex $argv 0]]
set projectDir [file dirname [file dirname [file normalize [info script]]]]
set classpath "bin:lib/NetRexxF.jar:target/dependency/sqlite-jdbc-3.45.1.0.jar:target/dependency/slf4j-api-1.7.36.jar"

# Set CLASSPATH env variable for the child processes
set env(CLASSPATH) "bin"

# Navigate to project root directory
cd $projectDir

set maxRetries 5
set retryCount 0

while {$retryCount < $maxRetries} {
    puts "=========================================================="
    puts "Iteration [expr {$retryCount + 1}] / $maxRetries"
    puts "=========================================================="

    puts "\[1/3\] Executing compiler validation..."
    
    # Run self_correct.
    set status [catch {exec bin/self_correct $nrxFile -cp $classpath} result]
    puts $result

    if {$status == 0} {
        set isTest [string match "*Test.nrx" $nrxFile]
        if {$isTest} {
            puts "\[1.5/3\] Running sandboxed verification sweep..."
            set className [file rootname [file tail $nrxFile]]
            set fd [open $nrxFile r]
            set content [read $fd]
            close $fd
            set pkg ""
            if {[regexp -line {^\s*package\s+([a-zA-Z0-9_\.]+)} $content match pkgName]} {
                set pkg "$pkgName."
            }
            set targetClass "$pkg$className"
            set sandboxCmd [list systemd-run --user --scope --description=Factory-Fuzzer-Sandbox -p MemoryMax=512M -p CPUQuota=50% -p TasksMax=100 ./bin/sandbox_exec.sh $targetClass]
            set sandboxStatus [catch {exec {*}$sandboxCmd 2>@1} sandboxResult]
            if {$sandboxStatus != 0} {
                puts "\[ERROR\] Sandboxed execution sweep failed: $sandboxResult"
                exit 1
            }
            puts "  -> Sandboxed execution sweep completed successfully."
        }

        puts "=========================================================="
        puts " \[SUCCESS\] Zero-error build achieved!"
        puts "=========================================================="
        exit 0
    }

    puts "\[2/3\] Intercepting compiler diagnostic and generating prompt..."
    
    # Read the generated prompt
    set promptFile [file join $projectDir ".context" "self_correct_prompt.txt"]
    if {![file exists $promptFile]} {
        puts "\[ERROR\] Prompt file does not exist. Aborting loop."
        exit 1
    }
    
    set fd [open $promptFile r]
    set prompt [read $fd]
    close $fd

    puts "  -> Dispatching self-correction prompt to remote model..."
    
    # Call remote model
    set modelStatus [catch {exec agy --print $prompt} modelRaw]
    if {$modelStatus != 0} {
        puts "\[ERROR\] Failed to execute agy CLI: $modelRaw"
        exit 1
    }

    # Extract block from code fences if present
    set parts [split $modelRaw "```"]
    if {[llength $parts] >= 3} {
        puts "  -> Extracting block from code fences..."
        set revisedBlock [lindex $parts 1]
        # Strip any leading 'rexx\n' or similar language identifier
        regsub -nocase {^(?:rexx)?\n} $revisedBlock "" revisedBlock
    } else {
        puts "  -> Using raw model output..."
        set revisedBlock $modelRaw
    }

    puts "\[3/3\] Applying patch to source file..."
    
    set originalBlockFile [file join $projectDir ".context" "errant_block_original.txt"]
    set revisedBlockFile [file join $projectDir ".context" "errant_block_revised.txt"]
    
    # Save revised block to disk for patch_source
    set fd [open $revisedBlockFile w]
    puts -nonewline $fd $revisedBlock
    close $fd

    # Call patch_source binary
    set patchStatus [catch {exec bin/patch_source $nrxFile $originalBlockFile $revisedBlockFile} patchResult]
    puts $patchResult
    if {$patchStatus != 0} {
        puts "\[ERROR\] Failed to patch source file: $patchResult"
        exit 1
    }

    puts "  -> Source file patched. Retrying compilation..."
    incr retryCount
}

puts "=========================================================="
puts " \[ERROR\] Failed to converge after $maxRetries attempts."
puts "=========================================================="
exit 1

