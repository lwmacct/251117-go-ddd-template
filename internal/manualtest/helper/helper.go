// Package helper 提供手动测试辅助函数和 HTTP 客户端。
package helper

import (
	"os"
	"testing"
)

const (
	DefaultBaseURL   = "http://localhost:40012"
	DefaultDevSecret = "dev-secret-change-me"
)

// SkipIfNotManual 如果 MANUAL 环境变量未设置则跳过测试。
func SkipIfNotManual(t *testing.T) {
	t.Helper()
	if os.Getenv("MANUAL") == "" {
		t.SkipNow()
	}
}

// NewClient 创建测试客户端，从环境变量读取配置。
func NewClient() *Client {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	devSecret := os.Getenv("DEV_SECRET")
	if devSecret == "" {
		devSecret = DefaultDevSecret
	}
	return newClient(baseURL, devSecret)
}
