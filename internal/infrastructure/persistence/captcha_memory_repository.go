package persistence

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/captcha"
)

const (
	// captchaDefaultExpiration 默认过期时间（5分钟）
	captchaDefaultExpiration = 5 * time.Minute
	// captchaCleanupInterval 清理间隔（10分钟）
	captchaCleanupInterval = 10 * time.Minute
	// captchaMaxSize 最大存储数量
	captchaMaxSize = 10000
)

// captchaMemoryRepository 内存验证码仓储实现
// 🔒 安全策略：
// - 并发安全（使用 sync.RWMutex）
// - 验证码一次性使用，验证后立即删除
// - 自动清理过期验证码
// - LRU 策略，防止内存溢出
type captchaMemoryRepository struct {
	data      map[string]*captcha.CaptchaData
	mu        sync.RWMutex
	stopClean chan struct{}
}

var (
	_ captcha.CommandRepository = (*captchaMemoryRepository)(nil)
	_ captcha.QueryRepository   = (*captchaMemoryRepository)(nil)
)

// NewCaptchaMemoryRepository 创建内存验证码仓储
func NewCaptchaMemoryRepository() *captchaMemoryRepository {
	repo := &captchaMemoryRepository{
		data:      make(map[string]*captcha.CaptchaData),
		stopClean: make(chan struct{}),
	}

	// 启动定期清理协程
	go repo.cleanupExpired()

	return repo
}

// Create 创建验证码并存储
func (r *captchaMemoryRepository) Create(ctx context.Context, captchaID string, code string, expiration time.Duration) error {
	if expiration <= 0 {
		expiration = captchaDefaultExpiration
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// LRU 策略：如果超过最大容量，删除最旧的项
	if len(r.data) >= captchaMaxSize {
		var oldestID string
		var oldestTime time.Time
		for id, data := range r.data {
			if oldestID == "" || data.CreatedAt.Before(oldestTime) {
				oldestID = id
				oldestTime = data.CreatedAt
			}
		}
		if oldestID != "" {
			delete(r.data, oldestID)
		}
	}

	// 存储新项（验证码值统一转换为小写）
	now := time.Now()
	r.data[captchaID] = &captcha.CaptchaData{
		Code:      strings.ToLower(code),
		ExpireAt:  now.Add(expiration),
		CreatedAt: now,
	}

	return nil
}

// Verify 验证验证码（不区分大小写，一次性使用）
func (r *captchaMemoryRepository) Verify(ctx context.Context, captchaID string, code string) (bool, error) {
	if captchaID == "" || code == "" {
		return false, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 获取验证码数据
	captchaData, exists := r.data[captchaID]
	if !exists {
		return false, nil
	}

	// 检查是否过期
	if captchaData.IsExpired() {
		delete(r.data, captchaID)
		return false, nil
	}

	// 比较验证码（不区分大小写）
	isValid := strings.EqualFold(captchaData.Code, code)

	// 无论验证成功或失败都删除（一次性使用）
	delete(r.data, captchaID)

	return isValid, nil
}

// Delete 删除验证码
func (r *captchaMemoryRepository) Delete(ctx context.Context, captchaID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, captchaID)
	return nil
}

// GetStats 获取统计信息
func (r *captchaMemoryRepository) GetStats(ctx context.Context) map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	expired := 0
	for _, data := range r.data {
		if data.IsExpired() {
			expired++
		}
	}

	return map[string]interface{}{
		"total":   len(r.data),
		"expired": expired,
		"active":  len(r.data) - expired,
	}
}

// cleanupExpired 定期清理过期验证码
func (r *captchaMemoryRepository) cleanupExpired() {
	ticker := time.NewTicker(captchaCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.mu.Lock()
			for id, data := range r.data {
				if data.IsExpired() {
					delete(r.data, id)
				}
			}
			r.mu.Unlock()
		case <-r.stopClean:
			return
		}
	}
}

// Close 关闭仓储（停止清理协程）
func (r *captchaMemoryRepository) Close() error {
	close(r.stopClean)
	return nil
}

// CaptchaError 验证码错误
type CaptchaError struct {
	Message string
}

func (e *CaptchaError) Error() string {
	return fmt.Sprintf("captcha error: %s", e.Message)
}
