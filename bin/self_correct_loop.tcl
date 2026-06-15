#!/usr/bin/env tclsh
# =================================================================
# Autonomous Self-Correction & Incremental Synthesis Loop (Tcl)
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

proc cleanUpMethodBlock {revisedBlock methodSig} {
    set revisedBlock [string trim $revisedBlock]
    
    # 1. Ensure method signature is present
    if {$methodSig != "" && ![regexp -nocase {^\s*method\s+} $revisedBlock]} {
        puts "  -> Prepending missing method signature..."
        set revisedBlock "$methodSig\n$revisedBlock"
    }
    
    # 2. Count opens and closes to strip extra trailing end
    set opens 0
    foreach line [split $revisedBlock "\n"] {
        set trimmedLine [string trim $line]
        if {[string match "--*" $trimmedLine] || [string match "/**" $trimmedLine]} {
            continue
        }
        set opens [expr {$opens + [regexp -all -nocase {\bdo\b} $line]}]
        set opens [expr {$opens + [regexp -all -nocase {\bloop\b} $line]}]
        set opens [expr {$opens + [regexp -all -nocase {\bselect\b} $line]}]
    }
    
    set closes 0
    foreach line [split $revisedBlock "\n"] {
        set trimmedLine [string trim $line]
        if {[string match "--*" $trimmedLine] || [string match "/**" $trimmedLine]} {
            continue
        }
        set closes [expr {$closes + [regexp -all -nocase {\bend\b} $line]}]
    }
    
    puts "  -> Open blocks: $opens, Close blocks: $closes"
    if {$closes > $opens} {
        # Check if last non-empty line is 'end'
        set lines [split $revisedBlock "\n"]
        set lastIdx [expr {[llength $lines] - 1}]
        while {$lastIdx >= 0 && [string trim [lindex $lines $lastIdx]] == ""} {
            incr lastIdx -1
        }
        if {$lastIdx >= 0} {
            set lastLine [string trim [lindex $lines $lastIdx]]
            if {[regexp -nocase {^end\b} $lastLine]} {
                puts "  -> Stripping extra trailing end: '$lastLine'"
                set lines [lreplace $lines $lastIdx $lastIdx]
                set revisedBlock [join $lines "\n"]
            }
        }
    }
    
    return $revisedBlock
}

# Navigate to project root directory
cd $projectDir

set isTest [string match "*Test.nrx" $nrxFile]
set runIncremental 0
if {!$isTest && [file exists [file join $projectDir ".context" "methods.txt"]]} {
    set runIncremental 1
}

if {$runIncremental} {
    puts "=========================================================="
    puts "Starting Incremental Assembly-Line Synthesis for:"
    puts "  $nrxFile"
    puts "=========================================================="

    # Read methods list
    set fd [open [file join $projectDir ".context" "methods.txt"] r]
    set methods [split [read $fd] "\n"]
    close $fd

    foreach method $methods {
        set method [string trim $method]
        if {$method == ""} continue

        puts "----------------------------------------------------------"
        puts "Synthesizing method: $method"
        puts "----------------------------------------------------------"

        set skeletonFile [file join $projectDir ".context" "skeleton_$method.txt"]
        set promptFile [file join $projectDir ".context" "prompt_$method.txt"]
        
        if {![file exists $skeletonFile] || ![file exists $promptFile]} {
            puts stderr "\[ERROR\] Missing skeleton or prompt file for method $method"
            exit 1
        }

        # Read original skeleton block
        set fd [open $skeletonFile r]
        set originalBlock [read $fd]
        close $fd

        set methodSig [string trim [lindex [split $originalBlock "\n"] 0]]

        # Read prompt
        set fd [open $promptFile r]
        set prompt [read $fd]
        close $fd

        # Start with the baseline skeleton block for patching
        set origFile [file join $projectDir ".context" "errant_block_original.txt"]
        set fd [open $origFile w]
        puts -nonewline $fd $originalBlock
        close $fd

        set methodRetry 0
        set methodMaxRetries 3
        set methodSynthesized 0

        while {$methodRetry < $methodMaxRetries && !$methodSynthesized} {
            if {$methodRetry == 0} {
                puts "  -> Dispatching method synthesis to remote model..."
                set modelStatus [catch {exec bin/llm --print $prompt} modelRaw]
            } else {
                # If we are retrying, read the self-correction prompt written by self_correct Go binary
                set scPromptFile [file join $projectDir ".context" "self_correct_prompt.txt"]
                if {![file exists $scPromptFile]} {
                    puts stderr "\[ERROR\] Self-correction prompt file missing for retry"
                    exit 1
                }
                set fd [open $scPromptFile r]
                set scPrompt [read $fd]
                close $fd

                set backoffDelay [expr {int(pow(2, $methodRetry) * 1000)}]
                puts "\[BACKOFF ACTIVE\]: Cooling pipeline for $backoffDelay ms before retrying..."
                after $backoffDelay

                puts "  -> Dispatching method repair request to remote model..."
                set modelStatus [catch {exec bin/llm --print $scPrompt} modelRaw]
            }

            if {$modelStatus != 0} {
                puts stderr "\[ERROR\] Failed to execute llm binary: $modelRaw"
                exit 1
            }

            # Extract block from code fences if present
            set parts [split $modelRaw "```"]
            if {[llength $parts] >= 3} {
                puts "  -> Extracting block from code fences..."
                set revisedBlock [lindex $parts 1]
                regsub -nocase {^(?:rexx)?\n} $revisedBlock "" revisedBlock
            } else {
                puts "  -> Using raw model output..."
                set revisedBlock $modelRaw
            }

            # Strip illegal 'returns void' or 'returns none' from method signatures
            regsub -all -nocase {\s+returns\s+void\b} $revisedBlock "" revisedBlock
            regsub -all -nocase {\s+returns\s+none\b} $revisedBlock "" revisedBlock

            # Clean up the method block using our cleanUpMethodBlock helper
            set revisedBlock [cleanUpMethodBlock $revisedBlock $methodSig]

            # Write revised block to errant_block_revised.txt
            set revFile [file join $projectDir ".context" "errant_block_revised.txt"]
            set fd [open $revFile w]
            puts -nonewline $fd $revisedBlock
            close $fd

            # Patch the source file
            set patchStatus [catch {exec bin/patch_source $nrxFile $origFile $revFile} patchResult]
            puts $patchResult
            if {$patchStatus != 0} {
                puts stderr "\[ERROR\] Failed to patch source file: $patchResult"
                exit 1
            }
            puts "--- SOURCE FILE CONTENT AFTER PATCHING $method ---"
            set fd_check [open $nrxFile r]
            puts [read $fd_check]
            close $fd_check
            puts "--------------------------------------------------"

            # Purge stale intermediate .java file
            set javaArtifact "[file rootname $nrxFile].java"
            file delete -force $javaArtifact

            # Execute compile check
            puts "  -> Executing compile check for method $method..."
            set compStatus [catch {exec bin/self_correct $nrxFile -cp $classpath} compResult]
            puts $compResult

            if {$compStatus == 0} {
                set methodSynthesized 1
                puts "  -> Method $method compiled successfully!"
            } else {
                puts "  -> Method $method failed compilation. Initiating repair turn..."
                # self_correct Go binary generates errant_block_original.txt in .context/
                # We use it directly as the new origFile baseline for the next patch iteration
                incr methodRetry
            }
        }

        if {!$methodSynthesized} {
            puts stderr "=========================================================="
            puts stderr " \[ERROR\] Failed to synthesize method $method after $methodMaxRetries attempts."
            puts stderr "=========================================================="
            exit 1
        }
    }

    puts "=========================================================="
    puts " \[SUCCESS\] All methods synthesized and compiled cleanly!"
    puts "=========================================================="
    exit 0

} else {
    # Standard compilation and self-correction loop (for Test files or fallback)
    set maxRetries 3
    set retryCount 0
    set compiledSuccessfully 0

    while {$retryCount < $maxRetries && !$compiledSuccessfully} {
        puts "=========================================================="
        puts "Iteration [expr {$retryCount + 1}] / $maxRetries"
        puts "=========================================================="

        puts "  -> Purging stale intermediate artifacts..."
        set javaArtifact "[file rootname $nrxFile].java"
        file delete -force $javaArtifact

        puts "\[1/3\] Executing compiler validation..."
        set status [catch {exec bin/self_correct $nrxFile -cp $classpath} result]
        puts $result

        if {$status == 0} {
            set compiledSuccessfully 1
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
                    puts stderr "\[ERROR\] Sandboxed execution sweep failed: $sandboxResult"
                    exit 1
                }
                puts "  -> Sandboxed execution sweep completed successfully."
            }

            puts "=========================================================="
            puts " \[SUCCESS\] Zero-error build achieved!"
            puts "=========================================================="
            exit 0
        }

        incr retryCount
        if {$retryCount >= $maxRetries} {
            puts stderr "\[CIRCUIT BREAKER FAILURE\]: Compilation retry ceiling breached ($maxRetries turns). Execution frozen."
            exit 1
        }

        set backoffDelay [expr {int(pow(2, $retryCount) * 1000)}]
        puts "\[BACKOFF ACTIVE\]: Cooling pipeline for $backoffDelay ms before dispatching repair request..."
        after $backoffDelay

        puts "\[2/3\] Intercepting compiler diagnostic and generating prompt..."
        set promptFile [file join $projectDir ".context" "self_correct_prompt.txt"]
        if {![file exists $promptFile]} {
            puts stderr "\[ERROR\] Prompt file does not exist. Aborting loop."
            exit 1
        }
        
        set fd [open $promptFile r]
        set prompt [read $fd]
        close $fd

        puts "  -> Dispatching self-correction prompt to remote model..."
        set modelStatus [catch {exec bin/llm --print $prompt} modelRaw]
        if {$modelStatus != 0} {
            puts stderr "\[ERROR\] Failed to execute llm: $modelRaw"
            exit 1
        }

        set parts [split $modelRaw "```"]
        if {[llength $parts] >= 3} {
            puts "  -> Extracting block from code fences..."
            set revisedBlock [lindex $parts 1]
            regsub -nocase {^(?:rexx)?\n} $revisedBlock "" revisedBlock
        } else {
            puts "  -> Using raw model output..."
            set revisedBlock $modelRaw
        }

        # Strip illegal 'returns void' or 'returns none' from method signatures
        regsub -all -nocase {\s+returns\s+void\b} $revisedBlock "" revisedBlock
        regsub -all -nocase {\s+returns\s+none\b} $revisedBlock "" revisedBlock

        # Try to extract signature from originalBlockFile
        set originalBlockFile [file join $projectDir ".context" "errant_block_original.txt"]
        set origSig ""
        if {[file exists $originalBlockFile]} {
            set fd_orig [open $originalBlockFile r]
            set origContent [read $fd_orig]
            close $fd_orig
            foreach line [split $origContent "\n"] {
                set trimLine [string trim $line]
                if {[regexp -nocase {^\s*method\s+} $trimLine]} {
                    set origSig $trimLine
                    break
                }
            }
        }

        # Clean up the method block using our cleanUpMethodBlock helper
        set revisedBlock [cleanUpMethodBlock $revisedBlock $origSig]

        puts "\[3/3\] Applying patch to source file..."
        set revisedBlockFile [file join $projectDir ".context" "errant_block_revised.txt"]
        
        set fd [open $revisedBlockFile w]
        puts -nonewline $fd $revisedBlock
        close $fd

        set patchStatus [catch {exec bin/patch_source $nrxFile $originalBlockFile $revisedBlockFile} patchResult]
        puts $patchResult
        if {$patchStatus != 0} {
            puts stderr "\[ERROR\] Failed to patch source file: $patchResult"
            exit 1
        }

        puts "  -> Source file patched. Retrying compilation..."
    }
}
exit 1
