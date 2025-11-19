# VitePress 2.0 升级总结

## 升级日期

2025-11-18

## 版本变更

- **旧版本**: VitePress 1.6.4
- **新版本**: VitePress 2.0.0-alpha.13

## 升级原因

- 使用最新的 Vite 7 构建工具链
- 体验 VitePress 2.0 的新特性和性能优化
- 改进的 CJK (中日韩) 语言支持
- 更好的开发体验和 HMR 支持

## 环境要求

### Node.js 版本

- **最低要求**: Node.js 20.19.0 或 22.12.0+
- **当前环境**: Node.js 24.11.1 ✅

### 依赖包变更

```json
{
  "devDependencies": {
    "vitepress": "^2.0.0-alpha.13", // 从 ^1.6.4 升级
    "vue": "^3.5.24" // 保持不变
  }
}
```

## 主要变更

### 1. 配置文件更新 (`docs/.vitepress/config.ts`)

#### 新增配置项

```typescript
markdown: {
  lineNumbers: true,
  // VitePress 2.0 新增：CJK 友好的强调语法 (默认启用)
  cjkFriendlyEmphasis: true,
  // 图片懒加载
  image: {
    lazyLoading: true,
  },
},

// Vite 配置 (如需自定义)
vite: {
  // Vite 7 配置选项
},
```

#### 配置说明

- `cjkFriendlyEmphasis`: 改进中文、日文、韩文的强调语法处理 (原 `cjkFriendly` 重命名)
- `image.lazyLoading`: 启用图片懒加载，优化页面性能
- `vite`: 可以直接配置 Vite 7 选项

### 2. package.json 更新

#### 版本号升级

```json
{
  "version": "2.0.0" // 从 1.0.0 升级
}
```

#### 新增 engines 字段

```json
{
  "engines": {
    "node": ">=20.19.0"
  }
}
```

#### 描述和关键词更新

```json
{
  "description": "基于 Go 的 DDD 模板应用文档 - VitePress 2.0",
  "keywords": [
    "vitepress",
    "vitepress-2.0", // 新增
    "documentation",
    "go",
    "ddd"
  ]
}
```

## VitePress 2.0 新特性

### 1. 性能优化

- **Git 时间戳批量获取**: 单次 git 调用获取所有文件时间戳 (原来是每个文件单独调用)
- **改进的 HMR**: 主题和配置文件的热更新更快速

### 2. Markdown 增强

- **CJK 友好**: 默认启用 `markdown-it-cjk-friendly` 插件
- **图片懒加载**: 支持原生的 `lazyLoading` 配置
- **代码块增强**: 支持自定义 display-name 和 Shell 符号保护

### 3. API 变更

- **配置重命名**: `cjkFriendly` → `cjkFriendlyEmphasis`
- **代码块属性**: 禁用 `markdown-it-attrs`，改用 Shiki transformers
- **PostCSS 样式隔离**: `postcssIsolateStyles` 默认值变更

### 4. 依赖更新

- **Vite 7**: 使用最新的 Vite 构建工具
- **ESM Only**: 仅支持 ESM 模块系统
- **DocSearch v4 Beta**: 更现代的搜索体验 (如果使用 Algolia)

## 破坏性变更

### 1. Node.js 版本要求

- ❌ **不再支持**: Node.js 18.x
- ✅ **最低要求**: Node.js 20.19.0 或 22.12.0+

### 2. 配置重命名

```typescript
// ❌ 旧配置 (不再支持)
markdown: {
  cjkFriendly: true;
}

// ✅ 新配置
markdown: {
  cjkFriendlyEmphasis: true;
}
```

### 3. PostCSS 配置

- `postcssIsolateStyles` 的 `transform` 和 `exclude` 选项不再支持

## 测试结果

### 构建测试

```bash
npm run docs:build
```

**结果**: ✅ 成功 (2.98 秒)

### 构建输出

```
vitepress v2.0.0-alpha.13
build complete in 2.98s.
✓ building client + server bundles...
✓ rendering pages...
```

## 后续建议

### 1. 开发服务器测试

```bash
npm run docs:dev
```

测试热更新和开发体验

### 2. 预览构建结果

```bash
npm run docs:preview
```

测试生产构建的输出

### 3. 性能对比

- 对比 1.6.4 和 2.0.0-alpha.13 的构建速度
- 测试 Git 时间戳获取的性能提升
- 验证图片懒加载的效果

### 4. 探索新特性

- 尝试使用 Shiki transformers 自定义代码块
- 测试改进的 CJK 语言支持
- 体验更快的 HMR

### 5. 监控稳定性

- 2.0.0-alpha.13 是测试版本
- 关注 GitHub Issues: https://github.com/vuejs/vitepress/issues
- 关注稳定版发布: https://github.com/vuejs/vitepress/releases

## 回滚方案

如果遇到问题，可以回滚到 1.6.4：

```bash
# 1. 回滚 VitePress 版本
npm install -D vitepress@^1.6.4

# 2. 恢复 package.json
git checkout package.json

# 3. 恢复配置文件
git checkout docs/.vitepress/config.ts

# 4. 重新安装依赖
npm install
```

## 相关资源

- VitePress 官方文档: https://vitepress.dev/
- VitePress 2.0 Changelog: https://github.com/vuejs/vitepress/blob/main/CHANGELOG.md
- VitePress GitHub: https://github.com/vuejs/vitepress
- Vite 7 文档: https://vite.dev/

## 升级清单

- [x] 检查 Node.js 版本 (>= 20.19.0)
- [x] 升级 VitePress 到 2.0.0-alpha.13
- [x] 更新配置文件 (cjkFriendlyEmphasis、lazyLoading)
- [x] 更新 package.json (engines、version、description)
- [x] 测试构建 (docs:build)
- [ ] 测试开发服务器 (docs:dev)
- [ ] 测试预览服务器 (docs:preview)
- [ ] 性能测试和对比
- [ ] 文档部署测试

## 升级总结

VitePress 2.0.0-alpha.13 升级顺利完成！主要变更包括：

1. ✅ Vite 7 构建工具链
2. ✅ 改进的 CJK 语言支持
3. ✅ 图片懒加载功能
4. ✅ 更快的 Git 时间戳获取
5. ✅ 更好的开发体验

**构建时间**: 2.98 秒
**依赖包变更**: +6 packages, -24 packages, changed 24 packages
**漏洞**: 0 vulnerabilities ✅

升级后的项目已准备好使用 VitePress 2.0 的新特性！
