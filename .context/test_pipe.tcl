set fd [open "|python3 -u generated/metrics/metrics.py" r+]
puts $fd "log test 42.0"
flush $fd
close $fd w
set output [read $fd]
close $fd
puts "output: [string trim $output]"
