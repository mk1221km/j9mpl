const std = @import("std");

pub const RingBuffer = struct {
	Buf: [1024]f64,
	Head: usize,
	Count: usize,

	const Self = @This();

	pub fn init() RingBuffer {
    return RingBuffer{};
}
	pub fn push(self: *Self, value: f64) void {
    self.Buf[self.Head] = value;
    self.Head = (self.Head + 1) % 1024;
    if (self.Count < 1024) {
        self.Count += 1;
    }
}
	pub fn readRange(self: *Self, dest: []f64, count: usize) []f64 {
    const n = @min(count, self.Count);
    if (n == 0) return dest[0..0];
    const start = if (self.Head >= n) self.Head - n else self.Head + 1024 - n;
    var i: usize = 0;
    while (i < n) : (i += 1) {
        dest[i] = self.Buf[(start + i) % 1024];
    }
    return dest[0..n];
}
	method avg() public static f64
if (self.Count == 0) return 0.0;
var sum: f64 = 0.0;
var i: usize = 0;
while (i < self.Count) : (i += 1) {
    sum += self.Buf[i];
}
return sum / @as(f64, @floatFromInt(self.Count));
	// SKELETON_utilization
};
