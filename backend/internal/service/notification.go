package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/model"
	"github.com/bookandmusic/love-girl/internal/repo"
)

type NotificationService struct {
	*BaseService
	NotificationRepo *repo.NotificationRepo
	FileService      *FileService
}

func NewNotificationService(log *log.Logger, notificationRepo *repo.NotificationRepo, fileService *FileService) *NotificationService {
	return &NotificationService{
		BaseService:      &BaseService{Log: log},
		NotificationRepo: notificationRepo,
		FileService:      fileService,
	}
}

type FrontendNotificationSender struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Avatar any    `json:"avatar,omitempty"`
}

type FrontendNotification struct {
	ID        uint64                     `json:"id"`
	Type      string                     `json:"type"`
	MomentID  uint64                     `json:"momentId"`
	CommentID uint64                     `json:"commentId"`
	Sender    FrontendNotificationSender `json:"sender"`
	Content   string                     `json:"content"`
	IsRead    bool                       `json:"isRead"`
	CreatedAt string                     `json:"createdAt"`
}

type NotificationListResponse struct {
	Notifications []*FrontendNotification `json:"notifications"`
	Total         int64                   `json:"total"`
	Page          int                     `json:"page"`
	Size          int                     `json:"size"`
}

func (s *NotificationService) convertToFrontendFormat(c *gin.Context, notification *model.Notification) *FrontendNotification {
	if notification == nil {
		return nil
	}

	sender := FrontendNotificationSender{}
	if notification.Sender != nil {
		sender.ID = notification.Sender.ID
		sender.Name = notification.Sender.Name
		sender.Avatar = s.FileService.BuildFileResponse(c, notification.Sender.Avatar)
	}

	return &FrontendNotification{
		ID:        notification.ID,
		Type:      string(notification.Type),
		MomentID:  notification.MomentID,
		CommentID: notification.CommentID,
		Sender:    sender,
		Content:   notification.Content,
		IsRead:    notification.IsRead,
		CreatedAt: notification.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID, senderID, momentID, commentID uint64, notificationType model.NotificationType, content string) error {
	notification := &model.Notification{
		UserID:    userID,
		SenderID:  senderID,
		MomentID:  momentID,
		CommentID: commentID,
		Type:      notificationType,
		Content:   content,
		IsRead:    false,
	}

	s.Log.Info("正在创建通知",
		"receiverID", userID,
		"senderID", senderID,
		"momentID", momentID,
		"commentID", commentID,
		"type", notificationType,
	)

	if err := s.NotificationRepo.Create(ctx, notification); err != nil {
		s.Log.Error("创建通知失败", "error", err)
		return fmt.Errorf("创建通知失败")
	}

	s.Log.Info("通知创建成功", "notificationID", notification.ID, "receiverID", userID)
	return nil
}

func (s *NotificationService) ListUnreadNotifications(c *gin.Context, userID uint64, page, size int) (*NotificationListResponse, error) {
	ctx := c.Request.Context()
	s.Log.Info("获取未读通知列表", "userID", userID, "page", page, "size", size)
	notifications, total, err := s.NotificationRepo.FindUnreadByUserID(ctx, userID, page, size)
	if err != nil {
		s.Log.Error("获取通知列表失败", "error", err, "userID", userID)
		return nil, fmt.Errorf("系统内部错误")
	}

	s.Log.Info("获取通知列表成功", "userID", userID, "total", total, "count", len(notifications))

	result := make([]*FrontendNotification, len(notifications))
	for i, notification := range notifications {
		result[i] = s.convertToFrontendFormat(c, &notification)
	}

	return &NotificationListResponse{
		Notifications: result,
		Total:         total,
		Page:          page,
		Size:          size,
	}, nil
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uint64) (int64, error) {
	return s.NotificationRepo.CountUnreadByUserID(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id uint64, userID uint64) error {
	return s.NotificationRepo.MarkAsRead(ctx, id)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uint64) error {
	return s.NotificationRepo.MarkAllAsRead(ctx, userID)
}
