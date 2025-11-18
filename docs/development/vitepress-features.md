# VitePress 2.0 功能展示

本页面展示 VitePress 2.0 的原生功能，无需安装任何插件。

## 📛 Badge 徽章

使用徽章标注版本、状态等信息。

### 用户管理 <Badge type="tip" text="v2.0" />

用户 CRUD 操作接口。

### 批量导入 <Badge type="info" text="新功能" />

支持批量导入用户数据。

### 旧版 API <Badge type="warning" text="已废弃" />

请使用新版 API。

### 实验性功能 <Badge type="danger" text="实验性" />

此功能仍在测试中。

## 📝 代码块高亮

### 行高亮

高亮指定行：

```go {2,4-6}
func CreateUser(user *User) error {
    // 验证用户数据
    if err := user.Validate(); err != nil {
        return err
    }
    // 保存到数据库
    return db.Create(user).Error
}
```

### 代码聚焦

聚焦重点代码：

```typescript
export default defineConfig({
  themeConfig: {
    search: { // [!code focus]
      provider: 'local' // [!code focus]
    } // [!code focus]
  }
})
```

### 代码差异

显示代码的增删改：

```go
func (r *UserRepository) Create(user *User) error {
    return r.db.Create(user).Error // [!code --]
    // 添加事务支持
    return r.db.Transaction(func(tx *gorm.DB) error { // [!code ++]
        return tx.Create(user).Error // [!code ++]
    }) // [!code ++]
}
```

### 错误和警告标记

```go
func ConnectDB(url string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(url), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent), // [!code warning]
    })

    if err != nil {
        panic(err) // [!code error]
    }

    return db, nil
}
```

## 📦 自定义容器

### 嵌套容器

::: details 点击查看完整配置
::: code-group

```yaml [开发环境]
server:
  addr: :8080
  debug: true

database:
  url: postgresql://localhost:5432/dev
```

```yaml [生产环境]
server:
  addr: :80
  debug: false

database:
  url: postgresql://db.example.com:5432/prod
```

:::
:::

### 自定义标题

::: tip 💡 最佳实践
始终在生产环境中禁用 debug 模式。
:::

::: warning ⚠️ 注意事项
修改配置后需要重启服务器。
:::

::: danger 🚨 安全警告
不要在代码中硬编码数据库密码！
:::

## 📄 文件名显示

代码块可以显示文件名：

```go [internal/domain/user/model.go]
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"uniqueIndex"`
    Email     string    `gorm:"uniqueIndex"`
    Password  string    `json:"-"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

```typescript [.vitepress/config.ts]
export default defineConfig({
  title: "Go DDD Template",
  description: "基于 Go 的 DDD 模板应用"
})
```

## 📊 表格对齐

| 功能 | 状态 | 优先级 |
| :--- | :---: | ---: |
| 用户认证 | ✅ 已完成 | 高 |
| 权限管理 | 🚧 进行中 | 高 |
| 日志系统 | 📋 计划中 | 中 |
| 监控告警 | 💭 待定 | 低 |

## 🎯 任务列表

- [x] 完成用户 CRUD 接口
- [x] 实现 JWT 认证
- [x] 集成 PostgreSQL
- [x] 集成 Redis
- [ ] 添加单元测试
- [ ] 添加集成测试
- [ ] 完善 API 文档
- [ ] 部署到生产环境

## 😊 Emoji 支持

:tada: 项目初始化
:rocket: 部署到生产
:bug: 修复认证 bug
:sparkles: 添加缓存功能
:fire: 移除废弃代码
:lock: 修复安全漏洞
:memo: 更新文档
:white_check_mark: 添加测试

## 🔗 链接和引用

### 内部链接

- [快速开始](/guide/getting-started)
- [项目架构](/guide/architecture)
- [API 文档](/api/)

### 外部链接

- [VitePress 官方文档](https://vitepress.dev/)
- [Go 官方网站](https://go.dev/)
- [GitHub 仓库](https://github.com/lwmacct/251117-go-ddd-template)

## 📸 图片

![Go Logo](https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png)

## 🎨 提示和警告

::: tip 提示
使用环境变量管理配置，避免硬编码。
:::

::: warning 警告
生产环境中务必关闭 debug 模式。
:::

::: danger 危险
不要将敏感信息提交到 Git 仓库。
:::

::: details 更多信息
VitePress 2.0 支持所有 Markdown 扩展语法，包括表格、任务列表、Emoji 等。
:::

## 💡 使用建议

以上所有功能都是 VitePress 2.0 **原生支持**的，无需安装任何插件！

- ✅ Badge 徽章：标注版本、状态
- ✅ 代码高亮：突出重点代码
- ✅ 代码差异：展示变更
- ✅ 自定义容器：组织内容
- ✅ 文件名显示：明确代码来源
- ✅ 任务列表：跟踪进度
- ✅ Emoji：增加趣味性

## 📚 参考资料

- [VitePress Markdown 扩展](https://vitepress.dev/guide/markdown)
- [默认主题配置](https://vitepress.dev/reference/default-theme-config)
