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

func TestListTokensHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockPATQryRepo := new(MockPATQueryRepository)

	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour)

	tokens := []*domainPAT.PersonalAccessToken{
		{
			ID:          1,
			UserID:      1,
			Name:        "Token 1",
			TokenPrefix: "pat_ABC12",
			Permissions: []string{"user:read"},
			Status:      "active",
			ExpiresAt:   &expiresAt,
			CreatedAt:   now.Add(-7 * 24 * time.Hour),
			UpdatedAt:   now.Add(-1 * 24 * time.Hour),
		},
		{
			ID:          2,
			UserID:      1,
			Name:        "Token 2",
			TokenPrefix: "pat_XYZ99",
			Permissions: []string{"user:read", "user:write"},
			Status:      "disabled",
			ExpiresAt:   nil, // 永久 token
			CreatedAt:   now.Add(-14 * 24 * time.Hour),
			UpdatedAt:   now.Add(-2 * 24 * time.Hour),
		},
	}

	mockPATQryRepo.On("ListByUser", mock.Anything, uint(1)).Return(tokens, nil)

	handler := NewListTokensHandler(mockPATQryRepo)

	// Act
	result, err := handler.Handle(context.Background(), ListTokensQuery{
		UserID: 1,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	// 验证第一个 token
	assert.Equal(t, uint(1), result[0].ID)
	assert.Equal(t, "Token 1", result[0].Name)
	assert.Equal(t, "pat_ABC12", result[0].TokenPrefix)
	assert.Equal(t, []string{"user:read"}, result[0].Permissions)
	assert.Equal(t, "active", result[0].Status)
	assert.NotNil(t, result[0].ExpiresAt)

	// 验证第二个 token
	assert.Equal(t, uint(2), result[1].ID)
	assert.Equal(t, "Token 2", result[1].Name)
	assert.Equal(t, "pat_XYZ99", result[1].TokenPrefix)
	assert.Equal(t, []string{"user:read", "user:write"}, result[1].Permissions)
	assert.Equal(t, "disabled", result[1].Status)
	assert.Nil(t, result[1].ExpiresAt) // 永久 token

	mockPATQryRepo.AssertExpectations(t)
}

func TestListTokensHandler_Handle_EmptyList(t *testing.T) {
	// Arrange
	mockPATQryRepo := new(MockPATQueryRepository)

	mockPATQryRepo.On("ListByUser", mock.Anything, uint(1)).Return([]*domainPAT.PersonalAccessToken{}, nil)

	handler := NewListTokensHandler(mockPATQryRepo)

	// Act
	result, err := handler.Handle(context.Background(), ListTokensQuery{
		UserID: 1,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)

	mockPATQryRepo.AssertExpectations(t)
}

func TestListTokensHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		query      ListTokensQuery
		setupMocks func(*MockPATQueryRepository)
		wantErr    string
	}{
		{
			name:  "查询失败",
			query: ListTokensQuery{UserID: 1},
			setupMocks: func(patQry *MockPATQueryRepository) {
				patQry.On("ListByUser", mock.Anything, uint(1)).Return(nil, errors.New("db connection error"))
			},
			wantErr: "failed to fetch tokens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPATQryRepo := new(MockPATQueryRepository)

			tt.setupMocks(mockPATQryRepo)

			handler := NewListTokensHandler(mockPATQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
