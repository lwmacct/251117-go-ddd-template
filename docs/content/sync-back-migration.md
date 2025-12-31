# 副本改动回迁方法论

## 概述

本文档定义了从副本项目（fork/基于模板创建的业务项目）将改动回迁到模板项目的标准流程和方法。

---

## 适用场景

| 场景          | 是否回迁 | 示例                   |
| ------------- | -------- | ---------------------- |
| 模板 Bug 修复 | ✅ 是    | 修复审计日志表名错误   |
| 通用功能      | ✅ 是    | 自定义验证器、工具函数 |
| 架构优化      | ✅ 是    | 简化依赖注入、优化模式 |
| 业务功能      | ❌ 否    | 特定业务逻辑、领域模型 |

---

## 标准流程

### Phase 1: 对比分析

```bash
# 1. 确定回迁范围（以模板基准提交为起点）
cd /path/to/template && git log -1 --oneline   # 模板基准: abc1234

# 2. 查看副本项目在基准之后的提交
cd /path/to/fork
git log --oneline abc1234..HEAD          # 基准之后的所有提交
git log --oneline abc1234..HEAD --graph  # 树形结构展示

# 3. 对比两个项目的差异
# 方法一：直接对比目录（跨项目）
git diff --stat /path/to/template /path/to/fork -- internal/
git diff /path/to/template /path/to/fork -- internal/application/org/

# 方法二：在同一目录对比（推荐）
cd /path/to/fork
git diff-template -- internal/  # 假设有 template 作为 remote
```

### Phase 2: 分类改动

将改动分为以下类别：

| 类别       | 说明                         |
| ---------- | ---------------------------- |
| **回迁类** | Bug 修复、通用功能、架构优化 |
| **保留类** | 业务逻辑、特定配置、领域模型 |

### Phase 3: 确定回迁范围

```
改动是否修复了模板问题？
├─ 是 → 回迁
└─ 否 → 改动是否对所有项目通用？
    ├─ 是 → 回迁
    └─ 否 → 不回迁（保留在副本）
```

### Phase 4: 执行回迁

按依赖顺序执行：

```
Infrastructure → Domain → Application → Adapters → Container
```

### Phase 5: 验证

```bash
# 编译检查
go build -o /dev/null ./...

# 单元测试
go test ./...

# 集成测试（需要服务运行）
MANUAL=1 go test -v ./internal/manualtest/...
```

---

## 回迁技巧

### 1. 文件复制与路径替换

**单文件**（一步到位）：

```bash
sed 's/fork-repo/template-repo/g' \
  fork/internal/file.go > template/internal/file.go
```

**批量复制目录**：

```bash
# 复制目录
cp -r fork/internal/manualtest/org template/internal/manualtest/

# 批量替换
find template/internal/manualtest/org -name "*.go" \
  -exec sed -i 's/fork-repo/template-repo/g' {} \;
```

### 2. 分层迁移策略

| 层级           | 迁移要点              | 验证方式   |
| -------------- | --------------------- | ---------- |
| Infrastructure | 最底层，先迁移        | `go build` |
| Domain         | 接口定义，小心修改    | `go build` |
| Application    | 业务逻辑，依赖 Domain | `go test`  |
| Adapters       | HTTP 暴露，最后改     | `go test`  |
| Container      | 依赖注入，最后调整    | `go build` |

### 3. 测试文件处理

测试文件需要修改：

1. **包路径** - 批量替换
2. **API 路由** - 如有变化需同步
3. **业务数据** - 去除特定字段，保留通用部分

### 4. 配置文件差异

配置文件通常**不回迁**，需手动检查：

- `config.yaml` - 环境特定配置
- `.env` - 不提交
- Docker 配置 - 可能需要合并

---

## 常见问题

### Q: 回迁后编译失败，循环依赖？

**A**: 检查依赖方向，DDD 架构要求：

```
正确: Application → Domain ← Infrastructure
错误: Application → Infrastructure → Domain
```

### Q: Container 层构造函数参数太多？

**A**: 使用仓储聚合模式：

```go
// 分散注入（参数多）
type Params struct {
    Repo1 persistence.Repo1
    Repo2 persistence.Repo2
    // ...
}

// 聚合注入（参数少）
type Params struct {
    Repos persistence.AllRepositories
}
```

### Q: 测试文件太多，逐个替换太慢？

**A**: 批量操作：

```bash
# 批量复制
cp -r fork/manualtest/* template/manualtest/

# 批量替换
find template/manualtest -name "*.go" -exec sed -i 's/fork/template/g' {} +
```

---

## 检查清单

### 回迁前

- [ ] 确认改动是通用功能，非业务特定
- [ ] 模板项目已拉取最新代码
- [ ] 创建备份分支

### 回迁中

- [ ] 按依赖顺序迁移
- [ ] 每层完成后编译验证
- [ ] 更新相关文档

### 回迁后

- [ ] 编译通过
- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 更新 CHANGELOG

---
