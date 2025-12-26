//nolint:forcetypeassert // Mock 返回值类型在测试中总是已知的
package twofa

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
)

// MockTwoFAService 2FA 服务 Mock
type MockTwoFAService struct {
	mock.Mock
}

// Setup Mock 实现
func (m *MockTwoFAService) Setup(ctx context.Context, userID uint) (*twofa.SetupResult, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*twofa.SetupResult), args.Error(1)
}

// VerifyAndEnable Mock 实现
func (m *MockTwoFAService) VerifyAndEnable(ctx context.Context, userID uint, code string) ([]string, error) {
	args := m.Called(ctx, userID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// Verify Mock 实现
func (m *MockTwoFAService) Verify(ctx context.Context, userID uint, code string) (bool, error) {
	args := m.Called(ctx, userID, code)
	return args.Bool(0), args.Error(1)
}

// Disable Mock 实现
func (m *MockTwoFAService) Disable(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// GetStatus Mock 实现
func (m *MockTwoFAService) GetStatus(ctx context.Context, userID uint) (bool, int, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Int(1), args.Error(2)
}
