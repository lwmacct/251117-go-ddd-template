---
paths: "**/*.go"
---

# 运行时规范

- `pkill -f go-build` - 杀掉临时编译进程

## 进程管理

- ❌ **禁止杀掉 air 进程** - air 是开发环境的热重载服务

### 验证编译

服务无响应或需要验证 DI 问题时，使用 `timeout` 命令：

```bash
# 测试服务编译和启动（3秒后自动终止）
timeout 3 go run . api

# 测试特定命令
timeout 3 go run . db migrate
```

**判断标准**：

- 看到 `Starting API server` → DI 成功
- 端口占用错误 (`Address already in use`) → 服务已在运行，编译成功
- 其他错误 → DI 或配置问题

### 服务状态检查

```bash
# 检查服务是否运行
ps aux | grep -E "go run|main|air" | grep -vE "grep|node"


```
