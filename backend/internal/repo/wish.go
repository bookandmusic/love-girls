package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/model"
)

// WishRepo 愿望仓库
type WishRepo struct {
	*BaseRepo[model.Wish]
}

// NewWishRepo 创建新的愿望仓库实例
func NewWishRepo(dbCli *gorm.DB) *WishRepo {
	return &WishRepo{
		BaseRepo: NewBaseRepo[model.Wish](dbCli),
	}
}

// UpdateApprovalStatus 更新愿望审核状态
func (r *WishRepo) UpdateApprovalStatus(ctx context.Context, id uint64, approved bool) error {
	wish, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	wish.Approved = approved
	return r.BaseRepo.Update(ctx, wish)
}
