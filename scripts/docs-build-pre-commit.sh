#!/usr/bin/env bash

# Pre-commit hook entry for building VitePress docs when documentation files change.
# This script can be invoked directly or via the pre-commit framework.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "Running docs build pre-commit hook..."
echo "--------------------------------------"

# Prefer incoming file list (pre-commit), otherwise fall back to staged diff.
if [ "$#" -gt 0 ]; then
    DOCS_CHANGED_FILES="$*"
else
    DOCS_CHANGED_FILES="$(git diff --cached --name-only -- 'docs/**' 'docs/*' || true)"
fi

if [ -z "$DOCS_CHANGED_FILES" ]; then
    echo "No staged docs changes detected; skipping docs build."
    exit 0
fi

if [ ! -d "$ROOT_DIR/docs" ]; then
    echo "docs/ directory not found; skipping docs build."
    exit 0
fi

if [ ! -f "$ROOT_DIR/docs/package.json" ]; then
    echo "docs/package.json not found; skipping docs build."
    exit 0
fi

if [ ! -d "$ROOT_DIR/docs/node_modules" ]; then
    echo "Installing docs dependencies (npm install)..."
    (cd "$ROOT_DIR/docs" && npm install)
fi

echo "Building docs..."
if (cd "$ROOT_DIR/docs" && npm run build); then
    echo "Docs build succeeded."
    exit 0
else
    BUILD_EXIT_CODE=$?
    echo "Docs build failed with exit code $BUILD_EXIT_CODE."
    exit "$BUILD_EXIT_CODE"
fi
