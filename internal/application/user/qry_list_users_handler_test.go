package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainUser "github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func TestListUsersHandler_Handle_Success(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		query ListUsersQuery
		users []*domainUser.User
		total int64
	}{
		{
			name:  "列出第一页用户",
			query: ListUsersQuery{Page: 1, Limit: 10},
			users: []*domainUser.User{
				{ID: 1, Username: "user1", Email: "user1@example.com", Status: "active", CreatedAt: now, UpdatedAt: now},
				{ID: 2, Username: "user2", Email: "user2@example.com", Status: "active", CreatedAt: now, UpdatedAt: now},
			},
			total: 2,
		},
		{
			name:  "空用户列表",
			query: ListUsersQuery{Page: 1, Limit: 10},
			users: []*domainUser.User{},
			total: 0,
		},
		{
			name:  "分页查询",
			query: ListUsersQuery{Page: 2, Limit: 5},
			users: []*domainUser.User{
				{ID: 6, Username: "user6", Email: "user6@example.com", Status: "active", CreatedAt: now, UpdatedAt: now},
			},
			total: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockUserQueryRepository)
			offset := (tt.query.Page - 1) * tt.query.Limit

			mockQryRepo.On("List", mock.Anything, offset, tt.query.Limit).Return(tt.users, nil)
			mockQryRepo.On("Count", mock.Anything).Return(tt.total, nil)

			handler := NewListUsersHandler(mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.total, result.Total)
			assert.Len(t, result.Users, len(tt.users))
			mockQryRepo.AssertExpectations(t)
		})
	}
}

func TestListUsersHandler_Handle_WithSearch(t *testing.T) {
	now := time.Now()

	// Arrange
	mockQryRepo := new(MockUserQueryRepository)

	users := []*domainUser.User{
		{ID: 1, Username: "john", Email: "john@example.com", Status: "active", CreatedAt: now, UpdatedAt: now},
	}

	mockQryRepo.On("Search", mock.Anything, "john", 0, 10).Return(users, nil)
	mockQryRepo.On("CountBySearch", mock.Anything, "john").Return(int64(1), nil)

	handler := NewListUsersHandler(mockQryRepo)

	// Act
	result, err := handler.Handle(context.Background(), ListUsersQuery{
		Page:   1,
		Limit:  10,
		Search: "john",
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.Total)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "john", result.Users[0].Username)
	mockQryRepo.AssertExpectations(t)
}

func TestListUsersHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		query      ListUsersQuery
		setupMocks func(*MockUserQueryRepository)
		wantErr    string
	}{
		{
			name:  "列表查询数据库错误",
			query: ListUsersQuery{Page: 1, Limit: 10},
			setupMocks: func(qryRepo *MockUserQueryRepository) {
				qryRepo.On("List", mock.Anything, 0, 10).Return(nil, errors.New("database error"))
			},
			wantErr: "database error",
		},
		{
			name:  "计数数据库错误",
			query: ListUsersQuery{Page: 1, Limit: 10},
			setupMocks: func(qryRepo *MockUserQueryRepository) {
				qryRepo.On("List", mock.Anything, 0, 10).Return([]*domainUser.User{}, nil)
				qryRepo.On("Count", mock.Anything).Return(int64(0), errors.New("count error"))
			},
			wantErr: "count error",
		},
		{
			name:  "搜索数据库错误",
			query: ListUsersQuery{Page: 1, Limit: 10, Search: "test"},
			setupMocks: func(qryRepo *MockUserQueryRepository) {
				qryRepo.On("Search", mock.Anything, "test", 0, 10).Return(nil, errors.New("search error"))
			},
			wantErr: "search error",
		},
		{
			name:  "搜索计数错误",
			query: ListUsersQuery{Page: 1, Limit: 10, Search: "john"},
			setupMocks: func(qryRepo *MockUserQueryRepository) {
				qryRepo.On("Search", mock.Anything, "john", 0, 10).Return([]*domainUser.User{}, nil)
				qryRepo.On("CountBySearch", mock.Anything, "john").Return(int64(0), errors.New("count search error"))
			},
			wantErr: "count search error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockQryRepo := new(MockUserQueryRepository)
			tt.setupMocks(mockQryRepo)

			handler := NewListUsersHandler(mockQryRepo)

			// Act
			result, err := handler.Handle(context.Background(), tt.query)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestListUsersQuery_GetOffset(t *testing.T) {
	tests := []struct {
		name   string
		query  ListUsersQuery
		offset int
	}{
		{name: "第一页", query: ListUsersQuery{Page: 1, Limit: 10}, offset: 0},
		{name: "第二页", query: ListUsersQuery{Page: 2, Limit: 10}, offset: 10},
		{name: "第三页", query: ListUsersQuery{Page: 3, Limit: 5}, offset: 10},
		{name: "页码为0时默认第一页", query: ListUsersQuery{Page: 0, Limit: 10}, offset: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.offset, tt.query.GetOffset())
		})
	}
}
