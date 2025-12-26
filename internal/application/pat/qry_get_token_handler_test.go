package pat

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainPAT "github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

func TestGetTokenHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockPATQryRepo := new(MockPATQueryRepository)

	createdAt := time.Now().Add(-7 * 24 * time.Hour)
	updatedAt := time.Now().Add(-1 * 24 * time.Hour)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	token := &domainPAT.PersonalAccessToken{
		ID:          1,
		UserID:      1,
		Name:        "Test Token",
		Token:       "hashed_token",
		TokenPrefix: "pat_ABC12",
		Permissions: []string{"user:read", "user:write"},
		Status:      "active",
		ExpiresAt:   &expiresAt,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		IPWhitelist: []string{"127.0.0.1"},
	}

	mockPATQryRepo.On("FindByID", mock.Anything, uint(1)).Return(token, nil)

	handler := NewGetTokenHandler(mockPATQryRepo)

	// Act
	result, err := handler.Handle(context.Background(), GetTokenQuery{
		UserID:  1,
		TokenID: 1,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, "Test Token", result.Name)
	assert.Equal(t, "pat_ABC12", result.TokenPrefix)
	assert.Equal(t, []string{"user:read", "user:write"}, result.Permissions)
	assert.Equal(t, "active", result.Status)
	assert.NotNil(t, result.ExpiresAt)
	assert.Equal(t, []string{"127.0.0.1"}, result.IPWhitelist)

	mockPATQryRepo.AssertExpectations(t)
}

func TestGetTokenHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		query      GetTokenQuery
		setupMocks func(*MockPATQueryRepository)
		wantErr    string
	}{
		{
			name:  "Token 不存在",
			query: GetTokenQuery{UserID: 1, TokenID: 999},
			setupMocks: func(patQry *MockPATQueryRepository) {
				patQry.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: "token not found",
		},
		{
			name:  "Token 不属于该用户",
			query: GetTokenQuery{UserID: 1, TokenID: 2},
			setupMocks: func(patQry *MockPATQueryRepository) {
				token := &domainPAT.PersonalAccessToken{
					ID:     2,
					UserID: 999, // 不同的用户
					Name:   "Other User Token",
					Status: "active",
				}
				patQry.On("FindByID", mock.Anything, uint(2)).Return(token, nil)
			},
			wantErr: "token does not belong to this user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPATQryRepo := new(MockPATQueryRepository)

			tt.setupMocks(mockPATQryRepo)

			handler := NewGetTokenHandler(mockPATQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
