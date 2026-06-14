#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TARGET_CLASS="${1}"
shift 1

# Dynamically locate active Java installation path and Java Home
JAVA_EXEC="$(readlink -f "$(which java)")"
JAVA_HOME_DIR="$(dirname "$(dirname "${JAVA_EXEC}")")"

# System binaries and runtime library mappings
BWRAP_ARGS=(
    --ro-bind /usr /usr
    --ro-bind /lib /lib
    --ro-bind /lib64 /lib64
    --ro-bind /etc/alternatives /etc/alternatives
    --ro-bind /etc/ssl /etc/ssl
    --ro-bind "${JAVA_HOME_DIR}" "${JAVA_HOME_DIR}"
    --proc /proc
    --dev /dev
    --tmpfs /tmp
)

# Workspace mapping - isolate dependencies while locking source modifications
BWRAP_ARGS+=(
    --ro-bind "${PROJECT_ROOT}/bin" /app/bin
    --ro-bind "${PROJECT_ROOT}/lib" /app/lib
    --ro-bind "${PROJECT_ROOT}/target/dependency" /app/target/dependency
    --dir /app/generated
    --tmpfs /app/generated
    --chdir /app
)

# Unshare all kernel namespaces to complete isolation profiles
BWRAP_ARGS+=(
    --unshare-all
    --hostname factory-sandbox
)

exec bwrap "${BWRAP_ARGS[@]}" \
    "${JAVA_EXEC}" -cp "bin:lib/NetRexxF.jar:target/dependency/*" "${TARGET_CLASS}" "$@"
