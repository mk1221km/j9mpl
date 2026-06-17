# Zig Fixed-Size Ring Buffer Specification

## 1. Objective
Design a zero-dependency, fixed-capacity ring buffer data structure in Zig that stores f64 values. The buffer overwrites the oldest entry when full and exposes utilization metrics.

## 2. Invariants & Data Layout
* **Language:** Zig
* **Fixed Capacity:** 1024 entries
* **Behavior on Full:** Overwrite oldest entry (circular overwrite)
* **Allocator:** None required (fixed-size array on stack)

## 3. Data Structures
A public struct `RingBuffer` must expose:
* `buf` ([1024]f64) — internal fixed-size array
* `head` (usize) — index of next write position
* `count` (usize) — number of elements currently stored

## 4. Required Interfaces
The struct `RingBuffer` must implement:
1. `init() RingBuffer`: Returns a new zero-initialized ring buffer.
2. `push(value f64) void`: Appends a value. If full, overwrites oldest.
3. `readRange(dest: []f64, count: usize) []f64`: Copies up to N entries into the caller-owned dest buffer and returns the slice of copied elements.
4. `avg() f64`: Returns arithmetic mean. Returns 0 if empty.
5. `utilization() f64`: Ratio of stored to capacity (0.0 to 1.0).
