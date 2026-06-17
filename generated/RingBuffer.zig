const std = @import("std");

pub const RingBuffer = struct {
	buf: [1024]f64,
	head: usize,
	count: usize,

	const Self = @This();

	// SKELETON_init
	// SKELETON_push
	// SKELETON_readRange
	// SKELETON_avg
	// SKELETON_utilization
};
