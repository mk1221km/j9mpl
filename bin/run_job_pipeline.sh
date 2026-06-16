#!/usr/bin/env bash
# Runs the full synthesis & verification pipeline for a single spec inside an isolated workspace
set -euo pipefail

SPEC_FILE="${1}" # e.g. generated/TransactionRouterSpec.md
CLASS_NAME="$(basename "${SPEC_FILE}" Spec.md)"

echo "=== Starting job pipeline for ${CLASS_NAME} ==="

# 1. Parse Spec
echo "[1/5] Parsing specification..."
./bin/spec_parser "${SPEC_FILE}"

# 2. Incremental synthesis prep
echo "[2/5] Incremental Synthesis mode active. Class skeleton generated."

# 3. Compile/self-correct main class
echo "[3/5] Cleaning stale generation artifacts..."
rm -f "generated/${CLASS_NAME}.java"

# Compile helper record classes first if they exist
for rec in generated/*Record.nrx; do
  if [ -f "$rec" ]; then
    echo "[3/5] Compiling helper record class ${rec}..."
    ./bin/self_correct_loop.tcl "$rec"
  fi
done

echo "[3/5] Compiling and self-correcting main class..."
./bin/self_correct_loop.tcl "generated/${CLASS_NAME}.nrx"

# 4. Generate tests
echo "[4/5] Generating boundary fuzzer tests..."
./bin/run_test_generator.sh "${CLASS_NAME}"

# 5. Compile/self-correct and run tests inside sandbox
echo "[5/5] Compiling and executing fuzzer tests..."
./bin/self_correct_loop.tcl "generated/${CLASS_NAME}Test.nrx"

echo "=== Pipeline completed successfully for ${CLASS_NAME} ==="
