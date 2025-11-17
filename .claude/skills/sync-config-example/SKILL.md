---
name: sync-config-example
description: 当修改 internal/infrastructure/config/config.go 中的 Config 结构体、defaultConfig() 函数，或用户要求"同步配置"、"更新配置示例"时使用。保持 configs/config.example.yaml 与代码定义一致。
---

# 配置示例文件同步器

从 `internal/infrastructure/config/config.go` 中的 `defaultConfig()` 函数提取默认配置，生成标准的 `configs/config.example.yaml` 示例文件。

## 使用场景

- 修改了 `defaultConfig()` 函数中的默认值
- 在 `Config` 结构体中添加、删除或修改了配置字段
- 更改了配置字段的类型或 koanf 标签
- 需要重新生成完整的配置示例文件

## 工作流程

### 1. 读取并解析配置代码
- 读取 `internal/infrastructure/config/config.go`
- 定位 `defaultConfig()` 函数
- 解析 Config 结构体及所有嵌套结构体的定义
- 提取结构体字段的注释作为文档

### 2. 提取默认值
- 从 `defaultConfig()` 函数体中提取所有字段的默认值
- 识别字段类型：string、int、bool、time.Duration 等
- 处理嵌套结构体的初始化

### 3. 生成 YAML 配置文件
生成符合以下规范的 `config.example.yaml`：

**文件头部：**
```yaml
# 示例配置文件
# 复制此文件为 config.yaml 并根据需要修改
#
# 此文件与 internal/infrastructure/config/config.go 中的 defaultConfig() 保持同步
# 所有配置项都可以通过环境变量覆盖
```

**配置内容：**
- 按结构体层级组织为 YAML 嵌套结构
- 每个配置段添加分组注释（如：`# 服务器配置`）
- 每个字段添加行内注释说明用途、可选值等
- 敏感字段（password、secret 等）添加环境变量建议
- 使用 2 空格缩进，保持格式整洁

### 4. 处理特殊类型

**time.Duration 转换：**
将 Go 的 Duration 表达式转换为字符串格式

| Go 代码 | YAML 值 |
|---------|---------|
| `time.Minute` | `"1m"` |
| `15 * time.Minute` | `"15m"` |
| `time.Hour` | `"1h"` |
| `24 * time.Hour` | `"24h"` |
| `7 * 24 * time.Hour` | `"168h"` |

**敏感字段处理：**
- 密码类字段（password）：值为空字符串，添加环境变量建议
- 密钥类字段（secret、key）：使用占位符值，添加安全警告

### 5. 写入文件
- 完全重写 `configs/config.example.yaml`（不保留原有内容）
- 确保 YAML 语法正确
- 文件使用 UTF-8 编码
- 文件末尾保留空行

## 字段映射规则

### Go → YAML 键名
- 使用结构体字段的 `koanf` 标签值作为 YAML 键
- 嵌套结构体对应 YAML 的层级结构
- 标签格式：`koanf:"field_name"`

**示例：**
```go
type Config struct {
    Server ServerConfig `koanf:"server"`
}

type ServerConfig struct {
    Addr string `koanf:"addr"`
    Port int    `koanf:"port"`
}
```

生成 YAML：
```yaml
server:
  addr: "0.0.0.0:8080"
  port: 8080
```

### 类型转换

| Go 类型 | YAML 格式 | 示例 |
|---------|-----------|------|
| `string` | 字符串（建议加引号） | `"localhost"` |
| `int`/`int64` | 整数 | `5432` |
| `time.Duration` | 时间字符串 | `"15m"`, `"24h"` |
| `bool` | 布尔值 | `true`, `false` |

## 注释生成规范

### 1. 文件头注释
- 说明文件用途（示例配置）
- 说明如何使用（复制为 config.yaml）
- 说明与代码的同步关系
- 说明环境变量覆盖机制

### 2. 分组注释
每个顶层配置段添加分组注释：
```yaml
# [功能名称] 配置
section_name:
  key: value
```

### 3. 字段注释
- 优先使用 Go 代码中的字段注释
- 对于枚举值，列出所有可选项（如：`development | production`）
- 对于格式要求，说明格式（如：`格式: host:port`）
- 对于默认值，解释其含义

### 4. 安全提示
敏感字段添加提示：
```yaml
password: "" # 建议通过环境变量 APP_DATABASE_PASSWORD 设置
secret: "change-me-in-production" # 生产环境务必修改!
```

## 验证检查清单

生成文件后，自动验证以下内容：

- [ ] YAML 语法正确（可解析）
- [ ] 所有 defaultConfig() 中的字段都已包含
- [ ] Duration 类型正确转换为字符串格式
- [ ] 敏感字段有安全提示注释
- [ ] 文件编码为 UTF-8
- [ ] 缩进统一使用 2 个空格
- [ ] 嵌套层级结构正确

## 执行示例

### 用户请求
```
同步配置文件
```

### 执行步骤
1. 读取 `internal/infrastructure/config/config.go`
2. 解析 `defaultConfig()` 函数和相关结构体
3. 提取所有默认值和字段信息
4. 转换特殊类型（Duration、敏感字段等）
5. 生成完整的 YAML 内容
6. 写入 `configs/config.example.yaml`
7. 报告完成情况

### 完成报告
```
已从 config.go 读取默认配置
已解析配置结构：Config 及其嵌套结构体
已转换 time.Duration 类型字段
已生成 configs/config.example.yaml
同步完成
```

## 错误处理

| 错误情况 | 处理方式 |
|---------|----------|
| 无法读取 config.go | 报告文件路径错误，停止执行 |
| 无法解析 defaultConfig() | 提示检查函数格式，报告解析位置 |
| Duration 转换失败 | 使用原始表达式，添加警告注释 |
| 写入文件失败 | 报告权限或路径错误 |
| YAML 格式错误 | 报告语法问题，提供修正建议 |

## 注意事项

- **完全重写**：此 skill 会完全重写 `config.example.yaml`，不保留手动添加的内容
- **代码优先**：以 Go 代码中的 defaultConfig() 为唯一数据源
- **仅作示例**：生成的文件仅供参考，不包含敏感信息的实际值
- **需要审查**：建议同步后检查生成的文件，确认注释和格式符合预期
- **环境变量**：确保文件头注释中正确说明环境变量前缀（通常为 `APP_`）
