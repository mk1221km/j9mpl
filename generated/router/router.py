#!/usr/bin/env python3
# Stateless transaction router — pure compute over wire protocol.
# Tcl owns persistence; this binary only classifies and counts.
# Protocol: space-separated fields.
#   route <txId> <type> <amount>  →  "routed <txId> <channel>"
#   rate                           →  "0.0"
#   count <status>                 →  integer count

import sys

def route(tx_type, amount):
    """Determine channel based on transaction type and amount."""
    # Known high-priority types escalate; everything else is standard
    is_high = tx_type.lower() == "high"
    if amount <= 1000.0:
        return "HIGH_PRIORITY_WIRE" if is_high else "CLEAR_ACH"
    else:
        return "HIGH_PRIORITY_WIRE"

def main():
    tx_count = 0

    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        parts = line.split()
        if not parts:
            continue

        cmd = parts[0]

        if cmd == "route":
            if len(parts) < 4:
                print("INVALID")
                continue
            tx_id = parts[1]
            tx_type = parts[2]
            try:
                amount = float(parts[3])
            except ValueError:
                print("INVALID")
                continue

            channel = route(tx_type, amount)
            tx_count += 1
            print(f"routed {tx_id} {channel}")

        elif cmd == "rate":
            print("0.0")

        elif cmd == "count":
            status = parts[1] if len(parts) > 1 else ""
            if status == "routed":
                print(tx_count)
            else:
                print("0")

if __name__ == "__main__":
    main()
