package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/model"
)

// AlbumRepo 专辑仓库，提供相册相关的数据操作
// 功能：
//   - 创建相册（使用 Create）
//   - 删除相册（使用 DeleteByID）
//   - 查询相册（支持单查、列表、分页）
//   - 更新相册信息（使用 Update）
//   - 设置封面图片
//   - 管理相册图片（添加、删除、列表查询、数量统计）
type AlbumRepo struct {
	*BaseRepo[model.Album]
	entityFileRepo *EntityFileRepo
}

// NewAlbumRepo 创建新的专辑仓库实例
func NewAlbumRepo(dbCli *gorm.DB) *AlbumRepo {
	return &AlbumRepo{
		BaseRepo:       NewBaseRepo[model.Album](dbCli),
		entityFileRepo: NewEntityFileRepo(dbCli),
	}
}

// FindByID 根据ID查找相册并预加载关联数据
// 参数：
//   - ctx: 上下文
//   - id: 相册ID
//
// 返回：相册实体、错误
func (r *AlbumRepo) FindByID(ctx context.Context, id uint64) (*model.Album, error) {
	return r.BaseRepo.FindByID(ctx, id,
		WithPreloads("CoverImage", "EntityFiles.File"),
		WithPreloadCond("EntityFiles", "entity_type = ?", "album"),
	)
}

// ListAlbums 分页展示所有相册，需要关联封面图片
// 参数：
//   - ctx: 上下文
//   - page: 页码，从1开始
//   - size: 每页数量
//   - opts: 查询选项（排序、过滤等）
//
// 返回：相册列表、总数、错误
func (r *AlbumRepo) ListAlbums(ctx context.Context, page, size int, opts ...QueryOption) ([]model.Album, int64, error) {
	// 确保预加载关联数据
	allOpts := append(opts,
		WithPreloads("CoverImage", "EntityFiles.File"),
		WithPreloadCond("EntityFiles", "entity_type = ?", "album"),
	)
	return r.BaseRepo.FindWithPagination(ctx, page, size, allOpts...)
}

// SetCoverImage 设置相册封面图片
// 参数：
//   - ctx: 上下文
//   - albumID: 相册ID
//   - imageID: 图片ID
//
// 返回：如果设置成功返回nil，否则返回错误
func (r *AlbumRepo) SetCoverImage(ctx context.Context, albumID uint64, imageID uint64) error {
	album, err := r.BaseRepo.FindByID(ctx, albumID)
	if err != nil {
		return err
	}

	album.CoverImageID = &imageID
	return r.BaseRepo.Update(ctx, album)
}

// ListAlbumPhotos 分页展示相册下的所有图片
// 参数：
//   - ctx: 上下文
//   - albumID: 相册ID
//   - page: 页码，从1开始
//   - size: 每页数量
//   - opts: 查询选项（排序等）
//
// 返回：图片列表、总数、错误
func (r *AlbumRepo) ListAlbumPhotos(ctx context.Context, albumID uint64, page, size int, opts ...QueryOption) ([]model.File, int64, error) {
	return r.entityFileRepo.GetAssociatedFilesWithPagination(ctx, albumID, "album", page, size, opts...)
}

// AssociateFiles 将文件列表关联到相册（事务）
// 参数：
//   - ctx: 上下文
//   - albumID: 相册ID
//   - fileIDs: 文件ID列表
//
// 返回：如果关联成功返回nil，否则返回错误
// 说明：此方法会替换相册当前的所有文件关联
func (r *AlbumRepo) AssociateFiles(ctx context.Context, albumID uint64, fileIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除旧的关联
		if err := tx.Where("entity_id = ? AND entity_type = ?", albumID, "album").Delete(&model.EntityFile{}).Error; err != nil {
			return err
		}
		// 创建新的关联
		if len(fileIDs) > 0 {
			var associations []model.EntityFile
			for _, fileID := range fileIDs {
				associations = append(associations, model.EntityFile{
					EntityID:   albumID,
					EntityType: "album",
					FileID:     fileID,
				})
			}
			if err := tx.Create(&associations).Error; err != nil {
				return err
			}
		}
		// 更新相册照片数量
		if err := tx.Model(&model.Album{}).Where("id = ?", albumID).Update("photo_count", len(fileIDs)).Error; err != nil {
			return err
		}
		return nil
	})
}

// AppendFiles 将文件列表追加到相册（事务）
// 参数：
//   - ctx: 上下文
//   - albumID: 相册ID
//   - fileIDs: 文件ID列表
//
// 返回：如果关联成功返回nil，否则返回错误
// 说明：此方法只添加新的文件关联，不删除已存在的关联
func (r *AlbumRepo) AppendFiles(ctx context.Context, albumID uint64, fileIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 查询已存在的文件ID
		var existingFileIDs []uint64
		if err := tx.Model(&model.EntityFile{}).
			Where("entity_id = ? AND entity_type = ?", albumID, "album").
			Pluck("file_id", &existingFileIDs).Error; err != nil {
			return err
		}

		// 创建一个map用于快速查找已存在的文件ID
		existingIDMap := make(map[uint64]bool)
		for _, id := range existingFileIDs {
			existingIDMap[id] = true
		}

		// 只添加不存在的文件
		var newAssociations []model.EntityFile
		newCount := 0
		for _, fileID := range fileIDs {
			if !existingIDMap[fileID] {
				newAssociations = append(newAssociations, model.EntityFile{
					EntityID:   albumID,
					EntityType: "album",
					FileID:     fileID,
				})
				newCount++
			}
		}

		// 创建新的关联
		if len(newAssociations) > 0 {
			if err := tx.Create(&newAssociations).Error; err != nil {
				return err
			}
		}

		// 更新相册照片数量（增量更新）
		if newCount > 0 {
			if err := tx.Model(&model.Album{}).Where("id = ?", albumID).
				Update("photo_count", gorm.Expr("photo_count + ?", newCount)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// RemovePhoto 从相册删除某个图片（事务）
//
// 返回：如果删除成功返回nil，否则返回错误
func (r *AlbumRepo) RemovePhoto(ctx context.Context, albumID uint64, fileID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先获取相册信息，检查是否要清除封面
		var album model.Album
		if err := tx.Where("id = ?", albumID).First(&album).Error; err != nil {
			return err
		}

		// 删除照片关联
		if err := tx.Where("entity_id = ? AND entity_type = ? AND file_id = ?", albumID, "album", fileID).Delete(&model.EntityFile{}).Error; err != nil {
			return err
		}

		// 更新相册照片数量（减1）
		if err := tx.Model(&model.Album{}).Where("id = ?", albumID).UpdateColumn("photo_count", gorm.Expr("CASE WHEN photo_count > 0 THEN photo_count - 1 ELSE 0 END")).Error; err != nil {
			return err
		}

		// 如果删除的是封面照片，清除封面
		if album.CoverImageID != nil && *album.CoverImageID == fileID {
			if err := tx.Model(&model.Album{}).Where("id = ?", albumID).Update("cover_image_id", nil).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetAlbumCount 查询所有相册的数量
// 参数：
//   - ctx: 上下文
//
// 返回：相册数量、错误
func (r *AlbumRepo) GetAlbumCount(ctx context.Context) (int, error) {
	count, err := r.BaseRepo.Count(ctx)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetPhotoCount 查询相册下图片的数量
// 参数：
//   - ctx: 上下文
//   - albumID: 相册ID，为0时返回所有相册的照片总数
//
// 返回：图片数量、错误
func (r *AlbumRepo) GetPhotoCount(ctx context.Context, albumID ...uint64) (int, error) {
	conditions := []FilterCondition{
		{Field: "entity_type", Operator: "eq", Value: "album"},
	}

	// 如果提供了albumID且不为0，则只统计该相册的照片
	if len(albumID) > 0 && albumID[0] > 0 {
		conditions = append(conditions, FilterCondition{
			Field: "entity_id", Operator: "eq", Value: albumID[0],
		})
	}

	entityFiles, err := r.entityFileRepo.List(ctx, WithConditions(conditions...))
	if err != nil {
		return 0, err
	}
	return len(entityFiles), nil
}
