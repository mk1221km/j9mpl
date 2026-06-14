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
