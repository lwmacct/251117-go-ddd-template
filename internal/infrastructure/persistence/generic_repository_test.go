package persistence

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupGenericTestDB 创建用于泛型测试的 SQLite in-memory 数据库。
func setupGenericTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "无法连接测试数据库")

	// 迁移角色和权限表
	err = db.AutoMigrate(&RoleModel{}, &PermissionModel{})
	require.NoError(t, err, "数据库迁移失败")

	return db
}

// TestGenericCommandRepository_Create 测试泛型写仓储的 Create 方法
func TestGenericCommandRepository_Create(t *testing.T) {
	ctx := context.Background()
	db := setupGenericTestDB(t)

	// 使用 Role 作为测试实体（因为它的 CRUD 完全由泛型提供）
	repo := NewRoleCommandRepository(db)

	t.Run("成功创建实体", func(t *testing.T) {
		r := &role.Role{
			Name:        "test-role",
			DisplayName: "Test Role",
			Description: "A test role for generic repository",
		}

		err := repo.Create(ctx, r)

		require.NoError(t, err)
		assert.NotZero(t, r.ID, "ID 应该被自动设置")
		assert.NotZero(t, r.CreatedAt, "CreatedAt 应该被自动设置")

		// 验证数据库记录
		var model RoleModel
		err = db.First(&model, r.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "test-role", model.Name)
		assert.Equal(t, "Test Role", model.DisplayName)
	})

	t.Run("唯一性约束违反", func(t *testing.T) {
		// 先创建一个角色
		r1 := &role.Role{
			Name:        "unique-role",
			DisplayName: "Unique Role",
		}
		err := repo.Create(ctx, r1)
		require.NoError(t, err)

		// 尝试创建同名角色
		r2 := &role.Role{
			Name:        "unique-role", // 相同的 name
			DisplayName: "Another Role",
		}
		err = repo.Create(ctx, r2)

		assert.Error(t, err, "应该因唯一性约束而失败")
	})
}

// TestGenericCommandRepository_Update 测试泛型写仓储的 Update 方法
func TestGenericCommandRepository_Update(t *testing.T) {
	ctx := context.Background()
	db := setupGenericTestDB(t)
	repo := NewRoleCommandRepository(db)

	t.Run("成功更新实体", func(t *testing.T) {
		// 先创建
		r := &role.Role{
			Name:        "update-test",
			DisplayName: "Original Name",
			Description: "Original description",
		}
		err := repo.Create(ctx, r)
		require.NoError(t, err)

		// 更新
		r.DisplayName = "Updated Name"
		r.Description = "Updated description"
		err = repo.Update(ctx, r)

		require.NoError(t, err)
		assert.NotZero(t, r.UpdatedAt, "UpdatedAt 应该被更新")

		// 验证数据库
		var model RoleModel
		err = db.First(&model, r.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", model.DisplayName)
		assert.Equal(t, "Updated description", model.Description)
	})
}

// TestGenericCommandRepository_Delete 测试泛型写仓储的 Delete 方法
func TestGenericCommandRepository_Delete(t *testing.T) {
	ctx := context.Background()
	db := setupGenericTestDB(t)
	repo := NewRoleCommandRepository(db)

	t.Run("成功软删除实体", func(t *testing.T) {
		// 先创建
		r := &role.Role{
			Name:        "delete-test",
			DisplayName: "To Delete",
		}
		err := repo.Create(ctx, r)
		require.NoError(t, err)
		id := r.ID

		// 删除
		err = repo.Delete(ctx, id)
		require.NoError(t, err)

		// 验证软删除 (不使用 Unscoped 应该找不到)
		var model RoleModel
		err = db.First(&model, id).Error
		require.ErrorIs(t, err, gorm.ErrRecordNotFound, "软删除后不应该能找到记录")

		// 使用 Unscoped 可以找到
		err = db.Unscoped().First(&model, id).Error
		require.NoError(t, err)
		assert.NotNil(t, model.DeletedAt.Time, "DeletedAt 应该被设置")
	})

	t.Run("删除不存在的实体不报错", func(t *testing.T) {
		err := repo.Delete(ctx, 99999)
		// GORM 的 Delete 对不存在的记录不会报错
		assert.NoError(t, err)
	})
}

// TestGenericQueryRepository_GetByID 测试泛型读仓储的 GetByID 方法
func TestGenericQueryRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	db := setupGenericTestDB(t)
	cmdRepo := NewRoleCommandRepository(db)
	queryRepo := NewRoleQueryRepository(db)

	t.Run("成功获取存在的实体", func(t *testing.T) {
		// 先创建
		r := &role.Role{
			Name:        "get-test",
			DisplayName: "Get Test",
		}
		err := cmdRepo.Create(ctx, r)
		require.NoError(t, err)

		// 获取
		found, err := queryRepo.FindByID(ctx, r.ID)

		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, r.ID, found.ID)
		assert.Equal(t, "get-test", found.Name)
		assert.Equal(t, "Get Test", found.DisplayName)
	})

	t.Run("获取不存在的实体返回 nil", func(t *testing.T) {
		// Role 的 FindByID 对不存在的记录返回 nil, nil
		found, err := queryRepo.FindByID(ctx, 99999)

		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

// TestGenericQueryRepository_Exists 测试泛型读仓储的 Exists 方法
func TestGenericQueryRepository_Exists(t *testing.T) {
	ctx := context.Background()
	db := setupGenericTestDB(t)
	cmdRepo := NewRoleCommandRepository(db)
	queryRepo := NewRoleQueryRepository(db)

	t.Run("存在的实体返回 true", func(t *testing.T) {
		r := &role.Role{
			Name:        "exists-test",
			DisplayName: "Exists Test",
		}
		err := cmdRepo.Create(ctx, r)
		require.NoError(t, err)

		exists, err := queryRepo.Exists(ctx, r.ID)

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("不存在的实体返回 false", func(t *testing.T) {
		exists, err := queryRepo.Exists(ctx, 99999)

		require.NoError(t, err)
		assert.False(t, exists)
	})
}

// TestPermissionCommandRepository_GenericCRUD 验证 Permission 完全使用泛型
func TestPermissionCommandRepository_GenericCRUD(t *testing.T) {
	ctx := context.Background()
	db := setupGenericTestDB(t)
	repo := NewPermissionCommandRepository(db)

	t.Run("Permission 使用泛型 Create/Update/Delete", func(t *testing.T) {
		// Create
		p := &role.Permission{
			Code:        "test:action:create",
			Domain:      "test",
			Resource:    "action",
			Action:      "create",
			Description: "Test permission for create",
		}
		err := repo.Create(ctx, p)
		require.NoError(t, err)
		assert.NotZero(t, p.ID)

		// Update
		p.Description = "Updated description"
		err = repo.Update(ctx, p)
		require.NoError(t, err)

		// 验证更新
		var model PermissionModel
		err = db.First(&model, p.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "Updated description", model.Description)

		// Delete
		err = repo.Delete(ctx, p.ID)
		require.NoError(t, err)

		// 验证软删除
		err = db.First(&model, p.ID).Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
