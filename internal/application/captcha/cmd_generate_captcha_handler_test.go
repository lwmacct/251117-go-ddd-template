package captcha

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGenerateCaptchaHandler_Handle_RandomMode_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCaptchaCommandRepository)
	mockService := new(MockCaptchaService)

	mockService.On("GenerateRandomCode").Return("captcha_123", "base64_image_data", "ABC123", nil)
	mockService.On("GetDefaultExpiration").Return(int64(300))
	mockRepo.On("Create", mock.Anything, "captcha_123", "ABC123", 300*time.Second).Return(nil)

	handler := NewGenerateCaptchaHandler(mockRepo, mockService)

	// Act
	result, err := handler.Handle(context.Background(), GenerateCaptchaCommand{
		DevMode: false,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "captcha_123", result.ID)
	assert.Equal(t, "base64_image_data", result.Image)
	assert.Empty(t, result.Code) // 非开发模式不返回验证码
	assert.Positive(t, result.ExpireAt)

	mockRepo.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGenerateCaptchaHandler_Handle_DevMode_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCaptchaCommandRepository)
	mockService := new(MockCaptchaService)

	mockService.On("GenerateRandomCode").Return("captcha_456", "base64_random_image", "XYZ789", nil)
	mockService.On("GetDefaultExpiration").Return(int64(300))
	mockRepo.On("Create", mock.Anything, "captcha_456", "XYZ789", 300*time.Second).Return(nil)

	handler := NewGenerateCaptchaHandler(mockRepo, mockService)

	// Act
	result, err := handler.Handle(context.Background(), GenerateCaptchaCommand{
		DevMode: true,
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "captcha_456", result.ID)
	assert.Equal(t, "base64_random_image", result.Image)
	assert.Equal(t, "XYZ789", result.Code) // 开发模式返回验证码
}

func TestGenerateCaptchaHandler_Handle_DevMode_CustomCode_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCaptchaCommandRepository)
	mockService := new(MockCaptchaService)

	mockService.On("GenerateDevCaptchaID").Return("dev_captcha_001")
	mockService.On("GenerateCustomCodeImage", "111111").Return("base64_custom_image", nil)
	mockService.On("GetDefaultExpiration").Return(int64(600))
	mockRepo.On("Create", mock.Anything, "dev_captcha_001", "111111", 600*time.Second).Return(nil)

	handler := NewGenerateCaptchaHandler(mockRepo, mockService)

	// Act
	result, err := handler.Handle(context.Background(), GenerateCaptchaCommand{
		DevMode:    true,
		CustomCode: "111111",
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "dev_captcha_001", result.ID)
	assert.Equal(t, "base64_custom_image", result.Image)
	assert.Equal(t, "111111", result.Code) // 开发模式返回自定义验证码

	mockRepo.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGenerateCaptchaHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        GenerateCaptchaCommand
		setupMocks func(*MockCaptchaCommandRepository, *MockCaptchaService)
		wantErr    string
	}{
		{
			name: "生成随机验证码失败",
			cmd:  GenerateCaptchaCommand{DevMode: false},
			setupMocks: func(repo *MockCaptchaCommandRepository, service *MockCaptchaService) {
				service.On("GenerateRandomCode").Return("", "", "", errors.New("generation failed"))
			},
			wantErr: "generation failed",
		},
		{
			name: "生成自定义验证码图片失败",
			cmd:  GenerateCaptchaCommand{DevMode: true, CustomCode: "123456"},
			setupMocks: func(repo *MockCaptchaCommandRepository, service *MockCaptchaService) {
				service.On("GenerateDevCaptchaID").Return("dev_id")
				service.On("GenerateCustomCodeImage", "123456").Return("", errors.New("image generation failed"))
			},
			wantErr: "image generation failed",
		},
		{
			name: "存储验证码失败",
			cmd:  GenerateCaptchaCommand{DevMode: false},
			setupMocks: func(repo *MockCaptchaCommandRepository, service *MockCaptchaService) {
				service.On("GenerateRandomCode").Return("id", "img", "code", nil)
				service.On("GetDefaultExpiration").Return(int64(300))
				repo.On("Create", mock.Anything, "id", "code", 300*time.Second).Return(errors.New("storage error"))
			},
			wantErr: "storage error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockCaptchaCommandRepository)
			mockService := new(MockCaptchaService)
			tt.setupMocks(mockRepo, mockService)

			handler := NewGenerateCaptchaHandler(mockRepo, mockService)

			// Act
			result, err := handler.Handle(context.Background(), tt.cmd)

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestGenerateCaptchaHandler_Handle_ExpirationTime(t *testing.T) {
	// 验证过期时间计算正确
	mockRepo := new(MockCaptchaCommandRepository)
	mockService := new(MockCaptchaService)

	expiration := int64(600) // 10 分钟
	mockService.On("GenerateRandomCode").Return("id", "img", "code", nil)
	mockService.On("GetDefaultExpiration").Return(expiration)
	mockRepo.On("Create", mock.Anything, "id", "code", time.Duration(expiration)*time.Second).Return(nil)

	handler := NewGenerateCaptchaHandler(mockRepo, mockService)
	beforeTime := time.Now().Unix()

	result, err := handler.Handle(context.Background(), GenerateCaptchaCommand{})

	afterTime := time.Now().Unix()

	require.NoError(t, err)
	assert.NotNil(t, result)
	// 过期时间应该在 [beforeTime + expiration, afterTime + expiration] 范围内
	assert.GreaterOrEqual(t, result.ExpireAt, beforeTime+expiration)
	assert.LessOrEqual(t, result.ExpireAt, afterTime+expiration)
}
