# Reasonix project memory

Notes the user pinned via the `#` prompt prefix. The whole file is
loaded into the immutable system prefix every session — keep it terse.

- Executive Memorandum: Primitive-First Topology vs. Generative Prior Constraints

The argument for a language-invariant primitive layout—reducing data structures to flat arrays, universal string slices, and scalar indices—is architecturally valid. Stripping out custom object wrappers and pointer-heavy indirection achieves optimal mechanical sympathy, unifies the data schema across disparate language backends, and ensures that the virtual memory footprint remains constant regardless of the runtime environment.

However, within an automated code-synthesis pipeline, the selection of the data layout paradigm cannot be decided solely by hardware alignment. It must balance **hardware efficiency** against **generative convergence velocity**. The current reliance on encapsulated structures with explicit method receivers or compact structs is a tactical concession to the token distribution constraints of the underlying large language model.

---

### 1. The Generative Bottleneck: Token Density vs. Ideological Purity

The primary execution engine of our software plant is a probabilistic token predictor. To achieve reliable, Turn-1 compilation convergence, the input specification must map onto patterns that occupy high-density regions of the model's training distribution.

* **The Encapsulated Baseline (High Density):** Modern open-source corpuses are overwhelmingly saturated with object-oriented or encapsulated type patterns (e.g., Go structs with receivers, Zig structs utilizing `@This()`). When the prompt compiler emits a skeleton bound to a struct namespace, it anchors the model's attention weights to highly standardized, well-represented code paradigms. The model can synthesize internal logic with an exceptionally low error rate because the boundary is rigidly constrained by the type system.
* **The Pure Primitive Baseline (Low Density):** Decomposing an application into stateless pipelines that pass raw, loose primitives (e.g., unboxed arrays alongside isolated cursor pointers) shifts the language model into an algorithmic reasoning mode. Instead of executing boilerplate slot-filling, the model must manually track index transformations and explicit pointer dereferencing (`buf`, `head`, and `count` passed as discrete arguments). This dramatically expands the semantic surface area for minor typographical and syntax errors, introducing significant non-deterministic noise into the local verification loop.

---

### 2. Comparative Analysis of State Topographies

To illustrate the structural implications of a pure primitive-first architecture across our active targets, consider the execution signatures required for a stateless, unboxed `Push` operation:

#### Go Primitive Transformation

```go
package ringbuffer

// Stateless primitive mutation forces explicit value tracking and return tuple unpacking
func Push(buf []float64, head int, count int, value float64) (int, int) {
	buf[head] = value
	nextHead := (head + 1) % len(buf)
	if count < len(buf) {
		count++
	}
	return nextHead, count
}

```

#### Zig Primitive Transformation

```zig
// Pure-memory slice mutation requires explicit pointer dereferencing on scalar indices
pub fn push(buf: []f64, head: *usize, count: *usize, value: f64) void {
    buf[head.*] = value;
    head.* = (head.* + 1) % buf.len;
    if (count.* < buf.len) {
        count.* += 1;
    }
}

```

While the Zig primitive approach strips away the encapsulated object layer, it introduces explicit pointer dereferencing (`head.*`, `count.*`) to mutate state across the function boundary. Our execution telemetry confirms that the model regularly stumbles on Zig’s strict reference mechanics, frequently failing compile-time checks due to missing or misplaced dereference tokens. The encapsulated struct method hides this language-specific pointer syntax behind an implicit or clean explicit receiver (`self.Head`), insulating the generator from syntax traps.

---

### 3. Synthesis Pipeline Impact Matrix

Flipping the entire software factory to a language-invariant primitive paradigm introduces structural trade-offs across the engineering track:

| Operational Axis | Encapsulated Struct Target | Language-Invariant Primitive Target |
| --- | --- | --- |
| **Model Ingestion Error Rate** | **Low.** Relies on standard object/method templates. | **High.** Prone to index tracking and reference syntax errors. |
| **Parser Complexity** | **Moderate.** Parses high-level Markdown blocks into skeletons. | **High.** Requires parsing complex procedural argument tables. |
| **Cross-Language Symmetry** | Low. Signatures vary by language design patterns. | **Absolute. Identical raw primitives across all backends.** |
| **Local Verification Latency** | Minor overhead from type-checking constraints. | **Zero abstraction overhead.** Direct compilation to raw iron. |

---

### Strategic Conclusion

The choice is clear: encapsulated signatures are used because they act as a stabilization mechanism for automated generation, not because they are computationally superior. They turn open-source coding conventions into a shield that protects our pipeline from compiler rejections.

Transitioning the factory to a pure primitive, language-invariant architecture is entirely possible, but it requires a significantly more robust, version-locked few-shot exemplar ledger to firmly hold the model's attention weights to stateless array mathematics.

Given that our current Zig skeleton has achieved stability and is converting scalar routines efficiently, should we finalize the full 5/5 method convergence for the current `RingBuffer` contract to establish our execution baseline, or halt the current run to restructure the specification grammar inside `bin/spec_parser.go` for a pure primitive, function-only pipeline?
