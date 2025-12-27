---
paths:
  - "internal/infrastructure/twofa/**/*.go"
---

# TwoFA Infrastructure 规范

## 核心职责

实现 Domain 层 `twofa.Service` 接口，提供基于 TOTP (RFC 6238) 的双因素认证功能。

**兼容认证器**：Google Authenticator、Microsoft Authenticator、Authy、1Password 等。

## 文件结构

| 文件         | 职责               |
| ------------ | ------------------ |
| `doc.go`     | 包文档（**必需**） |
| `service.go` | 2FA 服务完整实现   |

## 设计原则

- **接口实现**：实现 `domain/twofa.Service` 接口
- **依赖注入**：通过构造函数注入 Command/Query Repository 和 UserQueryRepo
- **二维码输出**：返回 Base64 PNG，可直接用于 `<img src="...">`

## 核心方法

| 方法              | 说明                         |
| ----------------- | ---------------------------- |
| `Setup`           | 生成 TOTP 密钥和二维码       |
| `VerifyAndEnable` | 首次验证并启用，返回恢复码   |
| `Verify`          | 验证 TOTP 码或恢复码         |
| `Disable`         | 禁用 2FA                     |
| `GetStatus`       | 查询启用状态和剩余恢复码数量 |

## 安全设计

| 特性         | 说明                         |
| ------------ | ---------------------------- |
| **密钥长度** | 80 位（10 字节）随机数       |
| **恢复码**   | 一次性使用，验证后自动删除   |
| **存储安全** | 密钥存储在数据库，不对外暴露 |

## 使用流程

1. `Setup(userID)` → 返回二维码
2. 用户扫码配置认证器
3. `VerifyAndEnable(code)` → 返回恢复码（用户需妥善保存）
4. 后续登录 `Verify(code)` → 验证 TOTP 或恢复码

## 依赖

- `github.com/pquerna/otp/totp`：TOTP 实现
- `github.com/skip2/go-qrcode`：二维码生成

## 依赖方向

```
domain/twofa.Service (接口)
           ↑
infrastructure/twofa (实现)
```
