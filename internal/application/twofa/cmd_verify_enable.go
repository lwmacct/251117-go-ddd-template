package twofa

// VerifyEnableCommand 验证并启用 2FA 命令
type VerifyEnableCommand struct {
	UserID uint
	Code   string
}
