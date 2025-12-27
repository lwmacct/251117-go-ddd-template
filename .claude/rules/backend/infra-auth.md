---
paths:
  - "internal/infrastructure/auth/**/*.go"
---

# Auth Infrastructure 规范

## 核心职责

实现 Domain 层 `auth.Service` 接口，处理认证相关技术细节（JWT、密码哈希、会话管理）。

## 文件结构

| 文件                          | 职责                       |
| ----------------------------- | -------------------------- |
| `auth_service_impl.go`        | 实现 `auth.Service` 接口   |
| `jwt.go`                      | JWT 生成与验证             |
| `token_generator.go`          | Token 生成器               |
| `pat_service.go`              | Personal Access Token 服务 |
| `login_session.go`            | 登录会话管理               |
| `permission_cache_service.go` | 权限缓存服务               |

## Domain Service 实现

```go
// auth_service_impl.go - 实现 domain/auth.Service 接口
type authService struct {
    jwtManager *JWTManager
}

func NewAuthService(jwtManager *JWTManager) auth.Service {
    return &authService{jwtManager: jwtManager}
}

func (s *authService) HashPassword(password string) (string, error) { ... }
func (s *authService) VerifyPassword(hashedPassword, password string) error { ... }
```

## 依赖方向

```
domain/auth.Service (接口)
         ↑
infrastructure/auth (实现)
```
