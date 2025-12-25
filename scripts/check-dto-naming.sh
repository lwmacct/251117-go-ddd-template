#!/usr/bin/env bash
#
# 检查 dto.go 文件中的结构体命名是否以 DTO 结尾
# 用于 pre-commit hook
#

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

errors=0
checked=0

for file in "$@"; do
    # 只处理 dto.go 文件
    if [[ ! "$file" =~ dto\.go$ ]]; then
        continue
    fi

    if [[ ! -f "$file" ]]; then
        continue
    fi

    checked=$((checked + 1))

    # 提取所有 struct 类型声明: type XxxName struct
    # 使用 grep 匹配 "type ... struct" 模式
    while IFS= read -r line; do
        # 提取类型名（type 和 struct 之间的单词）
        type_name=$(echo "$line" | sed -n 's/.*type[[:space:]]\+\([A-Za-z_][A-Za-z0-9_]*\)[[:space:]]\+struct.*/\1/p')

        if [[ -z "$type_name" ]]; then
            continue
        fi

        # 检查是否以 DTO 结尾
        if [[ ! "$type_name" =~ DTO$ ]]; then
            if [[ $errors -eq 0 ]]; then
                echo -e "${RED}DTO naming convention violation detected:${NC}"
            fi
            echo -e "  ${YELLOW}${file}${NC}: ${RED}${type_name}${NC} should end with 'DTO'"
            errors=$((errors + 1))
        fi
    done < <(grep -E '^[[:space:]]*type[[:space:]]+[A-Za-z_][A-Za-z0-9_]*[[:space:]]+struct' "$file" 2>/dev/null || true)
done

if [[ $errors -gt 0 ]]; then
    echo ""
    echo -e "${RED}Found ${errors} struct(s) not following DTO naming convention.${NC}"
    echo -e "All structs in dto.go files must end with 'DTO'."
    echo -e "Example: UserResponse -> UserResponseDTO"
    exit 1
fi

if [[ $checked -gt 0 ]]; then
    echo -e "${GREEN}All ${checked} dto.go file(s) follow naming convention.${NC}"
fi

exit 0
