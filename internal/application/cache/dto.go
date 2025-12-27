package cache

// SetDTO 设置缓存请求
type SetDTO struct {
	Key   string `json:"key" binding:"required"`
	Value any    `json:"value" binding:"required"`
	TTL   int    `json:"ttl"` // 秒，默认 60
}

// GetResultDTO 获取缓存结果
type GetResultDTO struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

// SetResultDTO 设置缓存结果
type SetResultDTO struct {
	Key string `json:"key"`
	TTL int    `json:"ttl"`
}

// DeleteResultDTO 删除缓存结果
type DeleteResultDTO struct {
	Key string `json:"key"`
}
