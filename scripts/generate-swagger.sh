#!/bin/bash
set -e

# 颜色输出
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🔄 正在生成 Swagger/OpenAPI 文档${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 检查 swag 是否已安装
if ! command -v swag &> /dev/null; then
    echo -e "${YELLOW}⚠️  未找到 swag 工具。正在安装...${NC}"
    go install github.com/swaggo/swag/cmd/swag@latest
    echo -e "${GREEN}✅ swag 工具已安装${NC}"
    echo ""
fi

# 生成 Swagger 文档
echo -e "${BLUE}📝 正在生成 API 文档...${NC}"
swag init \
    -g internal/commands/api/api.go \
    -o internal/adapters/http/docs \
    --parseDependency \
    --parseInternal

echo ""
echo -e "${GREEN}✅ Swagger 文档生成成功!${NC}"
echo ""
echo -e "${BLUE}📁 生成的文件:${NC}"
ls -lh internal/adapters/http/docs/

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🚀 查看文档的方式:${NC}"
echo -e "   1. 启动 API 服务器: ${YELLOW}go run main.go api${NC}"
echo -e "   2. 在浏览器中打开: ${YELLOW}http://localhost:40012/swagger/index.html${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
