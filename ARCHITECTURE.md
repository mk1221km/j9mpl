# Architecture

## Core Invariants

Three design rules apply to every module in this plant. Violating any of them reintroduces the complexity we spent this session eliminating.

### 1. Silent Writes

Write operations produce no output:

| Module | Write command | Output |
|---|---|---|
| RingBuffer | `push <value>` | *(silent)* |
| MetricsLogger | `log <name> <value>` | *(silent)* |
| TransactionRouter | `route <id> <type> <amount>` | *(silent, no — see rule 3)* |

This keeps output volume bounded by query volume, not input volume. A 10-million-line stream of `push` commands produces zero output and cannot deadlock a pipe buffer.

**Exception:** `route` produces a single-line response because routing is a classification, not a storage operation. The response is one line per command, not one per event — still bounded.

### 2. Test Suite Is the Specification

There are no spec documents that define correct behavior. The Tcl test scripts are the sole authority:

- `bin/test_harness.tcl` — RingBuffer (25 base + 5 stress)
- `bin/test_metrics.tcl` — MetricsLogger (15 tests)
- `bin/test_router.tcl` — TransactionRouter (15 tests)
- `bin/chaos_injector.tcl` — cross-module anomaly injection (12 tests)

A binary passes if it produces the expected output for the given input. Nothing else matters — not types, not documentation, not compiler warnings.

### 3. Verification by Query

Chaos resistance is validated by querying after corruption, not by parsing acknowledgment tokens. A binary proves it survived corrupt input by answering a subsequent `avg` or `rate` correctly. There is no `INVALID` token, no survival flag, no heartbeat.

---

## Protocol

Every binary reads space-separated commands from stdin and writes results to stdout. One command per line. One result per line. No prompts, no headers, no trailers.

### RingBuffer

| Command | Output |
|---|---|
| `push <value>` | *(silent)* |
| `avg` | `N.N` |
| `readRange [n]` | `N.N` per element (default n=1) |
| `utilization` | `N.NNNN` |

### MetricsLogger

| Command | Output |
|---|---|
| `log <name> <value>` | *(silent)* |
| `avg <name>` | `N.N` |
| `count <name>` | `N` |

### TransactionRouter

| Command | Output |
|---|---|
| `route <id> <type> <amount>` | `routed <id> <channel>` |
| `rate` | `N.N` |
| `count <status>` | `N` |

Routing rules (built-in):
- amount ≤ 1000, type `high` → `HIGH_PRIORITY_WIRE`
- amount ≤ 1000, other → `CLEAR_ACH`
- amount > 1000 → `HIGH_PRIORITY_WIRE`

---

## Module Inventory

| Module | Tests | Python | Go | LuaJIT | Zig | OCaml | Bun/JS |
|---|---|---|---|---|---|---|---|
| RingBuffer | 25+5 | ring.py | ring.go | ring.lua | RingBuffer.zig | ring.ml | ring.js |
| MetricsLogger | 15 | metrics.py | — | — | — | — | — |
| TransactionRouter | 15 | router.py | — | — | — | — | — |

The ring buffer has full multi-language coverage because it was the proof-of-concept. Metrics and router only have Python, which is sufficient — adding Go or Zig is an exercise in copying the same logic with different syntax.

---

## Transport Agnosticism

The binary reads stdin and writes stdout. That is the only contract. Tcl manages the transport layer — the binary never knows whether it was invoked via pipe, Unix socket, HTTP, or any other mechanism.

To switch transports, change the Tcl adapter, not the binary:

```tcl
proc pipe_transport {binary input} {
    set fd [open "|$binary" r+]
    puts $fd $input
    close $fd w
    return [string trim [read $fd]]
}

proc socket_transport {fd input} {
    puts $fd $input
    flush $fd
    gets $fd line
    return [string trim $line]
}
```

The same `metrics.py` or `ring.py` serves either transport without modification. Transport framing (HTTP headers, socket handshakes, checksums) is stripped by the Tcl adapter before the data reaches the binary. The binary sees only a clean line of text.

See `bin/transport_demo.tcl` for a working example.

## How to Add a Language

1. Write a program that reads stdin, writes stdout, and implements the protocol for a module.
2. Run the corresponding Tcl test suite against it.
3. If all tests pass, the implementation is correct.

The Tcl controller handles persistence (SQLite), orchestration, and verification. The binary handles nothing but text transformation. There is no build chain beyond what the language needs to produce an executable.

---

## Operational Summary

| Assertion | Status |
|---|---|
| Write operations produce no output | Enforced across all modules |
| Pipe protocol handles any input volume | Verified at 10,000 lines |
| Corrupt input does not corrupt state | Verified (NaN, Inf, missing fields, long names) |
| Floating-point epsilon boundaries | Verified at 999.999 / 1000.0 / 1000.001 |
| Same test suite works for all languages | Verified across 6 runtimes |
| No database drivers in binaries | Verified — all persistence in Tcl |
