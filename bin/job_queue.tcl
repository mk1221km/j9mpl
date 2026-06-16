#!/usr/bin/env tclsh
# =================================================================
# Parallel Job Queue Supervisor for Specification Ingestion
# =================================================================

if {[llength $argv] < 1} {
    puts "Usage: [file tail [info script]] <spec1.md> \[spec2.md ...\]"
    exit 1
}

set projectDir [file dirname [file dirname [file normalize [info script]]]]
cd $projectDir

# Enforce local bin precedence across all fork and exec boundaries
set env(PATH) "$projectDir/bin:$env(PATH)"

set maxWorkers 1
set activeCount 0
set queue $argv
set exitCode 0

array set active {}

# Setup jobs folder
file mkdir [file join $projectDir "jobs"]

proc verifyEnvironmentInvariants {projectDir} {
    set invariants [list \
        [file join $projectDir "lib" "NetRexxF.jar"] \
        [file join $projectDir "ecj-3.46.0.jar"] \
        [file join $projectDir "bin" "llm"] \
    ]
    
    foreach invariant $invariants {
        if {![file exists $invariant]} {
            return -code error "CRITICAL EXCEPTION: Starvation of core component -> $invariant"
        }
    }
    
    set depDir [file join $projectDir "target" "dependency"]
    if {![file exists [file join $depDir "sqlite-jdbc-3.45.1.0.jar"]] || \
        ![file exists [file join $depDir "slf4j-api-1.7.36.jar"]]} {
        return -code error "CRITICAL EXCEPTION: Runtime dependencies absent in target/dependency/"
    }
    return 1
}

proc setupWorkspace {className specPath} {
    global projectDir
    
    # Pre-flight environment integrity check
    if {[catch {verifyEnvironmentInvariants $projectDir} err]} {
        puts stderr "$err"
        exit 1
    }
    
    puts "  -> Setting up isolated workspace for $className..."
    set jobDir [file join $projectDir "jobs" $className]
    
    # Clean up previous directory if it exists
    if {[file exists $jobDir]} {
        cleanupWorkspace $className
    }
    
    file mkdir $jobDir
    file mkdir [file join $jobDir ".context"]
    file mkdir [file join $jobDir "generated"]
    file mkdir [file join $jobDir "bin"]
    
    # Copy SQLite database
    file copy [file join $projectDir ".context" "project_context.db"] [file join $jobDir ".context" "project_context.db"]
    
    # Copy spec markdown file
    file copy $specPath [file join $jobDir "generated" [file tail $specPath]]
    
    # Copy existing synthesized class file if present to preserve manual refinement
    set srcNrx [file join $projectDir "generated" "$className.nrx"]
    if {[file exists $srcNrx]} {
        file copy $srcNrx [file join $jobDir "generated" "$className.nrx"]
    }
    
    # Link read-only directories
    file link [file join $jobDir "lib"] [file join $projectDir "lib"]
    file link [file join $jobDir "target"] [file join $projectDir "target"]
    file link [file join $jobDir "rascal-shell-stable.jar"] [file join $projectDir "rascal-shell-stable.jar"]
    file link [file join $jobDir "ecj-3.46.0.jar"] [file join $projectDir "ecj-3.46.0.jar"]
    file link [file join $jobDir "src"] [file join $projectDir "src"]
    file link [file join $jobDir "scratch"] [file join $projectDir "scratch"]
    
    # Copy all files from main bin/ to job bin/ (excluding the com directory) to preserve isolated path resolution
    foreach f [glob -nocomplain -tails -directory [file join $projectDir "bin"] *] {
        if {$f != "com"} {
            file copy [file join $projectDir "bin" $f] [file join $jobDir "bin" $f]
        }
    }
    
    # Copy class file structure for write isolation
    set comPath [file join $projectDir "bin" "com"]
    if {[file exists $comPath]} {
        exec cp -r $comPath [file join $jobDir "bin"]
    }
}

proc mergeWorkspace {className} {
    global projectDir
    puts "  -> Merging $className artifacts back to main workspace..."
    set jobDir [file join $projectDir "jobs" $className]
    
    # Copy generated files back
    catch {exec cp -a [file join $jobDir "generated"]/. [file join $projectDir "generated"]}
    # Copy compiled classes back
    catch {exec cp -a [file join $jobDir "bin" "com"]/. [file join $projectDir "bin" "com"]}
}

proc cleanupWorkspace {className} {
    global projectDir
    puts "  -> Cleaning up workspace for $className..."
    set jobDir [file join $projectDir "jobs" $className]
    catch {file delete -force $jobDir}
}

proc updateMainLedger {} {
    global projectDir
    puts "\[INFO\] Re-indexing main ledger database..."
    set extractCmd [list java -cp "rascal-shell-stable.jar:target/dependency/*" org.rascalmpl.shell.RascalShell ContextExtractor $projectDir [file join $projectDir ".context" "extracted.sql"]]
    if {[catch {exec {*}$extractCmd 2>@1} msg]} {
        puts "\[WARNING\] ContextExtractor failed during re-index: $msg"
    } else {
        set dbPath [file join $projectDir ".context" "project_context.db"]
        set sqlPath [file join $projectDir ".context" "extracted.sql"]
        if {[catch {exec sqlite3 $dbPath < $sqlPath 2>@1} msg]} {
            puts "\[WARNING\] SQLite import failed during re-index: $msg"
        } else {
            puts "\[INFO\] Main ledger database re-indexed successfully."
        }
    }
}

proc handleOutput {pipe className} {
    global active activeCount exitCode
    if {[gets $pipe line] >= 0} {
        puts "\[$className\]: $line"
    }
    if {[eof $pipe]} {
        fconfigure $pipe -blocking 1
        set status [catch {close $pipe} err]
        unset active($pipe)
        incr activeCount -1
        
        if {$status == 0} {
            puts "\[SUCCESS\] Job $className finished successfully!"
            mergeWorkspace $className
            cleanupWorkspace $className
            updateMainLedger
        } else {
            puts "\[ERROR\] Job $className failed: $err"
            cleanupWorkspace $className
            set exitCode 1
        }
        startNextJob
    }
}

proc startNextJob {} {
    global queue maxWorkers activeCount active projectDir
    
    if {[llength $queue] == 0 && $activeCount == 0} {
        puts "=========================================================="
        puts "All jobs completed."
        puts "=========================================================="
        global exitCode
        exit $exitCode
    }
    
    if {$activeCount >= $maxWorkers || [llength $queue] == 0} {
        return
    }
    
    # Pop from queue
    set specPath [file normalize [lindex $queue 0]]
    set queue [lrange $queue 1 end]
    
    # Determine class name
    set specName [file tail $specPath]
    if {![regexp {^(\w+)Spec\.md$} $specName match className]} {
        set className [file rootname $specName]
    }
    
    setupWorkspace $className $specPath
    
    puts "\[INFO\] Spawning job for $className..."
    
    # Open async pipe to run the pipeline inside the isolated workspace directory
    set cmd [list ./bin/run_job_pipeline.sh [file join "generated" [file tail $specPath]]]
    set pipeStatus [catch {open "| bash -c \"cd [file join jobs $className] && $cmd\" 2>@1" r} pipe]
    
    if {$pipeStatus != 0} {
        puts "\[ERROR\] Failed to spawn pipeline for $className: $pipe"
        cleanupWorkspace $className
        global exitCode
        set exitCode 1
        startNextJob
        return
    }
    
    fconfigure $pipe -blocking 0
    fileevent $pipe readable [list handleOutput $pipe $className]
    
    set active($pipe) $className
    incr activeCount
    
    # Try to start another job if workers are available
    startNextJob
}

puts "=========================================================="
puts "Starting Parallel Job Queue Supervisor"
puts "Queue: $queue"
puts "Max Concurrent Workers: $maxWorkers"
puts "=========================================================="

startNextJob

# Enter the Tcl native event loop
vwait forever
