package cache

// SetCacheDTO 设置缓存请求
type SetCacheDTO struct {
	Key   string `json:"key" binding:"required"`
	Value any    `json:"value" binding:"required"`
	TTL   int    `json:"ttl"` // 秒，默认 60
}

// GetCacheResultDTO 获取缓存结果
type GetCacheResultDTO struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

// SetCacheResultDTO 设置缓存结果
type SetCacheResultDTO struct {
	Key string `json:"key"`
	TTL int    `json:"ttl"`
}

// DeleteCacheResultDTO 删除缓存结果
type DeleteCacheResultDTO struct {
	Key string `json:"key"`
}
