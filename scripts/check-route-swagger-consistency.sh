#!/usr/bin/env bash
#
# 检查 router.go 与 handler @Router/@x-permission 注解的一致性
# 用于 pre-commit hook，防止路由/权限不匹配
#

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

ROUTER_FILE="internal/adapters/http/router.go"
HANDLER_DIR="internal/adapters/http/handler"

errors=0

# 检查必要文件
if [[ ! -f "$ROUTER_FILE" ]]; then
    echo -e "${YELLOW}Warning: $ROUTER_FILE not found, skipping${NC}"
    exit 0
fi

if [[ ! -d "$HANDLER_DIR" ]]; then
    echo -e "${YELLOW}Warning: $HANDLER_DIR not found, skipping${NC}"
    exit 0
fi

echo -e "${CYAN}Checking route and permission consistency...${NC}"

# 临时文件
HANDLER_DATA=$(mktemp)
ROUTER_DATA=$(mktemp)
trap "rm -f $HANDLER_DATA $ROUTER_DATA" EXIT

# ============================================================
# 从 handler/*.go 提取 @Router + @x-permission
# 格式: method|path|permission
# ============================================================
extract_handler_routes() {
    for file in "$HANDLER_DIR"/*.go; do
        [[ -f "$file" ]] || continue
        [[ "$(basename "$file")" == "CLAUDE.md" ]] && continue

        awk '
        # 收集注释块中的 @Router 和 @x-permission
        /^\/\// {
            if (match($0, /@Router[[:space:]]+([^[:space:]]+)[[:space:]]+\[([a-z]+)\]/, arr)) {
                route_method = arr[2]
                route_path = arr[1]
            }
            if (match($0, /@x-permission[[:space:]]+\{"scope":"([^"]+)"\}/, arr)) {
                permission = arr[1]
            }
        }

        # 遇到 func 定义时输出收集的信息
        /^func[[:space:]]+\(/ {
            if (route_method != "" && route_path != "") {
                if (permission == "") permission = "-"
                print route_method "|" route_path "|" permission
            }
            route_method = ""
            route_path = ""
            permission = ""
        }
        ' "$file"
    done
}

# ============================================================
# 从 router.go 提取路由 + RequirePermission
# 格式: method|path|permission
# ============================================================
extract_router_routes() {
    awk '
    BEGIN {
        prefix["r"] = ""
        prefix["api"] = "/api"
    }

    # Group 定义
    /[a-zA-Z_][a-zA-Z0-9_]*[[:space:]]*:=[[:space:]]*[a-zA-Z_][a-zA-Z0-9_]*\.Group\(/ {
        match($0, /([a-zA-Z_][a-zA-Z0-9_]*)[[:space:]]*:=[[:space:]]*([a-zA-Z_][a-zA-Z0-9_]*)\.Group\("([^"]+)"\)/, arr)
        if (arr[1] != "" && arr[2] != "" && arr[3] != "") {
            if (arr[2] in prefix) {
                prefix[arr[1]] = prefix[arr[2]] arr[3]
            }
        }
    }

    # 路由定义
    /\.(GET|POST|PUT|DELETE|PATCH)\("/ {
        match($0, /([a-zA-Z_][a-zA-Z0-9_]*)\.(GET|POST|PUT|DELETE|PATCH)\("([^"]+)"/, arr)
        if (arr[1] != "" && arr[2] != "" && arr[3] != "") {
            var_name = arr[1]
            method = tolower(arr[2])
            route_path = arr[3]

            # 构建完整路径
            if (var_name in prefix) {
                full_path = prefix[var_name] route_path
            } else {
                full_path = route_path
            }

            # :param -> {param}
            while (match(full_path, /:([a-zA-Z_][a-zA-Z0-9_]*)/, param)) {
                gsub(":" param[1], "{" param[1] "}", full_path)
            }

            # 提取 RequirePermission
            permission = "-"
            if (match($0, /RequirePermission\("([^"]+)"\)/, perm)) {
                permission = perm[1]
            }

            print method "|" full_path "|" permission
        }
    }
    ' "$ROUTER_FILE"
}

# 提取数据
extract_handler_routes | sort -u > "$HANDLER_DATA"
extract_router_routes | sort -u > "$ROUTER_DATA"

# ============================================================
# 检查路由一致性
# ============================================================
check_routes() {
    local handler_routes router_routes

    # 提取纯路由（method|path）
    handler_routes=$(cut -d'|' -f1,2 "$HANDLER_DATA" | sort -u)
    router_routes=$(cut -d'|' -f1,2 "$ROUTER_DATA" | sort -u)

    # Handler 有但 Router 没有
    while IFS='|' read -r method path; do
        if ! echo "$router_routes" | grep -qxF "${method}|${path}"; then
            if [[ $errors -eq 0 ]]; then
                echo ""
                echo -e "${RED}Inconsistency detected:${NC}"
            fi
            echo -e "  ${YELLOW}[@Router only]${NC} $method $path"
            echo -e "    ${CYAN}→ Add route to router.go${NC}"
            ((errors++)) || true
        fi
    done <<< "$handler_routes"

    # Router 有但 Handler 没有
    while IFS='|' read -r method path; do
        [[ "$path" =~ swagger|/docs ]] && continue
        if ! echo "$handler_routes" | grep -qxF "${method}|${path}"; then
            if [[ $errors -eq 0 ]]; then
                echo ""
                echo -e "${RED}Inconsistency detected:${NC}"
            fi
            echo -e "  ${YELLOW}[router.go only]${NC} $method $path"
            echo -e "    ${CYAN}→ Add @Router annotation to handler${NC}"
            ((errors++)) || true
        fi
    done <<< "$router_routes"
}

# ============================================================
# 检查权限一致性
# ============================================================
check_permissions() {
    # 对于同一路由，比较权限
    while IFS='|' read -r method path handler_perm; do
        # 在 router 数据中查找对应路由
        router_line=$(grep "^${method}|${path}|" "$ROUTER_DATA" 2>/dev/null || true)
        if [[ -n "$router_line" ]]; then
            router_perm=$(echo "$router_line" | cut -d'|' -f3)

            # 比较权限
            if [[ "$handler_perm" != "$router_perm" ]]; then
                if [[ $errors -eq 0 ]]; then
                    echo ""
                    echo -e "${RED}Inconsistency detected:${NC}"
                fi
                echo -e "  ${YELLOW}[Permission mismatch]${NC} $method $path"
                echo -e "    ${CYAN}@x-permission:${NC} $handler_perm"
                echo -e "    ${CYAN}RequirePermission:${NC} $router_perm"
                ((errors++)) || true
            fi
        fi
    done < "$HANDLER_DATA"
}

# 执行检查
check_routes
check_permissions

# ============================================================
# 结果
# ============================================================
if [[ $errors -gt 0 ]]; then
    echo ""
    echo -e "${RED}Found $errors inconsistency(ies).${NC}"
    exit 1
fi

echo -e "${GREEN}All routes and permissions are consistent.${NC}"
exit 0
