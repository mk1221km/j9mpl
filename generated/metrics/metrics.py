#!/usr/bin/env python3
# Stateless metrics logger — pure compute over wire protocol.
# Tcl owns persistence; this binary only validates and aggregates.
# Protocol: space-separated fields.
#   log <name> <value> [timestamp]  →  (no output, like push in ring buffer)
#   avg <name>                      →  numeric value
#   count <name>                    →  integer count
#
# No per-line acknowledgments. Output only when queried.
# This prevents pipe buffer deadlock at high input volumes.

import sys
import math

def main():
    store = {}  # name -> list of (value, timestamp)

    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        parts = line.split()
        if not parts:
            continue

        cmd = parts[0]

        if cmd == "log":
            if len(parts) < 3:
                continue
            name = parts[1]
            try:
                value = float(parts[2])
            except ValueError:
                continue

            # Silently reject non-finite IEEE 754 values
            if math.isnan(value) or math.isinf(value):
                continue

            timestamp = parts[3] if len(parts) > 3 else "0"

            if name not in store:
                store[name] = []
            store[name].append((value, timestamp))

        elif cmd == "avg":
            name = parts[1] if len(parts) > 1 else ""
            values = [v for v, _ in store.get(name, [])]
            if not values:
                print("0.0")
            else:
                avg = sum(values) / len(values)
                print(f"{avg:.1f}")

        elif cmd == "count":
            name = parts[1] if len(parts) > 1 else ""
            values = store.get(name, [])
            print(len(values))

if __name__ == "__main__":
    main()
