package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	userdto "github.com/lwmacct/251117-go-ddd-template/internal/application/user"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler 用户处理器
type UserHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userCommandRepo user.CommandRepository, userQueryRepo user.QueryRepository) *UserHandler {
	return &UserHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
	}
}

// Create 创建用户
// POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var dto userdto.CreateUserDTO
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

	if err := h.userCommandRepo.Create(c.Request.Context(), newUser); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "user created successfully",
		"data":    userdto.ToUserResponse(newUser),
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

	u, err := h.userQueryRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": userdto.ToUserResponse(u),
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

	users, err := h.userQueryRepo.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应 DTO
	responses := make([]*userdto.UserResponse, len(users))
	for i, u := range users {
		responses[i] = userdto.ToUserResponse(u)
	}

	// 获取总数
	total, err := h.userQueryRepo.Count(c.Request.Context())
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

	var dto userdto.UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 获取现有用户
	u, err := h.userQueryRepo.GetByID(c.Request.Context(), uint(id))
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

	if err := h.userCommandRepo.Update(c.Request.Context(), u); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "user updated successfully",
		"data":    userdto.ToUserResponse(u),
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

	if err := h.userCommandRepo.Delete(c.Request.Context(), uint(id)); err != nil{
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "user deleted successfully",
	})
}
