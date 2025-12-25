#!/usr/bin/env bash
#
# 检查 Application 层文件和结构体命名规范
# 用于 pre-commit hook
#
# 规则：
# - cmd_*.go 文件中的 struct 必须以 Command 结尾
# - qry_*.go 文件中的 struct 必须以 Query 结尾
# - dto.go 文件中的 struct 必须以 DTO 结尾
# - 只检查 internal/application/ 目录下的文件

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

errors=0
checked=0

check_struct_suffix() {
    local file="$1"
    local required_suffix="$2"
    local file_type="$3"

    # 提取所有 struct 类型声明
    while IFS= read -r line; do
        # 提取类型名
        type_name=$(echo "$line" | sed -n 's/.*type[[:space:]]\+\([A-Za-z_][A-Za-z0-9_]*\)[[:space:]]\+struct.*/\1/p')

        if [[ -z "$type_name" ]]; then
            continue
        fi

        # 检查是否以要求的后缀结尾
        if [[ ! "$type_name" =~ ${required_suffix}$ ]]; then
            if [[ $errors -eq 0 ]]; then
                echo -e "${RED}Application layer naming convention violation detected:${NC}"
            fi
            echo -e "  ${YELLOW}${file}${NC}: ${RED}${type_name}${NC} should end with '${required_suffix}' (${file_type} file)"
            errors=$((errors + 1))
        fi
    done < <(grep -E '^[[:space:]]*type[[:space:]]+[A-Za-z_][A-Za-z0-9_]*[[:space:]]+struct' "$file" 2>/dev/null || true)
}

for file in "$@"; do
    # 只处理 internal/application 目录下的文件
    if [[ ! "$file" =~ ^internal/application/ ]]; then
        continue
    fi

    # 跳过测试文件
    if [[ "$file" =~ _test\.go$ ]]; then
        continue
    fi

    # 跳过 handler 文件（handler 文件的命名规范单独处理）
    if [[ "$file" =~ _handler\.go$ ]]; then
        continue
    fi

    if [[ ! -f "$file" ]]; then
        continue
    fi

    checked=$((checked + 1))
    filename=$(basename "$file")

    # 根据文件前缀/名称检查结构体命名
    if [[ "$filename" =~ ^cmd_ ]]; then
        check_struct_suffix "$file" "Command" "cmd_*"
    elif [[ "$filename" =~ ^qry_ ]]; then
        check_struct_suffix "$file" "Query" "qry_*"
    elif [[ "$filename" == "dto.go" ]]; then
        check_struct_suffix "$file" "DTO" "dto"
    fi
done

if [[ $errors -gt 0 ]]; then
    echo ""
    echo -e "${RED}Found ${errors} struct(s) not following naming convention.${NC}"
    echo -e "Rules:"
    echo -e "  - cmd_*.go files: structs must end with 'Command'"
    echo -e "  - qry_*.go files: structs must end with 'Query'"
    echo -e "  - dto.go files: structs must end with 'DTO'"
    exit 1
fi

if [[ $checked -gt 0 ]]; then
    echo -e "${GREEN}All ${checked} application layer file(s) follow naming convention.${NC}"
fi

exit 0
