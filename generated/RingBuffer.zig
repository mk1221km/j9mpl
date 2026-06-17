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
    const num_to_copy = @min(@min(count, self.Count), dest.len);
    const start = (self.Head + 1024 - self.Count) % 1024;
    for (0..num_to_copy) |i| {
        dest[i] = self.Buf[(start + i) % 1024];
    }
    return dest[0..num_to_copy];
}
	pub fn avg(self: *Self) f64 {
    if (self.Count == 0) return 0.0;
    var sum: f64 = 0.0;
    const oldest = (self.Head + 1024 - self.Count) % 1024;
    var i: usize = 0;
    while (i < self.Count) : (i += 1) {
        const idx = (oldest + i) % 1024;
        sum += self.Buf[idx];
    }
    return sum / @intToFloat(f64, self.Count);
}
	// SKELETON_utilization
};
