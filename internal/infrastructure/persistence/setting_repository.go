package persistence

import (
	"context"
	"errors"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type settingRepository struct {
	db *gorm.DB
}

// NewSettingRepository 创建 Setting 仓储实例（deprecated，使用 NewSettingCommandRepository/NewSettingQueryRepository）
func NewSettingRepository(db *gorm.DB) setting.Repository {
	return &settingRepository{db: db}
}

// NewSettingCommandRepository 创建配置写操作仓储实例
func NewSettingCommandRepository(db *gorm.DB) setting.CommandRepository {
	return &settingRepository{db: db}
}

// NewSettingQueryRepository 创建配置读操作仓储实例
func NewSettingQueryRepository(db *gorm.DB) setting.QueryRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) Create(ctx context.Context, s *setting.Setting) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *settingRepository) Update(ctx context.Context, s *setting.Setting) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *settingRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&setting.Setting{}, id).Error
}

func (r *settingRepository) FindByID(ctx context.Context, id uint) (*setting.Setting, error) {
	var s setting.Setting
	err := r.db.WithContext(ctx).First(&s, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *settingRepository) FindByKey(ctx context.Context, key string) (*setting.Setting, error) {
	var s setting.Setting
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *settingRepository) FindByCategory(ctx context.Context, category string) ([]*setting.Setting, error) {
	var settings []*setting.Setting
	err := r.db.WithContext(ctx).Where("category = ?", category).Order("key ASC").Find(&settings).Error
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *settingRepository) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	var settings []*setting.Setting
	err := r.db.WithContext(ctx).Order("category ASC, key ASC").Find(&settings).Error
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *settingRepository) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	// 使用 GORM 的 Clauses 实现 Upsert (ON CONFLICT DO UPDATE)
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(settings).Error
}
