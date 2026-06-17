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
	method readRange(count usize) public static []f64

	// SKELETON_avg
	// SKELETON_utilization
};
