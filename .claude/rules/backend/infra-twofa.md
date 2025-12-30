---
paths:
  - "internal/infrastructure/twofa/**/*.go"
---

# TwoFA Infrastructure 规范

基于 TOTP (RFC 6238) 的双因素认证。

## 目录结构

```
internal/infrastructure/twofa/
├── doc.go      # 包文档（必需）
└── service.go  # 服务实现
```

## 使用流程

1. `Setup()` → 生成密钥和二维码
2. 用户扫码配置认证器
3. `VerifyAndEnable()` → 返回恢复码
4. `Verify()` → 验证 TOTP 或恢复码

## 安全设计

- 密钥：80 位随机数
- 恢复码：一次性，验证后删除
- 二维码：Base64 PNG 输出
