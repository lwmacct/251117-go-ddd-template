//nolint:forcetypeassert,nonamedreturns // Mock 返回值类型在测试中总是已知的
package pat

import (
	"context"

	"github.com/stretchr/testify/mock"

	domainPAT "github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
	domainUser "github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

// ============================================================
// MockPATCommandRepository
// ============================================================

type MockPATCommandRepository struct {
	mock.Mock
}

func (m *MockPATCommandRepository) Create(ctx context.Context, pat *domainPAT.PersonalAccessToken) error {
	args := m.Called(ctx, pat)
	// 模拟数据库分配 ID
	if pat.ID == 0 {
		pat.ID = 1
	}
	return args.Error(0)
}

func (m *MockPATCommandRepository) Update(ctx context.Context, pat *domainPAT.PersonalAccessToken) error {
	args := m.Called(ctx, pat)
	return args.Error(0)
}

func (m *MockPATCommandRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPATCommandRepository) Disable(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPATCommandRepository) Enable(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPATCommandRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPATCommandRepository) CleanupExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ============================================================
// MockPATQueryRepository
// ============================================================

type MockPATQueryRepository struct {
	mock.Mock
}

func (m *MockPATQueryRepository) FindByToken(ctx context.Context, tokenHash string) (*domainPAT.PersonalAccessToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainPAT.PersonalAccessToken), args.Error(1)
}

func (m *MockPATQueryRepository) FindByID(ctx context.Context, id uint) (*domainPAT.PersonalAccessToken, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainPAT.PersonalAccessToken), args.Error(1)
}

func (m *MockPATQueryRepository) FindByPrefix(ctx context.Context, prefix string) (*domainPAT.PersonalAccessToken, error) {
	args := m.Called(ctx, prefix)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainPAT.PersonalAccessToken), args.Error(1)
}

func (m *MockPATQueryRepository) ListByUser(ctx context.Context, userID uint) ([]*domainPAT.PersonalAccessToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainPAT.PersonalAccessToken), args.Error(1)
}

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
// MockTokenGenerator
// ============================================================

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GeneratePAT() (plainToken, hashedToken, prefix string, err error) {
	args := m.Called()
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

func (m *MockTokenGenerator) HashToken(plainToken string) string {
	args := m.Called(plainToken)
	return args.String(0)
}

func (m *MockTokenGenerator) ValidateTokenFormat(token string) bool {
	args := m.Called(token)
	return args.Bool(0)
}

func (m *MockTokenGenerator) ExtractPrefix(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}
