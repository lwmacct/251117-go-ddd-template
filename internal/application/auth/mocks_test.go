//nolint:forcetypeassert // Mock 返回值类型在测试中总是已知的
package auth

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	domainAuth "github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	domainTwoFA "github.com/lwmacct/251117-go-ddd-template/internal/domain/twofa"
	domainUser "github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// ============================================================
// MockUserQueryRepository
// ============================================================

type MockUserQueryRepository struct {
	mock.Mock
}

func (m *MockUserQueryRepository) GetByID(ctx context.Context, id uint) (*domainUser.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserQueryRepository) GetByIDWithRoles(ctx context.Context, id uint) (*domainUser.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserQueryRepository) GetByUsername(ctx context.Context, username string) (*domainUser.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserQueryRepository) GetByUsernameWithRoles(ctx context.Context, username string) (*domainUser.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserQueryRepository) GetByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserQueryRepository) GetByEmailWithRoles(ctx context.Context, email string) (*domainUser.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserQueryRepository) List(ctx context.Context, offset, limit int) ([]*domainUser.User, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainUser.User), args.Error(1)
}

func (m *MockUserQueryRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserQueryRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*domainUser.User, error) {
	args := m.Called(ctx, keyword, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainUser.User), args.Error(1)
}

func (m *MockUserQueryRepository) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	args := m.Called(ctx, keyword)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserQueryRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserQueryRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserQueryRepository) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}

func (m *MockUserQueryRepository) GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}

// ============================================================
// MockUserCommandRepository
// ============================================================

type MockUserCommandRepository struct {
	mock.Mock
}

func (m *MockUserCommandRepository) Create(ctx context.Context, u *domainUser.User) error {
	args := m.Called(ctx, u)
	// 模拟数据库分配 ID
	if u.ID == 0 {
		u.ID = 1
	}
	return args.Error(0)
}

func (m *MockUserCommandRepository) Update(ctx context.Context, u *domainUser.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserCommandRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserCommandRepository) UpdatePassword(ctx context.Context, id uint, hashedPassword string) error {
	args := m.Called(ctx, id, hashedPassword)
	return args.Error(0)
}

func (m *MockUserCommandRepository) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	args := m.Called(ctx, userID, roleIDs)
	return args.Error(0)
}

func (m *MockUserCommandRepository) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	args := m.Called(ctx, userID, roleIDs)
	return args.Error(0)
}

func (m *MockUserCommandRepository) UpdateStatus(ctx context.Context, userID uint, status string) error {
	args := m.Called(ctx, userID, status)
	return args.Error(0)
}

// ============================================================
// MockAuthService
// ============================================================

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) error {
	args := m.Called(ctx, hashedPassword, plainPassword)
	return args.Error(0)
}

func (m *MockAuthService) GeneratePasswordHash(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidatePasswordPolicy(ctx context.Context, password string) error {
	args := m.Called(ctx, password)
	return args.Error(0)
}

func (m *MockAuthService) GenerateAccessToken(ctx context.Context, userID uint, username string) (string, time.Time, error) {
	args := m.Called(ctx, userID, username)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockAuthService) GenerateRefreshToken(ctx context.Context, userID uint) (string, time.Time, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockAuthService) ValidateAccessToken(ctx context.Context, token string) (*domainAuth.TokenClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainAuth.TokenClaims), args.Error(1)
}

func (m *MockAuthService) ValidateRefreshToken(ctx context.Context, token string) (uint, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockAuthService) GeneratePATToken(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) HashPATToken(ctx context.Context, token string) string {
	args := m.Called(ctx, token)
	return args.String(0)
}

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
// MockTwoFAQueryRepository
// ============================================================

type MockTwoFAQueryRepository struct {
	mock.Mock
}

func (m *MockTwoFAQueryRepository) FindByUserID(ctx context.Context, userID uint) (*domainTwoFA.TwoFA, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainTwoFA.TwoFA), args.Error(1)
}

func (m *MockTwoFAQueryRepository) IsEnabled(ctx context.Context, userID uint) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

// ============================================================
// MockLoginSessionService (测试用接口实现)
// ============================================================

type MockLoginSessionService struct {
	mock.Mock
}

func (m *MockLoginSessionService) GenerateSessionToken(ctx context.Context, userID uint, account string) (string, error) {
	args := m.Called(ctx, userID, account)
	return args.String(0), args.Error(1)
}

func (m *MockLoginSessionService) VerifySessionToken(ctx context.Context, token string) (*LoginSessionData, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoginSessionData), args.Error(1)
}

// LoginSessionData 测试用登录会话数据
type LoginSessionData struct {
	UserID  uint
	Account string
}

// ============================================================
// MockTwoFAService (测试用接口实现)
// ============================================================

type MockTwoFAService struct {
	mock.Mock
}

func (m *MockTwoFAService) Verify(ctx context.Context, userID uint, code string) (bool, error) {
	args := m.Called(ctx, userID, code)
	return args.Bool(0), args.Error(1)
}
