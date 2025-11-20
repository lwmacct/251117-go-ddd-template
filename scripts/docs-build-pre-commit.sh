#!/usr/bin/env bash

# 用于在文档文件变更时构建 VitePress 文档的 pre-commit 钩子入口。
# 此脚本可以直接调用或通过 pre-commit 框架调用。

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "正在运行文档构建 pre-commit 钩子..."
echo "--------------------------------------"

# 优先使用传入的文件列表 (pre-commit),否则回退到暂存区差异检查。
if [ "$#" -gt 0 ]; then
    DOCS_CHANGED_FILES="$*"
else
    DOCS_CHANGED_FILES="$(git diff --cached --name-only -- 'docs/**' 'docs/*' || true)"
fi

if [ -z "$DOCS_CHANGED_FILES" ]; then
    echo "未检测到暂存的文档变更；跳过文档构建。"
    exit 0
fi

if [ ! -d "$ROOT_DIR/docs" ]; then
    echo "未找到 docs/ 目录；跳过文档构建。"
    exit 0
fi

if [ ! -f "$ROOT_DIR/docs/package.json" ]; then
    echo "未找到 docs/package.json；跳过文档构建。"
    exit 0
fi

if [ ! -d "$ROOT_DIR/docs/node_modules" ]; then
    echo "正在安装文档依赖 (npm install)..."
    (cd "$ROOT_DIR/docs" && npm install)
fi

echo "正在构建文档..."
echo "提示: 首次构建或缓存失效时可能需要 30-60 秒，请耐心等待..."

# 禁用 npm 进度条和交互式输出，避免在 pre-commit 环境中卡住
export CI=true
export npm_config_progress=false
export npm_config_loglevel=error

if (cd "$ROOT_DIR/docs" && npm run build 2>&1); then
    echo "✅ 文档构建成功。"
    exit 0
else
    BUILD_EXIT_CODE=$?
    echo "❌ 文档构建失败，退出码 $BUILD_EXIT_CODE。"
    exit "$BUILD_EXIT_CODE"
fi
