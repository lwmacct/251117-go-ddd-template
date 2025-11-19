package twofa

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// Service 2FA 服务
type Service struct {
	twofaCommandRepo twofa.CommandRepository
	twofaQueryRepo   twofa.QueryRepository
	userQueryRepo    user.QueryRepository
	issuer           string // TOTP 发行者名称
}

// NewService 创建 2FA 服务
func NewService(twofaCommandRepo twofa.CommandRepository, twofaQueryRepo twofa.QueryRepository, userQueryRepo user.QueryRepository, issuer string) *Service {
	if issuer == "" {
		issuer = "Go-DDD-Template"
	}
	return &Service{
		twofaCommandRepo: twofaCommandRepo,
		twofaQueryRepo:   twofaQueryRepo,
		userQueryRepo:    userQueryRepo,
		issuer:           issuer,
	}
}

// SetupResponse 设置 2FA 返回数据
type SetupResponse struct {
	Secret    string `json:"secret"`      // TOTP 密钥（用户可手动输入）
	QRCodeURL string `json:"qrcode_url"`  // 二维码URL
	QRCodeImg string `json:"qrcode_img"`  // Base64 编码的二维码图片
}

// Setup 设置 2FA（生成密钥和二维码）
func (s *Service) Setup(ctx context.Context, userID uint) (*SetupResponse, error) {
	// 查找用户
	u, err := s.userQueryRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 生成 TOTP 密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: u.Username, // 使用用户名
		SecretSize:  20,         // 20 字节（160 位），Base32 编码后 32 个字符
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// 构建简化的二维码 URL（移除默认参数）
	secret := key.Secret()
	qrcodeURL := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		url.QueryEscape(s.issuer),
		url.QueryEscape(u.Username),
		secret,
		url.QueryEscape(s.issuer),
	)

	// 生成二维码图片
	qrCodeBytes, err := qrcode.Encode(qrcodeURL, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Base64 编码
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeBytes)

	// 存储密钥到数据库（未启用状态）
	tfa := &twofa.TwoFA{
		UserID:        userID,
		Enabled:       false,
		Secret:        secret,
		RecoveryCodes: twofa.RecoveryCodes{}, // 空恢复码，验证后生成
	}

	if err := s.twofaCommandRepo.CreateOrUpdate(ctx, tfa); err != nil {
		return nil, fmt.Errorf("failed to save 2FA secret: %w", err)
	}

	return &SetupResponse{
		Secret:    secret,
		QRCodeURL: qrcodeURL,
		QRCodeImg: "data:image/png;base64," + qrCodeBase64,
	}, nil
}

// VerifyAndEnable 验证 TOTP 代码并启用 2FA
// 返回恢复码列表
func (s *Service) VerifyAndEnable(ctx context.Context, userID uint, code string) ([]string, error) {
	// 查找 2FA 配置
	tfa, err := s.twofaQueryRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get 2FA config: %w", err)
	}

	if tfa == nil || tfa.Secret == "" {
		return nil, fmt.Errorf("please setup 2FA first")
	}

	// 验证 TOTP 代码
	valid := totp.Validate(code, tfa.Secret)
	if !valid {
		return nil, fmt.Errorf("invalid verification code")
	}

	// 生成恢复码（8 个）
	recoveryCodes, err := GenerateRecoveryCodes(8)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recovery codes: %w", err)
	}

	// 启用 2FA 并保存恢复码
	now := time.Now()
	tfa.Enabled = true
	tfa.RecoveryCodes = recoveryCodes
	tfa.SetupCompletedAt = &now

	if err := s.twofaCommandRepo.CreateOrUpdate(ctx, tfa); err != nil {
		return nil, fmt.Errorf("failed to enable 2FA: %w", err)
	}

	return recoveryCodes, nil
}

// Verify 验证 TOTP 代码或恢复码
func (s *Service) Verify(ctx context.Context, userID uint, code string) (bool, error) {
	// 查找 2FA 配置
	tfa, err := s.twofaQueryRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get 2FA config: %w", err)
	}

	if tfa == nil || !tfa.Enabled {
		return false, fmt.Errorf("2FA not enabled")
	}

	// 首先尝试 TOTP 验证
	if totp.Validate(code, tfa.Secret) {
		// 更新最后使用时间
		now := time.Now()
		tfa.LastUsedAt = &now
		_ = s.twofaCommandRepo.CreateOrUpdate(ctx, tfa)
		return true, nil
	}

	// TOTP 验证失败，尝试恢复码
	code = strings.TrimSpace(code)
	for i, recoveryCode := range tfa.RecoveryCodes {
		if recoveryCode == code {
			// 移除已使用的恢复码
			tfa.RecoveryCodes = append(tfa.RecoveryCodes[:i], tfa.RecoveryCodes[i+1:]...)
			now := time.Now()
			tfa.LastUsedAt = &now
			if err := s.twofaCommandRepo.CreateOrUpdate(ctx, tfa); err != nil {
				return false, fmt.Errorf("failed to update recovery codes: %w", err)
			}
			return true, nil
		}
	}

	return false, nil
}

// Disable 禁用 2FA
func (s *Service) Disable(ctx context.Context, userID uint) error {
	return s.twofaCommandRepo.Delete(ctx, userID)
}

// GetStatus 获取 2FA 状态
func (s *Service) GetStatus(ctx context.Context, userID uint) (bool, int, error) {
	tfa, err := s.twofaQueryRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, 0, fmt.Errorf("failed to get 2FA config: %w", err)
	}

	if tfa == nil {
		return false, 0, nil
	}

	return tfa.Enabled, len(tfa.RecoveryCodes), nil
}

// GenerateRecoveryCodes 生成恢复码
// 格式：xxxx-xxxx（8位数字，用连字符分隔）
func GenerateRecoveryCodes(count int) ([]string, error) {
	codes := make([]string, count)

	for i := 0; i < count; i++ {
		// 生成 8 位随机数字
		b := make([]byte, 4)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}

		// 转换为 8 位数字
		num := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
		num = num % 100000000 // 限制在 8 位数字

		// 格式化为 xxxx-xxxx
		codes[i] = fmt.Sprintf("%04d-%04d", num/10000, num%10000)
	}

	return codes, nil
}
