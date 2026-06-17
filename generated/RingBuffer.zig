const std = @import("std");

pub const RingBuffer = struct {
	Buf: [1024]f64,
	Head: usize,
	Count: usize,

	const Self = @This();

	pub fn init() RingBuffer {
    return RingBuffer{};
}
	method push(value f64) public static void
r.buf[r.head] = value;
r.head = (r.head + 1) % 1024;
if r.count < 1024 {
    r.count++;
}
	// SKELETON_readRange
	// SKELETON_avg
	// SKELETON_utilization
};
