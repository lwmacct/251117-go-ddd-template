//nolint:forcetypeassert // 测试中的类型断言是可控的
package auditlog

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainAuditLog "github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

func TestCreateLogHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockAuditLogCommandRepository)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*auditlog.AuditLog")).Return(nil)

	handler := NewCreateLogHandler(mockRepo)

	// Act
	err := handler.Handle(context.Background(), CreateLogCommand{
		UserID:     1,
		Username:   "admin",
		Action:     "login",
		Resource:   "auth",
		ResourceID: "",
		IPAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0",
		Details:    `{"event":"login_success"}`,
		Status:     "success",
	})

	// Assert
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// 验证传递给 Repository 的日志对象
	call := mockRepo.Calls[0]
	log := call.Arguments.Get(1).(*domainAuditLog.AuditLog)
	assert.Equal(t, uint(1), log.UserID)
	assert.Equal(t, "admin", log.Username)
	assert.Equal(t, "login", log.Action)
	assert.Equal(t, "auth", log.Resource)
	assert.Equal(t, "success", log.Status)
}

func TestCreateLogHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        CreateLogCommand
		setupMocks func(*MockAuditLogCommandRepository)
		wantErr    string
	}{
		{
			name: "存储失败",
			cmd: CreateLogCommand{
				UserID:   1,
				Username: "user",
				Action:   "create",
				Resource: "user",
				Status:   "success",
			},
			setupMocks: func(repo *MockAuditLogCommandRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*auditlog.AuditLog")).Return(errors.New("database error"))
			},
			wantErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockAuditLogCommandRepository)
			tt.setupMocks(mockRepo)

			handler := NewCreateLogHandler(mockRepo)

			// Act
			err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestCreateLogHandler_Handle_AllFields(t *testing.T) {
	// 测试所有字段都正确传递
	mockRepo := new(MockAuditLogCommandRepository)

	var capturedLog *domainAuditLog.AuditLog
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*auditlog.AuditLog")).Run(func(args mock.Arguments) {
		capturedLog = args.Get(1).(*domainAuditLog.AuditLog)
	}).Return(nil)

	handler := NewCreateLogHandler(mockRepo)

	cmd := CreateLogCommand{
		UserID:     123,
		Username:   "testuser",
		Action:     "delete",
		Resource:   "user",
		ResourceID: "456",
		IPAddress:  "10.0.0.1",
		UserAgent:  "CustomAgent/2.0",
		Details:    `{"reason":"admin_request"}`,
		Status:     "success",
	}

	err := handler.Handle(context.Background(), cmd)

	require.NoError(t, err)
	assert.NotNil(t, capturedLog)
	assert.Equal(t, cmd.UserID, capturedLog.UserID)
	assert.Equal(t, cmd.Username, capturedLog.Username)
	assert.Equal(t, cmd.Action, capturedLog.Action)
	assert.Equal(t, cmd.Resource, capturedLog.Resource)
	assert.Equal(t, cmd.ResourceID, capturedLog.ResourceID)
	assert.Equal(t, cmd.IPAddress, capturedLog.IPAddress)
	assert.Equal(t, cmd.UserAgent, capturedLog.UserAgent)
	assert.Equal(t, cmd.Details, capturedLog.Details)
	assert.Equal(t, cmd.Status, capturedLog.Status)
}
