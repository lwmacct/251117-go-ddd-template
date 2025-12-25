# 测试指南

本文档介绍项目的测试策略。

<!--TOC-->

## Table of Contents

- [测试框架](#测试框架) `:17+5`
- [测试类型](#测试类型) `:22+9`
- [运行测试](#运行测试) `:31+20`
- [测试规范](#测试规范) `:51+9`
- [测试覆盖率目标](#测试覆盖率目标) `:60+8`

<!--TOC-->

## 测试框架

- **testify**: 断言库和 Mock 框架
- **go test**: Go 标准测试工具

## 测试类型

| 类型            | 位置                                   | 说明                         |
| --------------- | -------------------------------------- | ---------------------------- |
| 单元测试        | `*_test.go` 同目录                     | 测试单个函数/方法            |
| Use Case 测试   | `application/*/command/*_test.go`      | Mock Repository 测试业务逻辑 |
| Repository 测试 | `infrastructure/persistence/*_test.go` | 真实数据库集成测试           |
| HTTP 测试       | `adapters/http/handler/*_test.go`      | Mock Handler 测试 HTTP 层    |

## 运行测试

```bash
# 运行所有测试
go test ./...

# 详细输出
go test -v ./...

# 特定目录
go test ./internal/domain/user/...

# 特定测试函数
go test -run TestUser_CanLogin ./internal/domain/user/

# 覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 测试规范

1. **AAA 模式**: Arrange（准备）→ Act（执行）→ Assert（断言）
2. **表驱动测试**: 使用 `[]struct` 定义多组测试用例
3. **断言选择**:
   - `require.*` - 前置条件失败时立即终止
   - `assert.*` - 验证结果，失败继续执行
4. **Mock 实现**: 使用 `testify/mock` 隔离依赖

## 测试覆盖率目标

| 层级              | 目标 |
| ----------------- | ---- |
| Domain 层         | 90%+ |
| Application 层    | 80%+ |
| Infrastructure 层 | 70%+ |
| Adapters 层       | 60%+ |
