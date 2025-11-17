# 部署 VitePress 文档到 GitHub Pages

这是一个快速部署检查清单。详细说明请查看[部署指南](/guide/deployment)。

## 快速检查清单

- [x] 创建 `.github/workflows/deploy-docs.yml` workflow 文件
- [x] 配置 VitePress `base: '/251117-bd-vmalert/'`
- [ ] 在 GitHub 仓库中启用 GitHub Pages（Settings → Pages → Source: GitHub Actions）
- [ ] 提交并推送代码到 `main` 分支
- [ ] 查看 Actions 标签页确认部署成功
- [ ] 访问 https://lwmacct.github.io/251117-bd-vmalert/

## 启用 GitHub Pages

1. 访问仓库：https://github.com/lwmacct/251117-bd-vmalert
2. 进入 **Settings** → **Pages**
3. 在 **Source** 下选择 **GitHub Actions**
4. 保存设置

## 首次部署命令

```bash
# 提交并推送所有文件
git add .
git commit -m "Add VitePress documentation with GitHub Pages deployment"
git push origin main
```

## 查看部署状态

访问：https://github.com/lwmacct/251117-bd-vmalert/actions

## 文档地址

部署成功后访问：https://lwmacct.github.io/251117-bd-vmalert/

---

完整部署文档请参考：[部署指南](/guide/deployment)
