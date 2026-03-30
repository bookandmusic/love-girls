package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/model"
)

type NotificationRepo struct {
	*BaseRepo[model.Notification]
}

func NewNotificationRepo(dbCli *gorm.DB) *NotificationRepo {
	return &NotificationRepo{
		BaseRepo: NewBaseRepo[model.Notification](dbCli),
	}
}

func WithNotificationPreloads() []QueryOption {
	return []QueryOption{
		WithPreloads("Sender"),
		WithPreloads("Sender.Avatar"),
	}
}

func (r *NotificationRepo) FindByID(ctx context.Context, id uint64) (*model.Notification, error) {
	return r.BaseRepo.FindByID(ctx, id, WithNotificationPreloads()...)
}

func (r *NotificationRepo) Create(ctx context.Context, notification *model.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *NotificationRepo) FindUnreadByUserID(ctx context.Context, userID uint64, page, size int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := db.Preload("Sender").Preload("Sender.Avatar").
		Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *NotificationRepo) CountUnreadByUserID(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *NotificationRepo) MarkAsRead(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *NotificationRepo) MarkAllAsRead(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}
