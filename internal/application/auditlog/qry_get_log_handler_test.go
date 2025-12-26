package auditlog

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetLogHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockAuditLogQueryRepository)

	expectedLog := newTestAuditLog(1)
	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(expectedLog, nil)

	handler := NewGetLogHandler(mockRepo)

	// Act
	result, err := handler.Handle(context.Background(), GetLogQuery{LogID: 1})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedLog.ID, result.ID)
	assert.Equal(t, expectedLog.UserID, result.UserID)
	assert.Equal(t, expectedLog.Action, result.Action)
	assert.Equal(t, expectedLog.Resource, result.Resource)
	assert.Equal(t, expectedLog.Status, result.Status)

	mockRepo.AssertExpectations(t)
}

func TestGetLogHandler_Handle_NotFound(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(*MockAuditLogQueryRepository)
	}{
		{
			name: "返回nil日志",
			setupMocks: func(repo *MockAuditLogQueryRepository) {
				repo.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)
			},
		},
		{
			name: "返回错误",
			setupMocks: func(repo *MockAuditLogQueryRepository) {
				repo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockAuditLogQueryRepository)
			tt.setupMocks(mockRepo)

			handler := NewGetLogHandler(mockRepo)

			// Act
			result, err := handler.Handle(context.Background(), GetLogQuery{LogID: 999})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "audit log not found")
		})
	}
}
