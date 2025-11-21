#!/bin/bash
set -e

# 颜色输出
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

DOCS_DIR="internal/adapters/http/docs"

# 1. 保存当前文档的 git hash
if [ -d "$DOCS_DIR" ]; then
    OLD_HASH=$(find "$DOCS_DIR" -type f \( -name "*.go" -o -name "*.json" -o -name "*.yaml" \) -exec md5sum {} \; 2>/dev/null | md5sum | cut -d' ' -f1)
else
    OLD_HASH=""
fi

# 2. 检查并安装 swag 工具
if ! command -v swag &> /dev/null; then
    echo -e "${YELLOW}⚠️  未找到 swag 工具。正在安装...${NC}"
    go install github.com/swaggo/swag/cmd/swag@latest
    echo -e "${GREEN}✅ swag 工具已安装${NC}"
fi

# 3. 生成 Swagger 文档
echo -e "${BLUE}📝 正在生成 Swagger 文档...${NC}"
swag init \
    -g internal/commands/api/api.go \
    -o internal/adapters/http/docs \
    --parseDependency \
    --parseInternal > /dev/null 2>&1 || {
    echo -e "${RED}❌ Swagger 文档生成失败${NC}"
    exit 1
}

# 4. 计算新的文档 hash
if [ -d "$DOCS_DIR" ]; then
    NEW_HASH=$(find "$DOCS_DIR" -type f \( -name "*.go" -o -name "*.json" -o -name "*.yaml" \) -exec md5sum {} \; 2>/dev/null | md5sum | cut -d' ' -f1)
else
    NEW_HASH=""
fi

# 5. 检查文档是否有变化
if [ "$OLD_HASH" != "$NEW_HASH" ]; then
    echo -e "${YELLOW}⚠️  Swagger 文档已更新，正在自动添加到暂存区...${NC}"
    git add "$DOCS_DIR"/*.go "$DOCS_DIR"/*.json "$DOCS_DIR"/*.yaml 2>/dev/null || true
    echo -e "${GREEN}✅ Swagger 文档已添加到此次提交${NC}"
else
    echo -e "${GREEN}✅ Swagger 文档无变化${NC}"
fi

exit 0
