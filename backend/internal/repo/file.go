package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/model"
)

// FileRepo 文件仓库，提供文件相关的数据操作
// 功能：
//   - 创建文件记录（使用 Create）
//   - 删除文件（使用 DeleteByID）
//   - 查询未使用的文件（用于清理）
//   - 根据哈希值查找文件（用于去重）
type FileRepo struct {
	*BaseRepo[model.File]
}

// NewFileRepo 创建新的文件仓库实例
func NewFileRepo(dbCli *gorm.DB) *FileRepo {
	return &FileRepo{
		BaseRepo: NewBaseRepo[model.File](dbCli),
	}
}

// FindByHash 根据哈希值查找文件
// 参数：
//   - ctx: 上下文
//   - hash: 文件哈希值
//
// 返回：文件实体、错误
func (r *FileRepo) FindByHash(ctx context.Context, hash string) (*model.File, error) {
	return r.BaseRepo.FindOne(ctx, WithConditions(
		FilterCondition{Field: "hash", Operator: "eq", Value: hash},
	))
}

// FindUnusedFiles 分页展示所有未使用的文件（没有在entity_files表中被引用的文件）
// 参数：
//   - ctx: 上下文
//   - page: 页码，从1开始
//   - size: 每页数量
//
// 返回：文件列表、总数、错误
func (r *FileRepo) FindUnusedFiles(ctx context.Context, page, size int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	// 构建查询，查找没有在entity_files表中出现的文件
	db := r.BaseRepo.DB().WithContext(ctx).Model(&model.File{}).
		Joins("LEFT JOIN entity_files ON files.id = entity_files.file_id").
		Where("entity_files.file_id IS NULL")

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}
