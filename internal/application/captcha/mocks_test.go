//nolint:forcetypeassert // Mock 返回值类型在测试中总是已知的
package captcha

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// ============================================================
// MockCaptchaCommandRepository
// ============================================================

type MockCaptchaCommandRepository struct {
	mock.Mock
}

func (m *MockCaptchaCommandRepository) Create(ctx context.Context, captchaID string, code string, expiration time.Duration) error {
	args := m.Called(ctx, captchaID, code, expiration)
	return args.Error(0)
}

func (m *MockCaptchaCommandRepository) Verify(ctx context.Context, captchaID string, code string) (bool, error) {
	args := m.Called(ctx, captchaID, code)
	return args.Bool(0), args.Error(1)
}

func (m *MockCaptchaCommandRepository) Delete(ctx context.Context, captchaID string) error {
	args := m.Called(ctx, captchaID)
	return args.Error(0)
}

// ============================================================
// MockCaptchaService
// ============================================================

type MockCaptchaService struct {
	mock.Mock
}

func (m *MockCaptchaService) GenerateRandomCode() (string, string, string, error) {
	args := m.Called()
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

func (m *MockCaptchaService) GenerateCustomCodeImage(text string) (string, error) {
	args := m.Called(text)
	return args.String(0), args.Error(1)
}

func (m *MockCaptchaService) GenerateDevCaptchaID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockCaptchaService) GetDefaultExpiration() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}
