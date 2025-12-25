package persistence

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB 创建测试用 SQLite in-memory 数据库。
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "无法连接测试数据库")

	// 迁移所有需要的表
	err = db.AutoMigrate(&UserModel{}, &RoleModel{}, &PermissionModel{})
	require.NoError(t, err, "数据库迁移失败")

	return db
}

// createTestRole 在数据库中创建测试角色。
func createTestRole(t *testing.T, db *gorm.DB, name string) *RoleModel {
	t.Helper()

	role := &RoleModel{
		Name:        name,
		DisplayName: name + " Display",
		Description: "Test role: " + name,
	}
	err := db.Create(role).Error
	require.NoError(t, err, "创建测试角色失败")
	return role
}

func TestUserCommandRepository_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("成功创建用户", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserCommandRepository(db)

		u := &user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed_password",
			FullName: "Test User",
			Status:   "active",
		}

		err := repo.Create(ctx, u)

		require.NoError(t, err)
		assert.NotZero(t, u.ID, "ID 应该被设置")
		assert.NotZero(t, u.CreatedAt, "CreatedAt 应该被设置")

		// 验证数据库中的记录
		var model UserModel
		err = db.First(&model, u.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "testuser", model.Username)
		assert.Equal(t, "test@example.com", model.Email)
	})

	t.Run("用户名唯一性约束", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserCommandRepository(db)

		// 创建第一个用户
		u1 := &user.User{
			Username: "duplicate",
			Email:    "first@example.com",
			Password: "password",
			Status:   "active",
		}
		err := repo.Create(ctx, u1)
		require.NoError(t, err)

		// 尝试创建相同用户名的用户
		u2 := &user.User{
			Username: "duplicate",
			Email:    "second@example.com",
			Password: "password",
			Status:   "active",
		}
		err = repo.Create(ctx, u2)

		assert.Error(t, err, "应该返回唯一性约束错误")
	})

	t.Run("邮箱唯一性约束", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserCommandRepository(db)

		u1 := &user.User{
			Username: "user1",
			Email:    "duplicate@example.com",
			Password: "password",
			Status:   "active",
		}
		err := repo.Create(ctx, u1)
		require.NoError(t, err)

		u2 := &user.User{
			Username: "user2",
			Email:    "duplicate@example.com",
			Password: "password",
			Status:   "active",
		}
		err = repo.Create(ctx, u2)

		assert.Error(t, err, "应该返回邮箱唯一性约束错误")
	})
}

func TestUserCommandRepository_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("成功更新用户", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserCommandRepository(db)

		// 先创建用户
		u := &user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			FullName: "Original Name",
			Status:   "active",
		}
		err := repo.Create(ctx, u)
		require.NoError(t, err)

		// 更新用户
		u.FullName = "Updated Name"
		u.Bio = "New bio"
		err = repo.Update(ctx, u)

		require.NoError(t, err)

		// 验证更新
		var model UserModel
		err = db.First(&model, u.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", model.FullName)
		assert.Equal(t, "New bio", model.Bio)
	})
}

func TestUserCommandRepository_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("成功软删除用户", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserCommandRepository(db)

		// 创建用户
		u := &user.User{
			Username: "todelete",
			Email:    "delete@example.com",
			Password: "password",
			Status:   "active",
		}
		err := repo.Create(ctx, u)
		require.NoError(t, err)

		// 删除用户
		err = repo.Delete(ctx, u.ID)
		require.NoError(t, err)

		// 验证软删除（默认查询不应该找到）
		var model UserModel
		err = db.First(&model, u.ID).Error
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// 但 Unscoped 可以找到
		err = db.Unscoped().First(&model, u.ID).Error
		require.NoError(t, err)
		assert.True(t, model.DeletedAt.Valid, "DeletedAt 应该被设置")
	})
}

func TestUserCommandRepository_AssignRoles(t *testing.T) {
	ctx := context.Background()

	t.Run("成功分配角色", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserCommandRepository(db)

		// 创建角色
		role1 := createTestRole(t, db, "admin")
		role2 := createTestRole(t, db, "editor")

		// 创建用户
		u := &user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Status:   "active",
		}
		err := repo.Create(ctx, u)
		require.NoError(t, err)

		// 分配角色
		err = repo.AssignRoles(ctx, u.ID, []uint{role1.ID, role2.ID})
		require.NoError(t, err)

		// 验证角色分配
		var model UserModel
		err = db.Preload("Roles").First(&model, u.ID).Error
		require.NoError(t, err)
		assert.Len(t, model.Roles, 2)
	})

	t.Run("替换现有角色", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserCommandRepository(db)

		role1 := createTestRole(t, db, "role1")
		role2 := createTestRole(t, db, "role2")
		role3 := createTestRole(t, db, "role3")

		u := &user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Status:   "active",
		}
		err := repo.Create(ctx, u)
		require.NoError(t, err)

		// 先分配 role1, role2
		err = repo.AssignRoles(ctx, u.ID, []uint{role1.ID, role2.ID})
		require.NoError(t, err)

		// 替换为 role3
		err = repo.AssignRoles(ctx, u.ID, []uint{role3.ID})
		require.NoError(t, err)

		// 验证只有 role3
		var model UserModel
		err = db.Preload("Roles").First(&model, u.ID).Error
		require.NoError(t, err)
		assert.Len(t, model.Roles, 1)
		assert.Equal(t, "role3", model.Roles[0].Name)
	})

	t.Run("用户不存在", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserCommandRepository(db)

		err := repo.AssignRoles(ctx, 99999, []uint{1})

		assert.ErrorIs(t, err, user.ErrUserNotFound)
	})
}

func TestUserCommandRepository_RemoveRoles(t *testing.T) {
	ctx := context.Background()

	t.Run("成功移除角色", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserCommandRepository(db)

		role1 := createTestRole(t, db, "role1")
		role2 := createTestRole(t, db, "role2")

		u := &user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Status:   "active",
		}
		err := repo.Create(ctx, u)
		require.NoError(t, err)

		// 分配角色
		err = repo.AssignRoles(ctx, u.ID, []uint{role1.ID, role2.ID})
		require.NoError(t, err)

		// 移除 role1
		err = repo.RemoveRoles(ctx, u.ID, []uint{role1.ID})
		require.NoError(t, err)

		// 验证只剩 role2
		var model UserModel
		err = db.Preload("Roles").First(&model, u.ID).Error
		require.NoError(t, err)
		assert.Len(t, model.Roles, 1)
		assert.Equal(t, "role2", model.Roles[0].Name)
	})
}

func TestUserCommandRepository_UpdatePassword(t *testing.T) {
	ctx := context.Background()

	t.Run("成功更新密码", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserCommandRepository(db)

		u := &user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "old_password",
			Status:   "active",
		}
		err := repo.Create(ctx, u)
		require.NoError(t, err)

		// 更新密码
		err = repo.UpdatePassword(ctx, u.ID, "new_password_hash")
		require.NoError(t, err)

		// 验证
		var model UserModel
		err = db.First(&model, u.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "new_password_hash", model.Password)
	})
}

func TestUserCommandRepository_UpdateStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("成功更新状态", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserCommandRepository(db)

		u := &user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Status:   "active",
		}
		err := repo.Create(ctx, u)
		require.NoError(t, err)

		// 更新状态
		err = repo.UpdateStatus(ctx, u.ID, "banned")
		require.NoError(t, err)

		// 验证
		var model UserModel
		err = db.First(&model, u.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "banned", model.Status)
	})
}

func TestUserQueryRepository_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取用户", func(t *testing.T) {
		db := setupTestDB(t)
		cmdRepo := NewUserCommandRepository(db)
		queryRepo := NewUserQueryRepository(db)

		// 创建用户
		u := &user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			FullName: "Test User",
			Status:   "active",
		}
		err := cmdRepo.Create(ctx, u)
		require.NoError(t, err)

		// 查询
		found, err := queryRepo.GetByID(ctx, u.ID)

		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, u.ID, found.ID)
		assert.Equal(t, "testuser", found.Username)
		assert.Equal(t, "Test User", found.FullName)
	})

	t.Run("用户不存在", func(t *testing.T) {
		db := setupTestDB(t)
		queryRepo := NewUserQueryRepository(db)

		_, err := queryRepo.GetByID(ctx, 99999)

		assert.ErrorIs(t, err, user.ErrUserNotFound)
	})
}

func TestUserQueryRepository_GetByUsername(t *testing.T) {
	ctx := context.Background()

	t.Run("成功通过用户名获取", func(t *testing.T) {
		db := setupTestDB(t)
		cmdRepo := NewUserCommandRepository(db)
		queryRepo := NewUserQueryRepository(db)

		u := &user.User{
			Username: "findme",
			Email:    "find@example.com",
			Password: "password",
			Status:   "active",
		}
		err := cmdRepo.Create(ctx, u)
		require.NoError(t, err)

		found, err := queryRepo.GetByUsername(ctx, "findme")

		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, "findme", found.Username)
	})

	t.Run("用户名不存在", func(t *testing.T) {
		db := setupTestDB(t)
		queryRepo := NewUserQueryRepository(db)

		_, err := queryRepo.GetByUsername(ctx, "nonexistent")

		assert.ErrorIs(t, err, user.ErrUserNotFound)
	})
}

func TestUserQueryRepository_GetByEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("成功通过邮箱获取", func(t *testing.T) {
		db := setupTestDB(t)
		cmdRepo := NewUserCommandRepository(db)
		queryRepo := NewUserQueryRepository(db)

		u := &user.User{
			Username: "emailuser",
			Email:    "unique@example.com",
			Password: "password",
			Status:   "active",
		}
		err := cmdRepo.Create(ctx, u)
		require.NoError(t, err)

		found, err := queryRepo.GetByEmail(ctx, "unique@example.com")

		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, "unique@example.com", found.Email)
	})
}

func TestUserQueryRepository_ExistsByUsername(t *testing.T) { //nolint:dupl // 测试代码结构相似是可接受的
	ctx := context.Background()

	t.Run("用户名存在", func(t *testing.T) {
		db := setupTestDB(t)
		cmdRepo := NewUserCommandRepository(db)
		queryRepo := NewUserQueryRepository(db)

		u := &user.User{
			Username: "exists",
			Email:    "exists@example.com",
			Password: "password",
			Status:   "active",
		}
		err := cmdRepo.Create(ctx, u)
		require.NoError(t, err)

		exists, err := queryRepo.ExistsByUsername(ctx, "exists")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("用户名不存在", func(t *testing.T) {
		db := setupTestDB(t)
		queryRepo := NewUserQueryRepository(db)

		exists, err := queryRepo.ExistsByUsername(ctx, "notexists")

		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestUserQueryRepository_ExistsByEmail(t *testing.T) { //nolint:dupl // 测试代码结构相似是可接受的
	ctx := context.Background()

	t.Run("邮箱存在", func(t *testing.T) {
		db := setupTestDB(t)
		cmdRepo := NewUserCommandRepository(db)
		queryRepo := NewUserQueryRepository(db)

		u := &user.User{
			Username: "emailtest",
			Email:    "check@example.com",
			Password: "password",
			Status:   "active",
		}
		err := cmdRepo.Create(ctx, u)
		require.NoError(t, err)

		exists, err := queryRepo.ExistsByEmail(ctx, "check@example.com")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("邮箱不存在", func(t *testing.T) {
		db := setupTestDB(t)
		queryRepo := NewUserQueryRepository(db)

		exists, err := queryRepo.ExistsByEmail(ctx, "notcheck@example.com")

		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestUserQueryRepository_List(t *testing.T) {
	ctx := context.Background()

	t.Run("分页查询用户列表", func(t *testing.T) {
		db := setupTestDB(t)
		cmdRepo := NewUserCommandRepository(db)
		queryRepo := NewUserQueryRepository(db)

		// 创建多个用户
		for i := 1; i <= 5; i++ {
			u := &user.User{
				Username: "user" + string(rune('0'+i)),
				Email:    "user" + string(rune('0'+i)) + "@example.com",
				Password: "password",
				Status:   "active",
			}
			err := cmdRepo.Create(ctx, u)
			require.NoError(t, err)
		}

		// 查询前 3 个
		users, err := queryRepo.List(ctx, 0, 3)
		require.NoError(t, err)
		assert.Len(t, users, 3)

		// 查询后 2 个
		users, err = queryRepo.List(ctx, 3, 10)
		require.NoError(t, err)
		assert.Len(t, users, 2)
	})
}

func TestUserQueryRepository_Count(t *testing.T) {
	ctx := context.Background()

	t.Run("统计用户数量", func(t *testing.T) {
		db := setupTestDB(t)
		cmdRepo := NewUserCommandRepository(db)
		queryRepo := NewUserQueryRepository(db)

		// 创建 3 个用户
		for i := 1; i <= 3; i++ {
			u := &user.User{
				Username: "countuser" + string(rune('0'+i)),
				Email:    "count" + string(rune('0'+i)) + "@example.com",
				Password: "password",
				Status:   "active",
			}
			err := cmdRepo.Create(ctx, u)
			require.NoError(t, err)
		}

		count, err := queryRepo.Count(ctx)

		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestUserQueryRepository_GetByIDWithRoles(t *testing.T) {
	ctx := context.Background()

	t.Run("获取用户及其角色", func(t *testing.T) {
		db := setupTestDB(t)
		cmdRepo := NewUserCommandRepository(db)
		queryRepo := NewUserQueryRepository(db)

		// 创建角色
		role := createTestRole(t, db, "admin")

		// 创建用户
		u := &user.User{
			Username: "roleuser",
			Email:    "role@example.com",
			Password: "password",
			Status:   "active",
		}
		err := cmdRepo.Create(ctx, u)
		require.NoError(t, err)

		// 分配角色
		err = cmdRepo.AssignRoles(ctx, u.ID, []uint{role.ID})
		require.NoError(t, err)

		// 查询
		found, err := queryRepo.GetByIDWithRoles(ctx, u.ID)

		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Len(t, found.Roles, 1)
		assert.Equal(t, "admin", found.Roles[0].Name)
	})
}

func TestUserModel_Mapping(t *testing.T) {
	t.Run("Entity to Model", func(t *testing.T) {
		entity := &user.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed",
			FullName: "Test User",
			Status:   "active",
		}

		model := newUserModelFromEntity(entity)

		assert.Equal(t, entity.ID, model.ID)
		assert.Equal(t, entity.Username, model.Username)
		assert.Equal(t, entity.Email, model.Email)
		assert.Equal(t, entity.Password, model.Password)
	})

	t.Run("Model to Entity", func(t *testing.T) {
		model := &UserModel{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed",
			FullName: "Test User",
			Status:   "active",
		}

		entity := model.ToEntity()

		assert.Equal(t, model.ID, entity.ID)
		assert.Equal(t, model.Username, entity.Username)
		assert.Equal(t, model.Email, entity.Email)
	})

	t.Run("Nil handling", func(t *testing.T) {
		var nilEntity *user.User
		model := newUserModelFromEntity(nilEntity)
		assert.Nil(t, model)

		var nilModel *UserModel
		entity := nilModel.ToEntity()
		assert.Nil(t, entity)
	})
}
