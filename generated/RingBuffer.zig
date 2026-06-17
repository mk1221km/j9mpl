// Data-oriented ring buffer — flat array, scalars, freestanding functions.
// No struct, no methods, no @This(), no pointer receivers.

const std = @import("std");

const capacity = 1024;

// ringPush writes value at head, returns new head and new count.
fn ringPush(buf: []f64, head: usize, count: usize, value: f64) struct { head: usize, count: usize } {
    buf[head] = value;
    return .{
        .head = (head + 1) % capacity,
        .count = if (count < capacity) count + 1 else count,
    };
}

// ringAvg returns arithmetic mean of stored elements. Returns 0 if empty.
fn ringAvg(buf: []f64, count: usize) f64 {
    if (count == 0) return 0.0;
    var sum: f64 = 0.0;
    var i: usize = 0;
    while (i < count) : (i += 1) {
        sum += buf[i];
    }
    return sum / @as(f64, @floatFromInt(count));
}

// ringReadRange copies up to n of the most recent elements into dest.
// Returns the number of elements written.
fn ringReadRange(buf: []f64, head: usize, count: usize, n: usize, dest: []f64) usize {
    const actual = @min(n, count, dest.len);
    if (actual == 0) return 0;
    var i: usize = 0;
    while (i < actual) : (i += 1) {
        const idx = (head + capacity - actual + i) % capacity;
        dest[i] = buf[idx];
    }
    return actual;
}

// ringUtilization returns ratio of stored count to capacity.
fn ringUtilization(count: usize) f64 {
    return @as(f64, @floatFromInt(count)) / @as(f64, @floatFromInt(capacity));
}

pub fn main() !void {
    var storage: [capacity]f64 = undefined;
    var head: usize = 0;
    var count: usize = 0;
    var readbuf: [1024 * 128]u8 = undefined;
    var total: usize = 0;

    // Read all input from stdin until EOF
    while (true) {
        const n = std.os.linux.read(0, @as([*]u8, @ptrCast(&readbuf[total])), readbuf.len - total);
        if (n == 0 or n > readbuf.len - total) break;
        total += n;
        if (total >= readbuf.len) break;
    }
    if (total == 0) return;
    const data = readbuf[0..total];

    // Process lines
    var line_start: usize = 0;
    while (line_start < total) {
        // Find end of line
        const line_end = std.mem.indexOfScalarPos(u8, data, line_start, '\n') orelse total;
        const line = std.mem.trim(u8, data[line_start..line_end], " \t\r");
        line_start = line_end + 1;

        if (line.len == 0) continue;

        var it = std.mem.tokenizeAny(u8, line, " \t");
        const cmd = it.next() orelse continue;

        if (std.mem.eql(u8, cmd, "push")) {
            const valueStr = it.next() orelse continue;
            const value = std.fmt.parseFloat(f64, valueStr) catch continue;
            const result = ringPush(storage[0..], head, count, value);
            head = result.head;
            count = result.count;
        } else if (std.mem.eql(u8, cmd, "avg")) {
            const avg_val = ringAvg(storage[0..], count);
            var outbuf: [64]u8 = undefined;
            const out = try std.fmt.bufPrint(&outbuf, "{d:.1}\n", .{avg_val});
            _ = std.os.linux.write(1, out.ptr, out.len);
        } else if (std.mem.eql(u8, cmd, "readRange")) {
            const nStr = it.next() orelse "1";
            const n = std.fmt.parseInt(usize, nStr, 10) catch 1;
            var dest: [capacity]f64 = undefined;
            const written = ringReadRange(storage[0..], head, count, n, dest[0..]);
            var i: usize = 0;
            while (i < written) : (i += 1) {
                var outbuf: [64]u8 = undefined;
                const out = try std.fmt.bufPrint(&outbuf, "{d:.1}\n", .{dest[i]});
                _ = std.os.linux.write(1, out.ptr, out.len);
            }
        } else if (std.mem.eql(u8, cmd, "utilization")) {
            const util = ringUtilization(count);
            var outbuf: [64]u8 = undefined;
            const out = try std.fmt.bufPrint(&outbuf, "{d:.4}\n", .{util});
            _ = std.os.linux.write(1, out.ptr, out.len);
        }
    }
}
