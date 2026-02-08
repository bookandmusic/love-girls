package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/model"
)

// MomentRepo 瞬间/动态仓库，提供动态相关的数据操作
// 功能：
//   - 创建动态（使用 CreateWithFiles）
//   - 删除动态（同时删除文件关联）
//   - 查询动态（支持单查）
//   - 更新动态（支持同时修改关联图片）
//   - 更新点赞数
//   - 更新公开状态
type MomentRepo struct {
	*BaseRepo[model.Moment]
	entityFileRepo *EntityFileRepo
}

const MomentEntityType = "moment"

func NewMomentRepo(dbCli *gorm.DB) *MomentRepo {
	return &MomentRepo{
		BaseRepo:       NewBaseRepo[model.Moment](dbCli),
		entityFileRepo: NewEntityFileRepo(dbCli),
	}
}

func WithMomentPreloads() []QueryOption {
	return []QueryOption{
		WithPreloads("User"),
		WithPreloads("EntityFiles"),
		WithPreloadCond("EntityFiles", "entity_type = ?", MomentEntityType),
		WithPreloads("EntityFiles.File"),
	}
}

// FindByID 根据ID查找动态并预加载关联数据
//
// 返回：动态实体、错误
func (r *MomentRepo) FindByID(ctx context.Context, id uint64) (*model.Moment, error) {
	return r.BaseRepo.FindByID(ctx, id, WithMomentPreloads()...)
}

// DeleteWithFiles 删除动态，需要删除动态和图片文件关联关系（事务）
//
// 返回：如果删除成功返回nil，否则返回错误
func (r *MomentRepo) DeleteWithFiles(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除动态与文件的关联
		if err := tx.Where("entity_id = ? AND entity_type = ?", id, "moment").Delete(&model.EntityFile{}).Error; err != nil {
			return err
		}
		// 再删除动态记录
		return tx.Delete(&model.Moment{}, id).Error
	})
}

// UpdateWithFiles 修改动态，支持同时修改关联图片（事务）
//
// 返回：如果更新成功返回nil，否则返回错误
func (r *MomentRepo) UpdateWithFiles(ctx context.Context, moment *model.Moment, fileIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新动态信息
		if err := tx.Save(moment).Error; err != nil {
			return err
		}
		// 删除旧的文件关联
		if err := tx.Where("entity_id = ? AND entity_type = ?", moment.ID, "moment").Delete(&model.EntityFile{}).Error; err != nil {
			return err
		}
		// 创建新的文件关联
		if len(fileIDs) > 0 {
			var associations []model.EntityFile
			for _, fileID := range fileIDs {
				associations = append(associations, model.EntityFile{
					EntityID:   moment.ID,
					EntityType: "moment",
					FileID:     fileID,
				})
			}
			if err := tx.Create(&associations).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// CreateWithFiles 创建动态并关联文件（事务）
//
// 返回：如果创建成功返回nil，否则返回错误
func (r *MomentRepo) CreateWithFiles(ctx context.Context, moment *model.Moment, fileIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建动态记录
		if err := tx.Create(moment).Error; err != nil {
			return err
		}
		// 创建文件关联
		if len(fileIDs) > 0 {
			var associations []model.EntityFile
			for _, fileID := range fileIDs {
				associations = append(associations, model.EntityFile{
					EntityID:   moment.ID,
					EntityType: "moment",
					FileID:     fileID,
				})
			}
			if err := tx.Create(&associations).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateLike 更新点赞数
//
// 返回：如果更新成功返回nil，否则返回错误
func (r *MomentRepo) UpdateLike(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Moment{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error
}

// UpdatePublicStatus 更新公开状态
//
// 返回：如果更新成功返回nil，否则返回错误
func (r *MomentRepo) UpdatePublicStatus(ctx context.Context, id uint64, status bool) error {
	moment, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	moment.IsPublic = status
	return r.BaseRepo.Update(ctx, moment)
}

// ListMoments 分页查询动态列表
func (r *MomentRepo) ListMoments(ctx context.Context, page, size int, conditions ...FilterCondition) ([]model.Moment, int64, error) {
	preloadOps := WithMomentPreloads()
	ops := append([]QueryOption{WithConditions(conditions...)}, preloadOps...)
	return r.BaseRepo.FindWithPagination(ctx, page, size, ops...)
}
