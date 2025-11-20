# Swagger/OpenAPI Documentation

这个目录包含自动生成的 Swagger/OpenAPI 文档文件。

## 📁 文件说明

- **`docs.go`** - Swagger 文档的 Go 入口文件（由 swag 自动生成）
- **`swagger.json`** - OpenAPI 3.0 JSON 格式规范
- **`swagger.yaml`** - OpenAPI 3.0 YAML 格式规范

## 🚀 如何使用

### 1. 查看在线文档

启动 API 服务器后，访问：

```bash
http://localhost:40012/swagger/index.html
```

### 2. 重新生成文档

当你修改了 handler 文件中的 Swagger 注解后，运行：

手动运行:

```bash
swag init \
    -g internal/commands/api/api.go \
    -o internal/adapters/http/docs \
    --parseDependency \
    --parseInternal
```

### 3. API 测试

在 Swagger UI 中：

1. 点击右上角 **"Authorize"** 按钮
2. 输入 Bearer Token：`Bearer your_jwt_token_here`
3. 选择任意 API 端点，点击 **"Try it out"** 按钮
4. 填写参数，点击 **"Execute"** 执行测试

## 📝 添加新的 API 文档

在 handler 文件中添加 Swagger 注解：

```go
// @Summary      接口简要说明
// @Description  接口详细描述
// @Tags         分组标签
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID"
// @Param        request body YourRequestDTO true "请求体"
// @Success      200 {object} response.Response "成功响应"
// @Failure      400 {object} response.ErrorResponse "错误响应"
// @Router       /api/your-endpoint [post]
// @x-permission {"scope":"your:permission:scope"}
func (h *YourHandler) YourMethod(c *gin.Context) {
    // ...
}
```

然后重新生成文档。

## 🔧 配置说明

### 主文档信息

主文档信息定义在 `internal/commands/api/api.go` 文件顶部：

```go
// @title           Go DDD Template API
// @version         1.0.3
// @description     基于 DDD+CQRS 架构的 Go Web 应用 API 文档
// @host            localhost:40012
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### 支持的注解标签

- `@Summary` - 接口简要说明
- `@Description` - 接口详细描述
- `@Tags` - API 分组
- `@Accept` - 接受的 Content-Type
- `@Produce` - 响应的 Content-Type
- `@Security` - 安全认证方式
- `@Param` - 参数定义（query/path/body/header）
- `@Success` - 成功响应
- `@Failure` - 错误响应
- `@Router` - 路由定义
- `@x-permission` - 自定义权限标签（本项目扩展）

## ⚠️ 注意事项

1. **不要手动编辑** `docs.go`、`swagger.json`、`swagger.yaml` 文件
2. **所有文档更改** 都应通过修改源码中的注解来完成
3. **提交前** 确保重新生成文档以保持同步
4. **测试环境** 的 `@host` 配置在 `api.go` 中修改

## 📚 参考资源

- [Swag 文档](https://github.com/swaggo/swag)
- [OpenAPI 3.0 规范](https://swagger.io/specification/)
- [Gin Swagger 中间件](https://github.com/swaggo/gin-swagger)
