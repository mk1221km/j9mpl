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

# Go native compilation — no classpath needed
set env(GOFLAGS) "-mod=mod"

proc cleanUpMethodBlock {revisedBlock methodSig} {
    set revisedBlock [string trim $revisedBlock]
    
    # 1. Ensure method signature is present (Go uses 'func', NetRexx uses 'method')
    if {$methodSig != "" && ![regexp -nocase {^\s*(pub fn|func|method)\s+} $revisedBlock]} {
        puts "  -> Prepending missing function signature..."
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

set isZig [string match "*.zig" $nrxFile]
set isTest [string match "*_test.go" $nrxFile]
set runIncremental 0
if {!$isTest && [file exists [file join $projectDir ".context" "methods.txt"]]} {
    set nrxContent ""
    if {[file exists $nrxFile]} {
        set fd [open $nrxFile r]
        set nrxContent [read $fd]
        close $fd
    }
    if {[string match "*SKELETON_*" $nrxContent]} {
        set runIncremental 1
    } else {
        puts "\[INFO\] No skeleton markers found in $nrxFile. Skipping incremental synthesis."
    }
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
                # Build Go-targeted self-correction prompt from compiler errors
                set goError $compResult
                set scPrompt "The following Go function failed to compile:\n\n"
                append scPrompt "COMPILER ERROR:\n$goError\n\n"
                append scPrompt "Fix the function to compile successfully. Output ONLY the corrected Go function body.\n"
                append scPrompt "Use standard Go patterns: func (s *Type) Name(args) (result, error) { ... }\n"
                append scPrompt "Do NOT use Java or NetRexx syntax. Do NOT wrap in markdown code blocks.\n"

                set backoffDelay [expr {int(pow(2, $methodRetry) * 1000)}]
                puts "\[BACKOFF ACTIVE\]: Cooling pipeline for $backoffDelay ms before retrying..."
                after $backoffDelay

                puts "  -> Dispatching Go repair request to remote model..."
                set modelStatus [catch {exec bin/llm --print $scPrompt} modelRaw]
            }

            if {$modelStatus != 0} {
                puts stderr "\[ERROR\] Failed to execute llm binary: $modelRaw"
                exit 1
            }

            puts "  -> Raw model output length: [string length $modelRaw]"
            puts "  -> Raw model output:\n$modelRaw\n----------------------"

            # Extract block from code fences if present
            set normalizedRaw [string map {"```" "\u0000"} $modelRaw]
            set parts [split $normalizedRaw "\u0000"]
            if {[llength $parts] >= 3} {
                puts "  -> Extracting block from code fences..."
                set revisedBlock [lindex $parts 1]
                puts "  -> Extracted block:\n$revisedBlock\n----------------------"
                regsub -nocase {^(?:rexx|netrexx|go)?\n} $revisedBlock "" revisedBlock
            } else {
                puts "  -> Using raw model output..."
                set revisedBlock $modelRaw
            }

            # Strip illegal 'returns void' or 'returns none' from method signatures
            regsub -all -nocase {\s+returns\s+void\b} $revisedBlock "" revisedBlock
            regsub -all -nocase {\s+returns\s+none\b} $revisedBlock "" revisedBlock

            # Clean up the method block using our cleanUpMethodBlock helper
            set revisedBlock [cleanUpMethodBlock $revisedBlock $methodSig]

            # For Go target: rebuild the function body from scratch instead of patching
            # Read the full source file
            set fd [open $nrxFile r]
            set fullSource [read $fd]
            close $fd

            # Replace the SKELETON marker line(s) with the generated function body
            set marker "SKELETON_$method"
            set pattern "// $marker"
            if {[string first $pattern $fullSource] >= 0} {
                set fullSource [string map [list $pattern $revisedBlock] $fullSource]
            } else {
                puts stderr "\[ERROR\] SKELETON marker '$pattern' not found in source file."
                exit 1
            }

            set fd [open $nrxFile w]
            puts -nonewline $fd $fullSource
            close $fd
            puts "  -> Source file patched successfully."
            puts "--- SOURCE FILE CONTENT AFTER PATCHING $method ---"
            set fd_check [open $nrxFile r]
            puts [read $fd_check]
            close $fd_check
            puts "--------------------------------------------------"

            # Purge stale intermediate .java file
            set javaArtifact "[file rootname $nrxFile].java"
            file delete -force $javaArtifact

            # Execute compile check
            if {$isZig} {
                puts "  -> Compiling Zig source for method $method..."
                set compStatus [catch {exec zig build-obj $nrxFile --name ringbuffer} compResult]
            } else {
                puts "  -> Compiling and formatting Go source for method $method..."
                catch {exec goimports -w $nrxFile} _
                set compStatus [catch {exec go build -o /dev/null $nrxFile} compResult]
            }
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
    set maxCompileRetries 5
    set maxFuzzerRetries 5
    set compileRetryCount 0
    set fuzzerRetryCount 0
    set compiledSuccessfully 0

    while {($compileRetryCount < $maxCompileRetries) && ($fuzzerRetryCount < $maxFuzzerRetries) && !$compiledSuccessfully} {
        puts "=========================================================="
        puts "Iteration [expr {$compileRetryCount + $fuzzerRetryCount + 1}] (Compile attempts: $compileRetryCount/$maxCompileRetries, Fuzzer attempts: $fuzzerRetryCount/$maxFuzzerRetries)"
        puts "=========================================================="

        puts "  -> Purging stale intermediate artifacts..."
        set javaArtifact "[file rootname $nrxFile].java"
        file delete -force $javaArtifact

        puts "\[1/3\] Formatting and compiling Go source..."
        if {$isZig} {
            set status [catch {exec zig build-obj $nrxFile --name ringbuffer} result]
        } else {
            catch {exec go fmt $nrxFile} _
            set status [catch {exec go build -o /dev/null $nrxFile} result]
        }
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
                    
                    set failedMethod ""
                    if {[regexp {Assertion Failure in (\w+):} $sandboxResult match failedMethod]} {
                        puts "  -> Detected fuzzer assertion failure in method '$failedMethod'."
                        
                        set prodNrx [string map {_test.go .go} $nrxFile]
                        if {[file exists $prodNrx]} {
                            puts "  -> Extracting method '$failedMethod' from production file: $prodNrx"
                            set fd [open $prodNrx r]
                            set lines [split [read $fd] "\n"]
                            close $fd
                            
                            set inBlock 0
                            set originalBlockLines {}
                            foreach line $lines {
                                set trimmed [string trim $line]
                                if {[string match "--*" $trimmed] || [string match "/*" $trimmed]} {
                                    if {$inBlock} {
                                        lappend originalBlockLines $line
                                    }
                                    continue
                                }
                                if {[regexp -nocase {^\s*method\s+(\w+)} $line match mName]} {
                                    if {$inBlock} {
                                        break
                                    }
                                    if {[string equal -nocase $mName $failedMethod]} {
                                        set inBlock 1
                                        lappend originalBlockLines $line
                                        continue
                                    }
                                }
                                if {$inBlock} {
                                    lappend originalBlockLines $line
                                }
                            }
                            set originalBlock [join $originalBlockLines "\n"]
                            
                            # Build the fuzzer self-repair prompt
                            set promptBuilder ""
                            append promptBuilder "\[SYSTEM INVARIANT VIOLATION\]\n"
                            append promptBuilder "The property-based fuzzer executed boundary tests on the compiled bytecode and detected a runtime verification failure.\n\n"
                            append promptBuilder "Target Method: $failedMethod\n"
                            append promptBuilder "Fuzzer Diagnostic:\n$sandboxResult\n\n"
                            append promptBuilder "Original NetRexx Code Context:\n"
                            append promptBuilder "```rexx\n"
                            append promptBuilder "$originalBlock\n"
                            append promptBuilder "```\n\n"
                            append promptBuilder "CORRECTION INSTRUCTIONS:\n"
                            append promptBuilder "1. You MUST validate all inputs (including DTO record fields) at the very beginning of the method before any JDBC or logic execution.\n"
                            append promptBuilder "2. To prevent NullPointerException on null parameters, you MUST isolate checks inside explicit, nested 'if' statements instead of combining them with logical OR '|' or AND '&' operators, because NetRexx logical operators do not short-circuit.\n"
                            append promptBuilder "   Example: \n"
                            append promptBuilder "     if dbPath \\== null then if dbPath.trim().length() == 0 then do\n"
                            append promptBuilder "       signal java.lang.IllegalArgumentException(\"Invalid input\")\n"
                            append promptBuilder "     end\n"
                            append promptBuilder "3. ALWAYS use strict reference comparison identity operators '== null' and '\\== null' when checking for Java null pointers. Never use '=' or '\\=' for null checking.\n"
                            append promptBuilder "4. Throw the expected exception types (java.lang.IllegalArgumentException, java.io.IOException, java.sql.SQLException, java.lang.NumberFormatException) according to the validation invariants.\n"
                            append promptBuilder "Output ONLY the revised, complete NetRexx source block starting with 'method $failedMethod'. Do not wrap in markdown fences or include explanations.\n"
                            
                            set promptFile [file join $projectDir ".context" "self_correct_prompt.txt"]
                            set fd [open $promptFile w]
                            puts -nonewline $fd $promptBuilder
                            close $fd
                            
                            set originalBlockFile [file join $projectDir ".context" "errant_block_original.txt"]
                            set fd [open $originalBlockFile w]
                            puts -nonewline $fd $originalBlock
                            close $fd
                            
                            incr fuzzerRetryCount
                            if {$fuzzerRetryCount >= $maxFuzzerRetries} {
                                puts stderr "\[CIRCUIT BREAKER FAILURE\]: Fuzzer repair retry ceiling breached ($maxFuzzerRetries turns). Execution frozen."
                                exit 1
                            }
                            
                            set backoffDelay [expr {int(pow(2, $fuzzerRetryCount) * 1000)}]
                            puts "\[BACKOFF ACTIVE\]: Cooling pipeline for $backoffDelay ms before dispatching fuzzer repair request..."
                            after $backoffDelay
                            
                            puts "  -> Dispatching fuzzer repair request to remote model..."
                            set modelStatus [catch {exec bin/llm --print $promptBuilder} modelRaw]
                            if {$modelStatus != 0} {
                                puts stderr "\[ERROR\] Failed to execute llm: $modelRaw"
                                exit 1
                            }
                            
                            # Clean up and extract block
                            set normalizedRaw [string map {"```" "\u0000"} $modelRaw]
                            set parts [split $normalizedRaw "\u0000"]
                            if {[llength $parts] >= 3} {
                                set revisedBlock [lindex $parts 1]
                                regsub -nocase {^(?:rexx|netrexx|go)?\n} $revisedBlock "" revisedBlock
                            } else {
                                set revisedBlock $modelRaw
                            }
                            
                            regsub -all -nocase {\s+returns\s+void\b} $revisedBlock "" revisedBlock
                            regsub -all -nocase {\s+returns\s+none\b} $revisedBlock "" revisedBlock
                            
                            # Try to extract signature
                            set origSig ""
                            foreach line [split $originalBlock "\n"] {
                                set trimLine [string trim $line]
                                if {[regexp -nocase {^\s*method\s+} $trimLine]} {
                                    set origSig $trimLine
                                    break
                                }
                            }
                            set revisedBlock [cleanUpMethodBlock $revisedBlock $origSig]
                            
                            set revisedBlockFile [file join $projectDir ".context" "errant_block_revised.txt"]
                            set fd [open $revisedBlockFile w]
                            puts -nonewline $fd $revisedBlock
                            close $fd
                            
                            puts "  -> Applying repair patch to production file: $prodNrx"
                            set patchStatus [catch {exec bin/patch_source $prodNrx $originalBlockFile $revisedBlockFile} patchResult]
                            puts $patchResult
                            if {$patchStatus != 0} {
                                puts stderr "\[ERROR\] Failed to patch production file: $patchResult"
                                exit 1
                            }
                            
                            puts "  -> Re-compiling production file..."
                            if {$isZig} {
                                set compStatus [catch {exec zig build-obj $prodNrx --name ringbuffer} compResult]
                            } else {
                                catch {exec goimports -w $prodNrx} _
                                set compStatus [catch {exec go build -o /dev/null $prodNrx} compResult]
                            }
                            puts $compResult
                            
                            set compiledSuccessfully 0
                            continue
                        }
                    }
                    exit 1
                }
                puts "  -> Sandboxed execution sweep completed successfully."
            }

            puts "=========================================================="
            puts " \[SUCCESS\] Zero-error build achieved!"
            puts "=========================================================="
            exit 0
        }

        incr compileRetryCount
        if {$compileRetryCount >= $maxCompileRetries} {
            puts stderr "\[CIRCUIT BREAKER FAILURE\]: Compilation retry ceiling breached ($maxCompileRetries turns). Execution frozen."
            exit 1
        }

        set backoffDelay [expr {int(pow(2, $compileRetryCount) * 1000)}]
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

        puts "  -> Raw model output length: [string length $modelRaw]"
        puts "  -> Raw model output:\n$modelRaw\n----------------------"

        set normalizedRaw [string map {"```" "\u0000"} $modelRaw]
        set parts [split $normalizedRaw "\u0000"]
        if {[llength $parts] >= 3} {
            puts "  -> Extracting block from code fences..."
            set revisedBlock [lindex $parts 1]
            puts "  -> Extracted block:\n$revisedBlock\n----------------------"
            regsub -nocase {^(?:rexx|netrexx|go)?\n} $revisedBlock "" revisedBlock
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
