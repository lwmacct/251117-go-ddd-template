---
paths:
  - "internal/infrastructure/validation/**/*.go"
---

# Validation Infrastructure 规范

## 核心职责

实现 Domain 层 `setting.Validator` 接口，为 Setting 模块提供基于 JSON Logic 的动态规则验证。

## 文件结构

| 文件                          | 职责               |
| ----------------------------- | ------------------ |
| `doc.go`                      | 包文档（**必需**） |
| `jsonlogic_validator.go`      | JSON Logic 验证器  |
| `jsonlogic_validator_test.go` | 单元测试           |

## 设计原则

- **接口实现**：实现 `domain/setting.Validator` 接口
- **双格式支持**：JSON Logic 格式（推荐）+ 简单格式（向后兼容）
- **跨字段验证**：支持引用其他设置值进行验证

## 规则格式

### JSON Logic 格式（推荐）

| 示例                                                 | 说明       |
| ---------------------------------------------------- | ---------- |
| `{">=": [{"var": "value"}, 6]}`                      | 最小值验证 |
| `{"and": [{">=": [...]}, {"<=": [...]}]}`            | 范围验证   |
| `{">": [{"var": "value"}, {"var": "settings.xxx"}]}` | 跨字段验证 |

### 简单格式（向后兼容）

| 字段        | 类型   | 说明           |
| ----------- | ------ | -------------- |
| `required`  | bool   | 是否必填       |
| `min`/`max` | number | 数值范围       |
| `minLength` | number | 字符串最小长度 |
| `maxLength` | number | 字符串最大长度 |
| `message`   | string | 自定义错误消息 |

## 核心方法

| 方法            | 说明       |
| --------------- | ---------- |
| `Validate`      | 验证单个值 |
| `ValidateBatch` | 批量验证   |

## 设计说明

- 本包专用于 Setting 模块的动态验证需求
- 规则存储在数据库中，支持运行时修改
- 如需通用表单验证，请使用 `go-playground/validator`

## 依赖

- `github.com/diegoholiveira/jsonlogic/v3`：JSON Logic 引擎

## 依赖方向

```
domain/setting.Validator (接口)
            ↑
infrastructure/validation (实现)
```
