package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/model"
)

type CommentRepo struct {
	*BaseRepo[model.Comment]
}

func NewCommentRepo(dbCli *gorm.DB) *CommentRepo {
	return &CommentRepo{
		BaseRepo: NewBaseRepo[model.Comment](dbCli),
	}
}

func WithCommentPreloads() []QueryOption {
	return []QueryOption{
		WithPreloads("User"),
		WithPreloads("User.Avatar"),
	}
}

func (r *CommentRepo) FindByID(ctx context.Context, id uint64) (*model.Comment, error) {
	return r.BaseRepo.FindByID(ctx, id, WithCommentPreloads()...)
}

func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *CommentRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

func (r *CommentRepo) FindByMomentID(ctx context.Context, momentID uint64, page, size int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Comment{}).Where("moment_id = ?", momentID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := db.Preload("User").Preload("User.Avatar").
		Order("created_at ASC").
		Offset(offset).Limit(size).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (r *CommentRepo) FindByParentID(ctx context.Context, parentID uint64) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Preload("User").Preload("User.Avatar").
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepo) FindByPathPrefix(ctx context.Context, pathPrefix string) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.db.WithContext(ctx).
		Where("path LIKE ?", pathPrefix+"%").
		Preload("User").Preload("User.Avatar").
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepo) CountByMomentID(ctx context.Context, momentID uint64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Comment{}).Where("moment_id = ?", momentID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CommentRepo) DeleteByMomentID(ctx context.Context, momentID uint64) error {
	return r.db.WithContext(ctx).Where("moment_id = ?", momentID).Delete(&model.Comment{}).Error
}

func (r *CommentRepo) Update(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *CommentRepo) FindByIDs(ctx context.Context, ids []uint64) ([]model.Comment, error) {
	if len(ids) == 0 {
		return []model.Comment{}, nil
	}
	var comments []model.Comment
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Preload("User").Preload("User.Avatar").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
