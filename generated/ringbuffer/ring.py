#!/usr/bin/env python3
# Data-oriented ring buffer in Python.
# No structs, no classes, no type system nonsense.
# Flat list, scalar head/count, freestanding functions.

import sys

CAPACITY = 1024

def ring_push(buf, head, count, value):
    buf[head] = value
    head = (head + 1) % CAPACITY
    if count < CAPACITY:
        count += 1
    return head, count

def ring_avg(buf, count):
    if count == 0:
        return 0.0
    total = 0.0
    for i in range(count):
        total += buf[i]
    return total / count

def ring_read_range(buf, head, count, n):
    if n <= 0 or count == 0:
        return []
    if n > count:
        n = count
    return [buf[(head - n + i) % CAPACITY] for i in range(n)]

def ring_utilization(count):
    return count / CAPACITY

def main():
    storage = [0.0] * CAPACITY
    head = 0
    count = 0

    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        parts = line.split()
        if not parts:
            continue

        cmd = parts[0]

        if cmd == "push":
            if len(parts) < 2:
                continue
            value = float(parts[1])
            head, count = ring_push(storage, head, count, value)
        elif cmd == "avg":
            print(f"{ring_avg(storage, count):.1f}")
        elif cmd == "readRange":
            n = int(parts[1]) if len(parts) > 1 else 1
            for v in ring_read_range(storage, head, count, n):
                print(f"{v:.1f}")
        elif cmd == "utilization":
            print(f"{ring_utilization(count):.4f}")

if __name__ == "__main__":
    main()
