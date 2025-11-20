#!/bin/bash
set -e

# 颜色输出
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🔄 Generating Swagger/OpenAPI Documentation${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 检查 swag 是否已安装
if ! command -v swag &> /dev/null; then
    echo -e "${YELLOW}⚠️  swag tool not found. Installing...${NC}"
    go install github.com/swaggo/swag/cmd/swag@latest
    echo -e "${GREEN}✅ swag tool installed${NC}"
    echo ""
fi

# 生成 Swagger 文档
echo -e "${BLUE}📝 Generating API documentation...${NC}"
swag init \
    -g internal/commands/api/api.go \
    -o internal/adapters/http/docs \
    --parseDependency \
    --parseInternal \
    --parseDepth 2 \
    --exclude internal/infrastructure/database,internal/infrastructure/queue,internal/infrastructure/redis,internal/bootstrap,internal/commands/migrate,internal/commands/seed,internal/commands/worker

echo ""
echo -e "${GREEN}✅ Swagger documentation generated successfully!${NC}"
echo ""
echo -e "${BLUE}📁 Generated files:${NC}"
ls -lh internal/adapters/http/docs/

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🚀 To view the documentation:${NC}"
echo -e "   1. Start the API server: ${YELLOW}go run main.go api${NC}"
echo -e "   2. Open in browser: ${YELLOW}http://localhost:40012/swagger/index.html${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
