package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/model"
)

type SettingRepo struct {
	*BaseRepo[model.Setting]
}

func NewSettingRepo(db *gorm.DB) *SettingRepo {
	return &SettingRepo{
		BaseRepo: NewBaseRepo[model.Setting](db),
	}
}

// GetSettingByKey 根据键名获取设置
func (r *SettingRepo) GetSettingByKey(ctx context.Context, key string) (*model.Setting, error) {
	opts := []QueryOption{
		WithConditions(FilterCondition{
			Field:    "key",
			Operator: "eq",
			Value:    key,
		}),
	}
	entity, err := r.FindOne(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// GetSettingsByGroup 根据分组获取设置列表
func (r *SettingRepo) GetSettingsByGroup(ctx context.Context, group string) ([]model.Setting, error) {
	opts := []QueryOption{
		WithConditions(FilterCondition{
			Field:    "group",
			Operator: "eq",
			Value:    group,
		}),
		WithOrder("key", false),
	}
	list, err := r.List(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateOrCreateByKey 根据键名更新或创建设置
func (r *SettingRepo) UpdateOrCreateByKey(ctx context.Context, key string, setting *model.Setting) error {
	existing, err := r.GetSettingByKey(ctx, key)
	if err != nil {
		// 如果记录不存在，创建新记录
		if err == gorm.ErrRecordNotFound {
			setting.Key = key
			return r.Create(ctx, setting)
		}
		return err
	}

	// 如果记录存在，更新现有记录（使用 ID 作为条件，排除 key 以避免唯一约束问题）
	return r.db.WithContext(ctx).Model(&model.Setting{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
		"value": setting.Value,
		"type":  setting.Type,
		"label": setting.Label,
		"group": setting.Group,
	}).Error
}
