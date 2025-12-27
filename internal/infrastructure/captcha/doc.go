// Package captcha 提供图形验证码的基础设施实现。
//
// 本包使用 [github.com/mojocn/base64Captcha] 库生成验证码图片，特点：
//   - 支持字符型验证码（a-z, 0-9）
//   - 输出 Base64 编码的 PNG 图片，便于前端直接展示
//   - 可配置的干扰线和噪点，平衡安全性与可读性
//
// # 组件职责
//
//   - [Service]: 验证码图片生成服务
//   - [NewRepository]: 验证码存储仓储（内存实现）
//
// # 使用示例
//
//	svc := captcha.NewService()
//	repo := captcha.NewRepository()
//
//	// 生成验证码
//	captchaID, imageBase64, code, _ := svc.GenerateRandomCode()
//
//	// 存储验证码
//	repo.Create(ctx, captchaID, code, 5*time.Minute)
//
//	// 验证
//	valid, _ := repo.Verify(ctx, captchaID, userInput)
package captcha
