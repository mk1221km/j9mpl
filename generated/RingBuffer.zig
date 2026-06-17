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
    const n = @min(count, self.Count, dest.len);
    if (n == 0) return dest[0..0];

    const start: usize = if (self.Count == 1024) self.Head else 0;
    const buf_len: usize = 1024;

    if (start + n <= buf_len) {
        @memcpy(dest[0..n], self.Buf[start..start+n]);
    } else {
        const first_len = buf_len - start;
        @memcpy(dest[0..first_len], self.Buf[start..buf_len]);
        @memcpy(dest[first_len..n], self.Buf[0..(n - first_len)]);
    }

    return dest[0..n];
}
	pub fn avg(self: *Self) f64 {
    if (self.Count == 0) return 0.0;
    var sum: f64 = 0.0;
    var i: usize = 0;
    while (i < self.Count) : (i += 1) {
        sum += self.Buf[i];
    }
    return sum / @as(f64, @floatFromInt(self.Count));
}
	pub fn utilization(self: *Self) f64 {
    return @as(f64, self.Count) / 1024.0;
}
};
