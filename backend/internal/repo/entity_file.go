package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/model"
)

// EntityFileRepo 实体文件关联仓库
type EntityFileRepo struct {
	*BaseRepo[model.EntityFile]
}

// NewEntityFileRepo 创建新的实体文件关联仓库实例
func NewEntityFileRepo(dbCli *gorm.DB) *EntityFileRepo {
	return &EntityFileRepo{
		BaseRepo: NewBaseRepo[model.EntityFile](dbCli),
	}
}

// AssociateFile 关联文件到实体
func (r *EntityFileRepo) AssociateFile(ctx context.Context, entityID uint64, entityType string, fileID uint64) error {
	entityFile := model.EntityFile{
		EntityID:   entityID,
		EntityType: entityType,
		FileID:     fileID,
	}

	return r.Create(ctx, &entityFile)
}

// BatchAssociateFiles 批量关联文件到实体
func (r *EntityFileRepo) BatchAssociateFiles(ctx context.Context, entityID uint64, entityType string, fileIDs []uint64) error {
	var associations []model.EntityFile
	for _, fileID := range fileIDs {
		associations = append(associations, model.EntityFile{
			EntityID:   entityID,
			EntityType: entityType,
			FileID:     fileID,
		})
	}

	return r.BaseRepo.DB().WithContext(ctx).Create(&associations).Error
}

// RemoveAssociation 移除实体与文件的关联
func (r *EntityFileRepo) RemoveAssociation(ctx context.Context, entityID uint64, entityType string, fileID uint64) error {
	var entityFile model.EntityFile
	db := r.BaseRepo.DB().WithContext(ctx)
	db = r.ApplyFilters(db, []FilterCondition{
		{Field: "entity_id", Operator: "eq", Value: entityID},
		{Field: "entity_type", Operator: "eq", Value: entityType},
		{Field: "file_id", Operator: "eq", Value: fileID},
	})
	return db.Delete(&entityFile).Error
}

// GetAssociatedFiles 获取实体关联的所有文件
func (r *EntityFileRepo) GetAssociatedFiles(ctx context.Context, entityID uint64, entityType string) ([]model.File, error) {
	var files []model.File
	err := r.BaseRepo.DB().WithContext(ctx).
		Joins("JOIN entity_files ON entity_files.file_id = files.id").
		Where("entity_files.entity_id = ? AND entity_files.entity_type = ?", entityID, entityType).
		Find(&files).Error

	return files, err
}

// GetAssociatedFilesWithPagination 分页获取实体关联的文件
func (r *EntityFileRepo) GetAssociatedFilesWithPagination(ctx context.Context, entityID uint64, entityType string, page, size int, opts ...QueryOption) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	query := r.BaseRepo.DB().WithContext(ctx).
		Model(&model.File{}).
		Joins("JOIN entity_files ON entity_files.file_id = files.id").
		Where("entity_files.entity_id = ? AND entity_files.entity_type = ? AND entity_files.deleted_at IS NULL", entityID, entityType)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply order from options
	options := &QueryOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if options.OrderBy != "" {
		if options.Desc {
			query = query.Order(options.OrderBy + " DESC")
		} else {
			query = query.Order(options.OrderBy + " ASC")
		}
	}

	// Pagination
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// GetAssociatedEntities 获取文件关联的所有实体
func (r *EntityFileRepo) GetAssociatedEntities(ctx context.Context, fileID uint64, entityType string) ([]model.EntityFile, error) {
	return r.BaseRepo.List(ctx, WithConditions(
		FilterCondition{Field: "file_id", Operator: "eq", Value: fileID},
		FilterCondition{Field: "entity_type", Operator: "eq", Value: entityType},
	))
}

// CreateBatch 批量创建关联记录
func (r *EntityFileRepo) CreateBatch(ctx context.Context, entityFiles []model.EntityFile) error {
	return r.BaseRepo.DB().WithContext(ctx).Create(&entityFiles).Error
}

// ListByEntity 查询实体的所有关联
func (r *EntityFileRepo) ListByEntity(ctx context.Context, entityID uint64, entityType string) ([]model.EntityFile, error) {
	return r.BaseRepo.List(ctx, WithConditions(
		FilterCondition{Field: "entity_id", Operator: "eq", Value: entityID},
		FilterCondition{Field: "entity_type", Operator: "eq", Value: entityType},
	))
}

// DeleteByEntity 删除实体的所有关联
func (r *EntityFileRepo) DeleteByEntity(ctx context.Context, entityID uint64, entityType string) error {
	var entityFile model.EntityFile
	db := r.BaseRepo.DB().WithContext(ctx)
	db = r.ApplyFilters(db, []FilterCondition{
		{Field: "entity_id", Operator: "eq", Value: entityID},
		{Field: "entity_type", Operator: "eq", Value: entityType},
	})
	return db.Delete(&entityFile).Error
}

// UpdateEntityFiles 替换实体的文件关联
func (r *EntityFileRepo) UpdateEntityFiles(ctx context.Context, entityID uint64, entityType string, fileIDs []uint64) error {
	// 先删除旧的关联
	if err := r.DeleteByEntity(ctx, entityID, entityType); err != nil {
		return err
	}

	// 创建新的关联
	return r.BatchAssociateFiles(ctx, entityID, entityType, fileIDs)
}

// CheckFileUsed 检查文件是否被使用
func (r *EntityFileRepo) CheckFileUsed(ctx context.Context, fileID uint64) (bool, error) {
	var count int64
	db := r.BaseRepo.DB().WithContext(ctx)
	db = r.ApplyFilters(db, []FilterCondition{
		{Field: "file_id", Operator: "eq", Value: fileID},
	})
	if err := db.Model(&model.EntityFile{}).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
