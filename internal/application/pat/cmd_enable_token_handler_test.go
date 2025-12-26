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

func TestEnableTokenHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockPATCmdRepo := new(MockPATCommandRepository)
	mockPATQryRepo := new(MockPATQueryRepository)

	futureTime := time.Now().Add(30 * 24 * time.Hour)
	token := &domainPAT.PersonalAccessToken{
		ID:        1,
		UserID:    1,
		Name:      "Test Token",
		Status:    "disabled",
		ExpiresAt: &futureTime, // 未过期
	}

	mockPATQryRepo.On("FindByID", mock.Anything, uint(1)).Return(token, nil)
	mockPATCmdRepo.On("Enable", mock.Anything, uint(1)).Return(nil)

	handler := NewEnableTokenHandler(mockPATCmdRepo, mockPATQryRepo)

	// Act
	err := handler.Handle(context.Background(), EnableTokenCommand{
		UserID:  1,
		TokenID: 1,
	})

	// Assert
	require.NoError(t, err)

	mockPATQryRepo.AssertExpectations(t)
	mockPATCmdRepo.AssertExpectations(t)
}

func TestEnableTokenHandler_Handle_Success_NoExpiration(t *testing.T) {
	// Arrange
	mockPATCmdRepo := new(MockPATCommandRepository)
	mockPATQryRepo := new(MockPATQueryRepository)

	token := &domainPAT.PersonalAccessToken{
		ID:        1,
		UserID:    1,
		Name:      "Test Token",
		Status:    "disabled",
		ExpiresAt: nil, // 永久 token
	}

	mockPATQryRepo.On("FindByID", mock.Anything, uint(1)).Return(token, nil)
	mockPATCmdRepo.On("Enable", mock.Anything, uint(1)).Return(nil)

	handler := NewEnableTokenHandler(mockPATCmdRepo, mockPATQryRepo)

	// Act
	err := handler.Handle(context.Background(), EnableTokenCommand{
		UserID:  1,
		TokenID: 1,
	})

	// Assert
	require.NoError(t, err)

	mockPATQryRepo.AssertExpectations(t)
	mockPATCmdRepo.AssertExpectations(t)
}

func TestEnableTokenHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        EnableTokenCommand
		setupMocks func(*MockPATCommandRepository, *MockPATQueryRepository)
		wantErr    string
	}{
		{
			name: "Token 不存在",
			cmd:  EnableTokenCommand{UserID: 1, TokenID: 999},
			setupMocks: func(patCmd *MockPATCommandRepository, patQry *MockPATQueryRepository) {
				patQry.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: "token not found",
		},
		{
			name: "Token 不属于该用户",
			cmd:  EnableTokenCommand{UserID: 1, TokenID: 2},
			setupMocks: func(patCmd *MockPATCommandRepository, patQry *MockPATQueryRepository) {
				token := &domainPAT.PersonalAccessToken{
					ID:     2,
					UserID: 999, // 不同的用户
					Name:   "Other User Token",
					Status: "disabled",
				}
				patQry.On("FindByID", mock.Anything, uint(2)).Return(token, nil)
			},
			wantErr: "token does not belong to this user",
		},
		{
			name: "Token 已过期，无法启用",
			cmd:  EnableTokenCommand{UserID: 1, TokenID: 1},
			setupMocks: func(patCmd *MockPATCommandRepository, patQry *MockPATQueryRepository) {
				pastTime := time.Now().Add(-24 * time.Hour) // 昨天过期
				token := &domainPAT.PersonalAccessToken{
					ID:        1,
					UserID:    1,
					Name:      "Expired Token",
					Status:    "disabled",
					ExpiresAt: &pastTime,
				}
				patQry.On("FindByID", mock.Anything, uint(1)).Return(token, nil)
			},
			wantErr: "token is expired and cannot be enabled",
		},
		{
			name: "启用失败",
			cmd:  EnableTokenCommand{UserID: 1, TokenID: 1},
			setupMocks: func(patCmd *MockPATCommandRepository, patQry *MockPATQueryRepository) {
				futureTime := time.Now().Add(30 * 24 * time.Hour)
				token := &domainPAT.PersonalAccessToken{
					ID:        1,
					UserID:    1,
					Name:      "Test Token",
					Status:    "disabled",
					ExpiresAt: &futureTime,
				}
				patQry.On("FindByID", mock.Anything, uint(1)).Return(token, nil)
				patCmd.On("Enable", mock.Anything, uint(1)).Return(errors.New("db error"))
			},
			wantErr: "failed to enable token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPATCmdRepo := new(MockPATCommandRepository)
			mockPATQryRepo := new(MockPATQueryRepository)

			tt.setupMocks(mockPATCmdRepo, mockPATQryRepo)

			handler := NewEnableTokenHandler(mockPATCmdRepo, mockPATQryRepo)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
