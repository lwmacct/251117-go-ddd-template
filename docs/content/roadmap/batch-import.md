# 批量用户导入功能

> **状态**: ✅ 前后端均已完成
> **优先级**: 中
> **完成日期**: 2024-11-30 (前端) / 2025-11-30 (后端)

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:26+4`
- [技术方案](#技术方案) `:30+33`
  - [后端 API](#后端-api) `:32+16`
  - [前端实现](#前端实现) `:48+7`
  - [CSV 格式](#csv-格式) `:55+8`
- [开发计划](#开发计划) `:63+27`
  - [步骤 1: CSV 解析工具](#步骤-1-csv-解析工具) `:65+6`
  - [步骤 2: 导入对话框组件](#步骤-2-导入对话框组件) `:71+7`
  - [步骤 3: 页面集成](#步骤-3-页面集成) `:78+6`
  - [步骤 4: API 集成](#步骤-4-api-集成) `:84+6`
- [进度记录](#进度记录) `:90+10`
- [相关文件](#相关文件) `:100+13`

<!--TOC-->

## 需求背景

管理员需要批量导入用户，从 CSV 文件读取用户数据并批量创建。当前只能逐个创建用户，效率较低。

## 技术方案

### 后端 API

> **已完成**: 后端 API 已实现 (2025-11-30)

```
POST /api/admin/users/batch
Request Body: { "users": [...] }
Response: { "success": number, "failed": number, "errors": [...] }
```

**后端实现文件**:

- `internal/application/user/command/batch_create_users.go` - Command 定义
- `internal/application/user/command/batch_create_users_handler.go` - Handler 实现
- `internal/application/user/command/batch_create_users_handler_test.go` - 单元测试

### 前端实现

1. 创建 CSV 解析工具
2. 创建导入对话框组件
3. 显示导入预览和验证结果
4. 调用 API 批量创建

### CSV 格式

```csv
username,email,password,full_name,status
user1,user1@example.com,Password123,用户一,active
user2,user2@example.com,Password123,用户二,active
```

## 开发计划

### 步骤 1: CSV 解析工具

- [x] 在 `web/src/utils/import.ts` 创建解析函数
- [x] 支持 UTF-8 编码
- [x] 数据验证（用户名、邮箱、密码格式）

### 步骤 2: 导入对话框组件

- [x] 创建 `UserImportDialog.vue`
- [x] 文件上传和拖拽预览
- [x] 验证结果展示（统计 + 错误详情）
- [x] 三步向导流程（上传 → 预览 → 结果）

### 步骤 3: 页面集成

- [x] 添加导入按钮
- [x] 集成对话框
- [x] 导入完成后刷新列表

### 步骤 4: API 集成

- [x] 添加批量创建 API 类型定义
- [x] 错误处理
- [x] 后端实现 `POST /api/admin/users/batch`

## 进度记录

| 日期       | 步骤                    | 状态    | 备注                                |
| ---------- | ----------------------- | ------- | ----------------------------------- |
| 2024-11-30 | 步骤 1: CSV 解析工具    | ✅ 完成 | 支持 RFC 4180 规范、数据验证        |
| 2024-11-30 | 步骤 2: 导入对话框      | ✅ 完成 | 三步向导、拖拽上传                  |
| 2024-11-30 | 步骤 3: 页面集成        | ✅ 完成 | 添加按钮、自动刷新                  |
| 2024-11-30 | 步骤 4: API 集成 (前端) | ✅ 完成 | 前端 API 调用已就绪                 |
| 2025-11-30 | 步骤 5: 后端实现        | ✅ 完成 | DDD+CQRS 模式实现，支持部分失败处理 |

## 相关文件

```
web/src/
├── utils/
│   └── import.ts                    # CSV 解析工具（新建）
├── api/admin/
│   └── users.ts                     # 添加批量导入 API
├── pages/admin/users/
│   ├── index.vue                    # 添加导入按钮
│   └── components/
│       └── UserImportDialog.vue     # 导入对话框（新建）
```
