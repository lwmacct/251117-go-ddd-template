# 缓存接口

本文档说明缓存操作的 REST API 接口。

## 概述

缓存接口提供基本的 Redis 缓存操作功能，支持设置、获取和删除缓存数据。

**Base URL**: `http://localhost:8080`

::: tip 注意
这些接口主要用于演示和调试。生产环境中，缓存操作通常封装在业务逻辑内部，不直接暴露给前端。
:::

## 接口列表

| 方法   | 路径              | 描述     | 认证 |
| ------ | ----------------- | -------- | ---- |
| POST   | `/api/cache`      | 设置缓存 | 无   |
| GET    | `/api/cache/:key` | 获取缓存 | 无   |
| DELETE | `/api/cache/:key` | 删除缓存 | 无   |

## 设置缓存

设置键值对缓存，支持自定义过期时间。

### 请求

- **方法**: `POST`
- **路径**: `/api/cache`
- **Content-Type**: `application/json`

#### 请求参数

| 参数  | 类型    | 必填 | 默认值 | 说明                        |
| ----- | ------- | ---- | ------ | --------------------------- |
| key   | string  | 是   | -      | 缓存键名                    |
| value | any     | 是   | -      | 缓存值 (支持任意 JSON 类型) |
| ttl   | integer | 否   | 60     | 过期时间 (秒)               |

#### 请求示例

**字符串值**：

```json
{
  "key": "greeting",
  "value": "Hello, World!",
  "ttl": 300
}
```

**对象值**：

```json
{
  "key": "user:1001",
  "value": {
    "id": 1001,
    "username": "testuser",
    "email": "test@example.com"
  },
  "ttl": 600
}
```

**数组值**：

```json
{
  "key": "tags",
  "value": ["go", "redis", "cache"],
  "ttl": 120
}
```

### 响应

#### 成功响应 (200)

```json
{
  "message": "cache set successfully",
  "key": "greeting",
  "ttl": 300
}
```

#### 错误响应 (400)

```json
{
  "error": "Key: 'key' Error:Field validation for 'key' failed on the 'required' tag"
}
```

#### 错误响应 (500)

```json
{
  "error": "redis: connection refused"
}
```

### 代码示例

**cURL**：

```bash
curl -X POST http://localhost:8080/api/cache \
  -H "Content-Type: application/json" \
  -d '{
    "key": "test_key",
    "value": "test_value",
    "ttl": 300
  }'
```

**JavaScript (Fetch)**：

```javascript
fetch("http://localhost:8080/api/cache", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    key: "test_key",
    value: "test_value",
    ttl: 300,
  }),
})
  .then((response) => response.json())
  .then((data) => console.log(data));
```

**Go**：

```go
import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SetCacheRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl"`
}

func setCache() error {
	req := SetCacheRequest{
		Key:   "test_key",
		Value: "test_value",
		TTL:   300,
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post(
		"http://localhost:8080/api/cache",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
```

**Python**：

```python
import requests

data = {
    "key": "test_key",
    "value": "test_value",
    "ttl": 300
}

response = requests.post(
    'http://localhost:8080/api/cache',
    json=data
)
print(response.json())
```

## 获取缓存

根据键名获取缓存值。

### 请求

- **方法**: `GET`
- **路径**: `/api/cache/:key`

#### 路径参数

| 参数 | 类型   | 必填 | 说明     |
| ---- | ------ | ---- | -------- |
| key  | string | 是   | 缓存键名 |

#### 请求示例

```bash
GET /api/cache/greeting
```

### 响应

#### 成功响应 (200)

**字符串值**：

```json
{
  "key": "greeting",
  "value": "Hello, World!"
}
```

**对象值**：

```json
{
  "key": "user:1001",
  "value": {
    "id": 1001,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

**数组值**：

```json
{
  "key": "tags",
  "value": ["go", "redis", "cache"]
}
```

#### 错误响应 (400)

```json
{
  "error": "key is required"
}
```

#### 错误响应 (404)

```json
{
  "error": "redis: nil"
}
```

::: info 说明
当缓存键不存在或已过期时，返回 404 错误。
:::

### 代码示例

**cURL**：

```bash
curl http://localhost:8080/api/cache/greeting
```

**JavaScript (Fetch)**：

```javascript
fetch("http://localhost:8080/api/cache/greeting")
  .then((response) => {
    if (!response.ok) {
      throw new Error("Cache not found");
    }
    return response.json();
  })
  .then((data) => console.log(data.value));
```

**Go**：

```go
import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetCacheResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func getCache(key string) (*GetCacheResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/cache/%s", key))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("cache not found")
	}

	var result GetCacheResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
```

**Python**：

```python
import requests

response = requests.get('http://localhost:8080/api/cache/greeting')

if response.status_code == 200:
    data = response.json()
    print(f"Value: {data['value']}")
elif response.status_code == 404:
    print("Cache not found or expired")
else:
    print(f"Error: {response.json()['error']}")
```

## 删除缓存

删除指定键的缓存。

### 请求

- **方法**: `DELETE`
- **路径**: `/api/cache/:key`

#### 路径参数

| 参数 | 类型   | 必填 | 说明     |
| ---- | ------ | ---- | -------- |
| key  | string | 是   | 缓存键名 |

#### 请求示例

```bash
DELETE /api/cache/greeting
```

### 响应

#### 成功响应 (200)

```json
{
  "message": "cache deleted successfully",
  "key": "greeting"
}
```

::: info 说明
即使键不存在，删除操作也会返回成功 (幂等操作) 。
:::

#### 错误响应 (400)

```json
{
  "error": "key is required"
}
```

#### 错误响应 (500)

```json
{
  "error": "redis: connection refused"
}
```

### 代码示例

**cURL**：

```bash
curl -X DELETE http://localhost:8080/api/cache/greeting
```

**JavaScript (Fetch)**：

```javascript
fetch("http://localhost:8080/api/cache/greeting", {
  method: "DELETE",
})
  .then((response) => response.json())
  .then((data) => console.log(data.message));
```

**Go**：

```go
import (
	"fmt"
	"net/http"
)

func deleteCache(key string) error {
	req, _ := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("http://localhost:8080/api/cache/%s", key),
		nil,
	)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
```

**Python**：

```python
import requests

response = requests.delete('http://localhost:8080/api/cache/greeting')

if response.status_code == 200:
    print(response.json()['message'])
else:
    print(f"Error: {response.json()['error']}")
```

## 错误码

| 状态码 | 说明       | 原因                     |
| ------ | ---------- | ------------------------ |
| 200    | 成功       | 操作成功                 |
| 400    | 请求错误   | 参数缺失或格式错误       |
| 404    | 未找到     | 缓存键不存在或已过期     |
| 500    | 服务器错误 | Redis 连接失败或内部错误 |

## 使用场景

### 1. 会话存储

```bash
# 设置会话
curl -X POST http://localhost:8080/api/cache \
  -H "Content-Type: application/json" \
  -d '{
    "key": "session:abc123",
    "value": {"user_id": 1, "expires": 1234567890},
    "ttl": 3600
  }'

# 获取会话
curl http://localhost:8080/api/cache/session:abc123

# 删除会话 (登出)
curl -X DELETE http://localhost:8080/api/cache/session:abc123
```

### 2. 临时数据缓存

```bash
# 缓存API响应
curl -X POST http://localhost:8080/api/cache \
  -H "Content-Type: application/json" \
  -d '{
    "key": "api:weather:beijing",
    "value": {"temp": 25, "humidity": 60},
    "ttl": 300
  }'
```

### 3. 限流计数器

```bash
# 记录请求次数 (配合 TTL 实现滑动窗口)
curl -X POST http://localhost:8080/api/cache \
  -H "Content-Type: application/json" \
  -d '{
    "key": "rate:192.168.1.1:minute",
    "value": 1,
    "ttl": 60
  }'
```

## 最佳实践

### 1. 键命名规范

使用命名空间分隔键名，便于管理：

```
格式: <namespace>:<entity>:<id>

示例:
- user:profile:1001
- session:token:abc123
- cache:api:weather:beijing
- rate:limit:192.168.1.1
```

### 2. TTL 设置建议

| 数据类型 | 推荐 TTL  | 说明             |
| -------- | --------- | ---------------- |
| 会话数据 | 1-24 小时 | 根据业务需求     |
| API 缓存 | 5-30 分钟 | 根据数据更新频率 |
| 限流计数 | 1-60 秒   | 时间窗口大小     |
| 临时数据 | 1-5 分钟  | 短期缓存         |

### 3. 值类型选择

```json
// ✅ 推荐：结构化数据
{
  "key": "user:1001",
  "value": {
    "id": 1001,
    "name": "test",
    "timestamp": 1234567890
  }
}

// ❌ 避免：过大的对象
{
  "key": "user:1001:full",
  "value": {
    // 避免存储大量数据 (>1MB)
  }
}
```

### 4. 错误处理

```javascript
async function getCacheWithFallback(key) {
  try {
    const response = await fetch(`http://localhost:8080/api/cache/${key}`);
    if (response.ok) {
      const data = await response.json();
      return data.value;
    }
  } catch (error) {
    console.error("Cache error:", error);
  }

  // 缓存未命中，从数据库加载
  return fetchFromDatabase(key);
}
```

## 性能考虑

### 1. 批量操作

::: warning 注意
当前接口不支持批量操作。如需批量设置/获取缓存，建议直接使用 CacheRepository。
:::

### 2. 大对象缓存

::: warning 警告
避免缓存超过 **1MB** 的对象，可能导致性能问题。大对象建议拆分或使用对象存储。
:::

### 3. 高频访问优化

对于高频访问的缓存键：

- 设置合理的 TTL，避免频繁更新
- 考虑使用本地内存缓存 (双层缓存)
- 使用 Redis 集群分散压力

## 相关链接

- [Redis 集成指南](/guide/redis) - Redis 配置和高级用法
- [认证接口](/api/auth) - JWT 认证接口
- [用户接口](/api/users) - 用户管理接口
- [配置系统](/guide/configuration) - Redis 连接配置

## 技术实现

缓存接口基于以下技术实现：

- **Redis 客户端**: `github.com/redis/go-redis/v9`
- **HTTP 框架**: Gin
- **Handler 位置**: `internal/adapters/http/handler/cache.go:26`
- **路由配置**: `internal/adapters/http/router.go`

查看源码了解更多实现细节。
