# Autonomous Self-Correction Loop: Architectural Blueprint & Implementation

This document presents the design and implementation of the **Autonomous Self-Correction Loop** (Phase I Deployment). The goal of this loop is to automatically ingest compilation and translation diagnostics, isolate the errant NetRexx block via database coordinate lookup or text block boundary scanning, dispatch a target repair prompt to the remote LLM, and apply the patch back to the source code without human intervention.

---

## 1. The Error Ingestion & Repair Protocol

The self-correction pipeline handles two distinct failure phases:

1. **Phase A: NetRexx Translation Faults (`nrc` stage)**
   - When a syntax or scope error (such as an undeclared variable) is encountered during NetRexx-to-Java translation, `nrc` reports the error line number of the `.nrx` file directly.
   - The Go engine intercepts this stream, strips ANSI color escapes, and scans the `.nrx` source file line-by-line to locate the enclosing class or method block.
   - The exact source block is extracted and packaged into a targeted prompt.

2. **Phase B: Standalone Java Compilation Faults (`ecj` stage)**
   - When translation succeeds but the ECJ compiler fails due to type checks or symbol resolution issues, ECJ outputs the error line number in the generated `.java` file.
   - The Go engine intercepts ECJ diagnostics, queries the **SQLite context ledger** (`project_context.db`) to map the Java line coordinate to a logical symbol URI (e.g. `java+method://...`), resolves that URI to its NetRexx source path, and extracts the corresponding block.

```
                      +---------------------------------------+
                      |       Input NetRexx Source Code       |
                      +---------------------------------------+
                                          |
                                          v
                      +---------------------------------------+
                      |    1. Run NRC Translation Utility     |
                      +---------------------------------------+
                                          |
                      +-------------------+-------------------+
                      | (nrc fails)                           | (nrc succeeds)
                      v                                       v
        +----------------------------+          +-----------------------------+
        | Parse Translator Diagnostic|          | 2. Run Headless ECJ Compiler|
        |  (Strip ANSI escapes &     |          +-----------------------------+
        |   get nrx line coordinate) |                        |
        +----------------------------+       +----------------+----------------+
                      |                      | (ecj succeeds)                  | (ecj fails)
                      |                      v                                 v
                      |          +------------------------+      +-----------------------------+
                      |          |  SUCCESS: Deploy Class |      | Parse ECJ Diagnostics       |
                      |          |  files to OpenJ9 JRE   |      |  (Get java line coordinate) |
                      |          +------------------------+      +-----------------------------+
                      |                                                        |
                      |                                                        v
                      |                                          +-----------------------------+
                      |                                          | Query SQLite context ledger |
                      |                                          |  to resolve logical symbol  |
                      |                                          +-----------------------------+
                      |                                                        |
                      v                                                        v
        +------------------------------------------------------------------------------+
        | 3. Extract target NetRexx block context (lines/method)                       |
        +------------------------------------------------------------------------------+
                                          |
                                          v
        +------------------------------------------------------------------------------+
        | 4. Package prompt (.context/self_correct_prompt.txt) & original block        |
        +------------------------------------------------------------------------------+
                                          |
                                          v
        +------------------------------------------------------------------------------+
        | 5. Dispatch prompt to remote model via 'agy --print'                         |
        +------------------------------------------------------------------------------+
                                          |
                                          v
        +------------------------------------------------------------------------------+
        | 6. Apply patch back to NetRexx source file using patch_source                |
        +------------------------------------------------------------------------------+
```

---

## 2. Go Orchestration Engine: `self_correct.go`

This binary sits at the center of the compilation boundary, executing translation/compilation, intercepting outputs, and packaging prompt payloads.

- **Source Location:** [bin/self_correct.go](file:///home/me/code/j9mpl/bin/self_correct.go)
- **Compiled Binary:** `bin/self_correct`

```go
// (Refer to bin/self_correct.go for full implementation details)
```

### Highlights:
- **ANSI Escape Sequence Stripping:** Cleans console outputs using `\x1b\[[0-9;]*[a-zA-Z]` to prevent parsing failures caused by colorized compiler streams.
- **Relational Context Mapping:** Queries the SQLite coordinate span database (`declarations` table) using `start_line` and `end_line` ranges to accurately resolve which logical class/method symbol contains the error.
- **Minimalist Footprint:** Runs natively in milliseconds without external heavyweight library dependencies.

---

## 3. Patching Utility: `patch_source.go`

This tool ensures deterministic updates to the source files by finding the exact original block text in the `.nrx` file and replacing it with the corrected block returned from the model.

- **Source Location:** [bin/patch_source.go](file:///home/me/code/j9mpl/bin/patch_source.go)
- **Compiled Binary:** `bin/patch_source`

```go
// (Refer to bin/patch_source.go for full implementation details)
```

---

## 4. Bash Loop Orchestrator: `self_correct_loop.sh`

This script ties the entire self-correction pipeline together, running the validation engine, feeding the prompt to the remote LLM using the `agy` CLI, patching the file, and retrying until compilation succeeds.

- **Location:** [bin/self_correct_loop.sh](file:///home/me/code/j9mpl/bin/self_correct_loop.sh)

```bash
#!/usr/bin/env bash
# Autonomous Self-Correction Loop for NetRexx / ECJ Toolchain
set -euo pipefail

if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <target_nrx_file>"
    exit 1
fi

NRX_FILE="${1}"
# Resolve absolute path
NRX_FILE="$(realpath "${NRX_FILE}")"
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CLASSPATH="lib/NetRexxF.jar:target/dependency/sqlite-jdbc-3.45.1.0.jar:target/dependency/slf4j-api-1.7.36.jar"

MAX_RETRIES=5
RETRY_COUNT=0

while [ "${RETRY_COUNT}" -lt "${MAX_RETRIES}" ]; do
    echo "=========================================================="
    echo "Iteration $((RETRY_COUNT+1)) / ${MAX_RETRIES}"
    echo "=========================================================="

    echo "[1/3] Executing compiler validation..."
    # Run self_correct. If it succeeds, exit code is 0.
    if "${PROJECT_DIR}/bin/self_correct" "${NRX_FILE}" -cp "${CLASSPATH}"; then
        echo "=========================================================="
        echo " [SUCCESS] Zero-error build achieved!"
        echo "=========================================================="
        exit 0
    fi

    echo "[2/3] Intercepting compiler diagnostic and generating prompt..."
    PROMPT_FILE="${PROJECT_DIR}/.context/self_correct_prompt.txt"
    ORIGINAL_BLOCK_FILE="${PROJECT_DIR}/.context/errant_block_original.txt"
    REVISED_BLOCK_FILE="${PROJECT_DIR}/.context/errant_block_revised.txt"
    TEMP_OUT="${PROJECT_DIR}/.context/model_output_raw.txt"

    echo "  -> Dispatching self-correction prompt to remote model..."
    # Run agy in non-interactive print mode with prompt content
    agy --print "$(cat "${PROMPT_FILE}")" > "${TEMP_OUT}"

    # Clean up model output to isolate the raw NetRexx block
    if grep -q '\`\`\`' "${TEMP_OUT}"; then
        echo "  -> Extracting block from code fences..."
        # Extract lines between first and last ``` block and strip the fences
        sed -n '/^```/,/^```/p' "${TEMP_OUT}" | grep -E -v '^```' > "${REVISED_BLOCK_FILE}"
    else
        echo "  -> Using raw model output..."
        cat "${TEMP_OUT}" > "${REVISED_BLOCK_FILE}"
    fi

    echo "[3/3] Applying patch to source file..."
    if ! "${PROJECT_DIR}/bin/patch_source" "${NRX_FILE}" "${ORIGINAL_BLOCK_FILE}" "${REVISED_BLOCK_FILE}"; then
        echo "[ERROR] Failed to patch source file. Self-correction loop aborted."
        exit 1
    fi

    echo "  -> Source file patched. Retrying compilation..."
    RETRY_COUNT=$((RETRY_COUNT+1))
done

echo "=========================================================="
echo " [ERROR] Failed to converge after ${MAX_RETRIES} attempts."
echo "=========================================================="
exit 1
```

---

## 5. Deployment Advantages

- **Model Context Efficiency:** Confining the model's focus strictly to the errant method/class block maximizes prefix cache hits on remote APIs, reducing latency and token costs.
- **Deterministic compilation constraints:** Combining NetRexx's clean language design with standalone ECJ's `-proceedOnError` validation ensures that syntax repairs converge predictably without introducing compilation bloat.
- **Zero-Dependency runtime footprint:** The orchestration loop uses native Go compiled binaries and standard bash commands, preserving the local-first, low-overhead characteristics of the factory stack.
