# Fixed-Size Ring Buffer Specification

## 1. Objective
Design a zero-dependency, fixed-capacity ring buffer data structure in Go that stores float64 values. The buffer overwrites the oldest entry when full, supports atomic push and read-range operations, and exposes utilization metrics.

## 2. Invariants & Data Layout
* **Storage Engine:** In-memory contiguous array (no database)
* **Fixed Capacity:** 1024 entries
* **Behavior on Full:** Overwrite oldest entry (circular overwrite)
* **Thread Safety:** None required (single-threaded)

## 3. Data Structures
A public struct `RingBuffer` must expose:
* `buf` ([1024]float64) — internal fixed-size array
* `head` (int) — index of next write position
* `count` (int) — number of elements currently stored

## 4. Required Interfaces
The struct `RingBuffer` must implement:
1. `NewRingBuffer() *RingBuffer`: Returns a new zero-initialized ring buffer with capacity 1024.
2. `Push(value float64)`: Appends a value. If the buffer is full, overwrites the oldest entry.
3. `ReadRange(count int) []float64`: Returns the last N entries without removing them. If count exceeds stored elements, return all available.
4. `Avg() float64`: Returns the arithmetic mean of all stored values. Returns 0 if empty.
5. `Utilization() float64`: Returns the ratio of stored elements to capacity (0.0 to 1.0).
