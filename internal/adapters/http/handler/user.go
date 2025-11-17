package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-bd-vmalert/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler 用户处理器
type UserHandler struct {
	userRepo user.Repository
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userRepo user.Repository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// Create 创建用户
// POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var dto user.UserCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to hash password"})
		return
	}

	newUser := &user.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: string(hashedPassword),
		FullName: dto.FullName,
		Status:   "active",
	}

	if err := h.userRepo.Create(c.Request.Context(), newUser); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "user created successfully",
		"data":    newUser.ToResponse(),
	})
}

// GetByID 获取用户详情
// GET /api/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	u, err := h.userRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": u.ToResponse(),
	})
}

// List 获取用户列表
// GET /api/users?page=1&limit=10
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := h.userRepo.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应 DTO
	responses := make([]*user.UserResponse, len(users))
	for i, u := range users {
		responses[i] = u.ToResponse()
	}

	// 获取总数
	total, err := h.userRepo.Count(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// Update 更新用户
// PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	var dto user.UserUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 获取现有用户
	u, err := h.userRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	if dto.FullName != nil {
		u.FullName = *dto.FullName
	}
	if dto.Avatar != nil {
		u.Avatar = *dto.Avatar
	}
	if dto.Bio != nil {
		u.Bio = *dto.Bio
	}
	if dto.Status != nil {
		u.Status = *dto.Status
	}

	if err := h.userRepo.Update(c.Request.Context(), u); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "user updated successfully",
		"data":    u.ToResponse(),
	})
}

// Delete 删除用户
// DELETE /api/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.userRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "user deleted successfully",
	})
}
