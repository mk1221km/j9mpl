#!/usr/bin/env bash
# Property-Based Testing Verification Loop generator (Phase III)
set -euo pipefail

if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <target_class_name>"
    exit 1
fi

CLASS_NAME="${1}"
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
DB_PATH="${PROJECT_DIR}/.context/project_context.db"
DECLS_CSV="${PROJECT_DIR}/.context/declarations.csv"
OUTPUT_NRX="${PROJECT_DIR}/generated/${CLASS_NAME}Test.nrx"

echo "=========================================================="
echo "Phase III: Property-Based Verification Loop Generator"
echo "=========================================================="
echo "Target Class: ${CLASS_NAME}"
echo "Ledger Path:  ${DB_PATH}"
echo "=========================================================="

# 1. Export type & declaration facts from SQLite ledger to CSV
echo "[1/3] Exporting declarations from context ledger..."
sqlite3 -csv "${DB_PATH}" "SELECT symbol_uri, file_path, start_line, end_line FROM declarations;" > "${DECLS_CSV}"

# 2. Invoke Rascal generator
echo "[2/3] Executing Rascal TestGenerator..."
java -jar "${PROJECT_DIR}/rascal-shell-stable.jar" TestGenerator "${CLASS_NAME}" "${DECLS_CSV}" "${OUTPUT_NRX}"

# Boundaries are now natively generated and injected by the Rascal generator from .context/fuzzer_boundaries.json

# 3. Clean up
echo "[3/3] Cleaning up temporary files..."
rm -f "${DECLS_CSV}"

echo "=========================================================="
echo "Harness Generation Complete: ${OUTPUT_NRX}"
echo "=========================================================="
