//nolint:forcetypeassert // 测试中的类型断言是可控的
package auditlog

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainAuditLog "github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

func TestListLogsHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockAuditLogQueryRepository)

	logs := []domainAuditLog.AuditLog{
		*newTestAuditLog(1),
		*newTestAuditLog(2),
	}

	mockRepo.On("List", mock.Anything, mock.AnythingOfType("auditlog.FilterOptions")).Return(logs, int64(2), nil)

	handler := NewListLogsHandler(mockRepo)

	// Act
	result, err := handler.Handle(context.Background(), ListLogsQuery{
		Page:  1,
		Limit: 10,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.Logs, 2)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Limit)

	mockRepo.AssertExpectations(t)
}

func TestListLogsHandler_Handle_EmptyList(t *testing.T) {
	// Arrange
	mockRepo := new(MockAuditLogQueryRepository)

	mockRepo.On("List", mock.Anything, mock.AnythingOfType("auditlog.FilterOptions")).Return([]domainAuditLog.AuditLog{}, int64(0), nil)

	handler := NewListLogsHandler(mockRepo)

	// Act
	result, err := handler.Handle(context.Background(), ListLogsQuery{
		Page:  1,
		Limit: 10,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.Total)
	assert.Empty(t, result.Logs)
}

func TestListLogsHandler_Handle_WithFilters(t *testing.T) {
	// Arrange
	mockRepo := new(MockAuditLogQueryRepository)

	userID := uint(1)
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	var capturedFilter domainAuditLog.FilterOptions
	mockRepo.On("List", mock.Anything, mock.AnythingOfType("auditlog.FilterOptions")).Run(func(args mock.Arguments) {
		capturedFilter = args.Get(1).(domainAuditLog.FilterOptions)
	}).Return([]domainAuditLog.AuditLog{}, int64(0), nil)

	handler := NewListLogsHandler(mockRepo)

	// Act
	_, err := handler.Handle(context.Background(), ListLogsQuery{
		Page:      2,
		Limit:     20,
		UserID:    &userID,
		Action:    "login",
		Resource:  "auth",
		Status:    "success",
		StartDate: &startDate,
		EndDate:   &endDate,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 2, capturedFilter.Page)
	assert.Equal(t, 20, capturedFilter.Limit)
	assert.Equal(t, &userID, capturedFilter.UserID)
	assert.Equal(t, "login", capturedFilter.Action)
	assert.Equal(t, "auth", capturedFilter.Resource)
	assert.Equal(t, "success", capturedFilter.Status)
	assert.Equal(t, &startDate, capturedFilter.StartDate)
	assert.Equal(t, &endDate, capturedFilter.EndDate)
}

func TestListLogsHandler_Handle_Error(t *testing.T) {
	// Arrange
	mockRepo := new(MockAuditLogQueryRepository)

	mockRepo.On("List", mock.Anything, mock.AnythingOfType("auditlog.FilterOptions")).Return(nil, int64(0), errors.New("database error"))

	handler := NewListLogsHandler(mockRepo)

	// Act
	result, err := handler.Handle(context.Background(), ListLogsQuery{
		Page:  1,
		Limit: 10,
	})

	// Assert
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list audit logs")
}

func TestListLogsHandler_Handle_DTOConversion(t *testing.T) {
	// 验证 DTO 转换正确
	mockRepo := new(MockAuditLogQueryRepository)

	log := newTestAuditLog(1)
	log.Action = "create"
	log.Resource = "user"
	log.Details = `{"user_id":123}`

	mockRepo.On("List", mock.Anything, mock.AnythingOfType("auditlog.FilterOptions")).Return([]domainAuditLog.AuditLog{*log}, int64(1), nil)

	handler := NewListLogsHandler(mockRepo)

	result, err := handler.Handle(context.Background(), ListLogsQuery{Page: 1, Limit: 10})

	require.NoError(t, err)
	assert.Len(t, result.Logs, 1)

	dto := result.Logs[0]
	assert.Equal(t, log.ID, dto.ID)
	assert.Equal(t, log.UserID, dto.UserID)
	assert.Equal(t, log.Action, dto.Action)
	assert.Equal(t, log.Resource, dto.Resource)
	assert.Equal(t, log.Details, dto.Details)
	assert.Equal(t, log.Status, dto.Status)
}
