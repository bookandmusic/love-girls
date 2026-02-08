package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/model"
)

// PlaceRepo 地点仓库，提供地点相关的数据操作
// 功能：
//   - 创建地点（可同时创建图片关联）
//   - 删除地点（同时删除文件关联）
//   - 查询地点（支持单查、列表、分页）
//   - 更新地点（支持同时修改关联图片）
//
// 说明：地点数据的图片是允许不存在的
type PlaceRepo struct {
	*BaseRepo[model.Place]
	entityFileRepo *EntityFileRepo
}

// NewPlaceRepo 创建新的地点仓库实例
func NewPlaceRepo(dbCli *gorm.DB) *PlaceRepo {
	return &PlaceRepo{
		BaseRepo:       NewBaseRepo[model.Place](dbCli),
		entityFileRepo: NewEntityFileRepo(dbCli),
	}
}

// ListPlaces 分页展示所有地点数据，需要关联每个地点的图片，使用外键关联的图片
// 参数：
//   - ctx: 上下文
//   - page: 页码，从1开始
//   - size: 每页数量
//   - opts: 查询选项（排序、过滤等）
//
// 返回：地点列表、总数、错误
func (r *PlaceRepo) ListPlaces(ctx context.Context, page, size int, opts ...QueryOption) ([]model.Place, int64, error) {
	// 确保预加载关联数据
	allOpts := append(opts, WithPreloads("Image"))
	return r.BaseRepo.FindWithPagination(ctx, page, size, allOpts...)
}

// CreateWithImage 创建地点记录，如果存在图片，需要同步添加一条关联关系
// 参数：
//   - ctx: 上下文
//   - place: 地点实体（包含ImageID字段）
//
// 返回：如果创建成功返回nil，否则返回错误
func (r *PlaceRepo) CreateWithImage(ctx context.Context, place *model.Place) error {
	// 创建地点记录
	if err := r.BaseRepo.Create(ctx, place); err != nil {
		return err
	}

	// 如果存在图片ID，创建关联关系，以便跟踪图片使用情况
	if place.ImageID != nil {
		return r.entityFileRepo.AssociateFile(ctx, place.ID, "place", *place.ImageID)
	}

	return nil
}

// UpdateWithImage 修改地点，支持同时修改关联图片关系
// 参数：
//   - ctx: 上下文
//   - place: 地点实体（包含ImageID字段）
//
// 返回：如果更新成功返回nil，否则返回错误
func (r *PlaceRepo) UpdateWithImage(ctx context.Context, place *model.Place) error {
	// 更新地点信息
	if err := r.BaseRepo.Update(ctx, place); err != nil {
		return err
	}

	// 更新图片关联关系，以便跟踪图片使用情况
	var fileIDs []uint64
	if place.ImageID != nil {
		fileIDs = []uint64{*place.ImageID}
	}
	return r.entityFileRepo.UpdateEntityFiles(ctx, place.ID, "place", fileIDs)
}

// DeleteWithImage 删除地点，需要删除地点和图片文件关联关系
// 参数：
//   - ctx: 上下文
//   - id: 地点ID
//
// 返回：如果删除成功返回nil，否则返回错误
func (r *PlaceRepo) DeleteWithImage(ctx context.Context, id uint64) error {
	// 先删除地点与文件的关联
	if err := r.entityFileRepo.DeleteByEntity(ctx, id, "place"); err != nil {
		return err
	}

	// 再删除地点记录
	return r.BaseRepo.DeleteByID(ctx, id)
}

// GetPlaceCount 查询所有地点的数量
// 参数：
//   - ctx: 上下文
//
// 返回：地点数量、错误
func (r *PlaceRepo) GetPlaceCount(ctx context.Context) (int, error) {
	count, err := r.BaseRepo.Count(ctx)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
