//nolint:dupl // 与 delete_token_handler_test.go 测试结构相似，但测试不同 Handler
package pat

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainPAT "github.com/lwmacct/251117-go-ddd-template/internal/domain/pat"
)

func TestDisableTokenHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockPATCmdRepo := new(MockPATCommandRepository)
	mockPATQryRepo := new(MockPATQueryRepository)

	token := &domainPAT.PersonalAccessToken{
		ID:     1,
		UserID: 1,
		Name:   "Test Token",
		Status: "active",
	}

	mockPATQryRepo.On("FindByID", mock.Anything, uint(1)).Return(token, nil)
	mockPATCmdRepo.On("Disable", mock.Anything, uint(1)).Return(nil)

	handler := NewDisableTokenHandler(mockPATCmdRepo, mockPATQryRepo)

	// Act
	err := handler.Handle(context.Background(), DisableTokenCommand{
		UserID:  1,
		TokenID: 1,
	})

	// Assert
	require.NoError(t, err)

	mockPATQryRepo.AssertExpectations(t)
	mockPATCmdRepo.AssertExpectations(t)
}

func TestDisableTokenHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        DisableTokenCommand
		setupMocks func(*MockPATCommandRepository, *MockPATQueryRepository)
		wantErr    string
	}{
		{
			name: "Token 不存在",
			cmd:  DisableTokenCommand{UserID: 1, TokenID: 999},
			setupMocks: func(patCmd *MockPATCommandRepository, patQry *MockPATQueryRepository) {
				patQry.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: "token not found",
		},
		{
			name: "Token 不属于该用户",
			cmd:  DisableTokenCommand{UserID: 1, TokenID: 2},
			setupMocks: func(patCmd *MockPATCommandRepository, patQry *MockPATQueryRepository) {
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
		{
			name: "禁用失败",
			cmd:  DisableTokenCommand{UserID: 1, TokenID: 1},
			setupMocks: func(patCmd *MockPATCommandRepository, patQry *MockPATQueryRepository) {
				token := &domainPAT.PersonalAccessToken{
					ID:     1,
					UserID: 1,
					Name:   "Test Token",
					Status: "active",
				}
				patQry.On("FindByID", mock.Anything, uint(1)).Return(token, nil)
				patCmd.On("Disable", mock.Anything, uint(1)).Return(errors.New("db error"))
			},
			wantErr: "failed to disable token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPATCmdRepo := new(MockPATCommandRepository)
			mockPATQryRepo := new(MockPATQueryRepository)

			tt.setupMocks(mockPATCmdRepo, mockPATQryRepo)

			handler := NewDisableTokenHandler(mockPATCmdRepo, mockPATQryRepo)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
