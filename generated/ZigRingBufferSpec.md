# Ring Buffer — Data-Oriented Primitive Contract

## 1. Objective
Fixed-capacity ring buffer storing f64 values. Overwrites oldest entry when full. Same contract across all languages (Go, LuaJIT, Zig, etc.). Verified by runtime tests, not compile-time type checks.

## 2. Data Layout (language-invariant)
State is three independent values — no wrapping struct or object:
* `buf` — contiguous array of f64, capacity 1024
* `head` — integer index of next write position
* `count` — integer number of elements currently stored (0..1024)

## 3. Operations (freestanding functions, not methods)
Each operation takes flat primitives and returns primitives. No struct receiver, no `self`, no `this`.

1. **push(buf, head, count, value)** → new_head, new_count
   Write value at head position. Advance head (circular). Increment count up to capacity.

2. **avg(buf, count)** → f64
   Compute arithmetic mean of stored elements. Return 0 if count is 0. Summation order does not matter — sum is commutative.

3. **readRange(buf, head, count, n)** → output_values
   Return the n most recent stored elements, oldest-first within the returned set. If n exceeds count, return all available elements. If count is 0 or n is 0, return nothing.

   Index computation (0-indexed proof):
   `idx = (head - n + i + capacity) % capacity` for i in 0..n-1

4. **utilization(count)** → f64
   Return count / capacity as a float in [0.0, 1.0].

## 4. Runtime Verification (shift-right)
All verification happens at runtime — no compile-time type safety assumed.
* Behavioral test suite drives the compiled binary via stdin/stdout
* Test suite is language-agnostic (Tcl), same tests for any implementation
* Edge cases verified: empty buffer, single element, full buffer, wraparound after 1024+ writes, readRange bounds, zero/negative/large values

## 5. Stdin/stdout Wire Protocol
```
push <value>       - store value, no output
avg                - print one line: "%.1f"
readRange [n]      - print n lines "%.1f" (default n=1)
utilization        - print one line: "%.4f"
```
Each command is one line. Each response is one value per line. No prompts, no headers.
