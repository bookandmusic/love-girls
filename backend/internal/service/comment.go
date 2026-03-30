package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/model"
	"github.com/bookandmusic/love-girl/internal/repo"
)

type CommentService struct {
	*BaseService
	CommentRepo      *repo.CommentRepo
	MomentRepo       *repo.MomentRepo
	NotificationRepo *repo.NotificationRepo
	FileService      *FileService
	NotificationSvc  *NotificationService
}

func NewCommentService(log *log.Logger, commentRepo *repo.CommentRepo, momentRepo *repo.MomentRepo, notificationRepo *repo.NotificationRepo, fileService *FileService, notificationService *NotificationService) *CommentService {
	return &CommentService{
		BaseService:      &BaseService{Log: log},
		CommentRepo:      commentRepo,
		MomentRepo:       momentRepo,
		NotificationRepo: notificationRepo,
		FileService:      fileService,
		NotificationSvc:  notificationService,
	}
}

type FrontendAuthorInfo struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Avatar any    `json:"avatar,omitempty"`
}

type FrontendComment struct {
	ID        uint64              `json:"id"`
	Content   string              `json:"content"`
	MomentID  uint64              `json:"momentId"`
	ParentID  *uint64             `json:"parentId"`
	ReplyToID *uint64             `json:"replyToId"`
	ReplyTo   *FrontendAuthorInfo `json:"replyTo,omitempty"`
	UserID    uint64              `json:"userId"`
	Author    FrontendAuthor      `json:"author"`
	Depth     int                 `json:"depth"`
	CreatedAt string              `json:"createdAt"`
	Children  []*FrontendComment  `json:"children,omitempty"`
}

type CommentCreateRequest struct {
	Content   string  `json:"content" binding:"required"`
	MomentID  uint64  `json:"momentId"`
	ParentID  *uint64 `json:"parentId"`
	ReplyToID *uint64 `json:"replyToId"`
	UserID    uint64  `json:"userId"`
}

type CommentListResponse struct {
	Comments []*FrontendComment `json:"comments"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	Size     int                `json:"size"`
}

func (s *CommentService) convertToFrontendFormat(c *gin.Context, comment *model.Comment, replyToUsersMap map[uint64]*model.User) *FrontendComment {
	if comment == nil {
		return nil
	}

	author := FrontendAuthor{}
	if comment.User != nil {
		author.ID = comment.User.ID
		author.Name = comment.User.Name
		author.Avatar = s.FileService.BuildFileResponse(c, comment.User.Avatar)
	}

	var replyTo *FrontendAuthorInfo
	if comment.ReplyToID != nil {
		if replyToUser, ok := replyToUsersMap[*comment.ReplyToID]; ok && replyToUser != nil {
			replyTo = &FrontendAuthorInfo{
				ID:     replyToUser.ID,
				Name:   replyToUser.Name,
				Avatar: s.FileService.BuildFileResponse(c, replyToUser.Avatar),
			}
		}
	}

	return &FrontendComment{
		ID:        comment.ID,
		Content:   comment.Content,
		MomentID:  comment.MomentID,
		ParentID:  comment.ParentID,
		ReplyToID: comment.ReplyToID,
		ReplyTo:   replyTo,
		UserID:    comment.UserID,
		Author:    author,
		Depth:     comment.Depth,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *CommentService) CreateComment(c *gin.Context, req *CommentCreateRequest) (*FrontendComment, error) {
	ctx := c.Request.Context()

	comment := &model.Comment{
		Content:   req.Content,
		MomentID:  req.MomentID,
		ParentID:  req.ParentID,
		ReplyToID: req.ReplyToID,
		UserID:    req.UserID,
	}

	if req.ParentID != nil {
		parentComment, err := s.CommentRepo.FindByID(ctx, *req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("父评论不存在")
			}
			return nil, fmt.Errorf("系统内部错误")
		}
		comment.Depth = parentComment.Depth + 1
		comment.Path = parentComment.Path + "/" + strconv.FormatUint(comment.ID, 10)
	} else {
		comment.Depth = 0
	}

	if err := s.CommentRepo.Create(ctx, comment); err != nil {
		s.Log.Error("创建评论失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	if comment.ParentID == nil {
		comment.Path = strconv.FormatUint(comment.ID, 10)
		if err := s.CommentRepo.Update(ctx, comment); err != nil {
			s.Log.Error("更新评论路径失败", "error", err)
		}
	} else {
		comment.Path = comment.Path[:strings.LastIndex(comment.Path, "/")+1] + strconv.FormatUint(comment.ID, 10)
		if err := s.CommentRepo.Update(ctx, comment); err != nil {
			s.Log.Error("更新评论路径失败", "error", err)
		}
	}

	createdComment, err := s.CommentRepo.FindByID(ctx, comment.ID)
	if err != nil {
		s.Log.Error("查询刚创建的评论失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	replyToUsersMap := make(map[uint64]*model.User)
	if req.ReplyToID != nil {
		replyToComment, err := s.CommentRepo.FindByID(ctx, *req.ReplyToID)
		if err == nil && replyToComment.User != nil {
			replyToUsersMap[*req.ReplyToID] = replyToComment.User
		}
	}

	// 创建通知
	s.createNotification(ctx, req, createdComment)

	return s.convertToFrontendFormat(c, createdComment, replyToUsersMap), nil
}

func (s *CommentService) createNotification(ctx context.Context, req *CommentCreateRequest, comment *model.Comment) {
	// 获取动态信息以获取动态作者
	moment, err := s.MomentRepo.FindByID(ctx, req.MomentID)
	if err != nil {
		s.Log.Error("获取动态信息失败", "error", err, "momentID", req.MomentID)
		return
	}

	s.Log.Info("创建通知检查",
		"momentUserID", moment.UserID,
		"commentUserID", req.UserID,
		"momentID", req.MomentID,
		"commentID", comment.ID,
		"parentID", req.ParentID,
		"replyToID", req.ReplyToID,
	)

	// 不给自己发通知
	if moment.UserID == req.UserID {
		s.Log.Info("跳过通知：评论自己的动态", "userID", req.UserID)
		return
	}

	var notificationType model.NotificationType
	var receiverID uint64

	if req.ParentID == nil {
		// 一级评论：通知动态作者
		notificationType = model.NotificationTypeComment
		receiverID = moment.UserID
		s.Log.Info("一级评论：通知动态作者", "receiverID", receiverID)
	} else {
		// 回复评论：通知被回复的评论作者
		notificationType = model.NotificationTypeReply
		if req.ReplyToID != nil {
			replyToComment, err := s.CommentRepo.FindByID(ctx, *req.ReplyToID)
			if err != nil || replyToComment == nil {
				s.Log.Error("找不到被回复的评论", "replyToID", *req.ReplyToID, "error", err)
				return
			}
			receiverID = replyToComment.UserID
			s.Log.Info("回复评论：通知被回复者", "receiverID", receiverID)
		} else {
			parentComment, err := s.CommentRepo.FindByID(ctx, *req.ParentID)
			if err != nil || parentComment == nil {
				s.Log.Error("找不到父评论", "parentID", *req.ParentID, "error", err)
				return
			}
			receiverID = parentComment.UserID
			s.Log.Info("回复评论：通知父评论作者", "receiverID", receiverID)
		}
	}

	// 不给自己发通知
	if receiverID == req.UserID {
		s.Log.Info("跳过通知：回复自己的评论", "userID", req.UserID)
		return
	}

	// 截取评论内容作为通知内容
	content := req.Content
	if len(content) > 100 {
		content = content[:100] + "..."
	}

	s.Log.Info("创建通知", "receiverID", receiverID, "senderID", req.UserID, "type", notificationType)
	if err := s.NotificationSvc.CreateNotification(ctx, receiverID, req.UserID, req.MomentID, comment.ID, notificationType, content); err != nil {
		s.Log.Error("创建通知失败", "error", err)
	}
}

func (s *CommentService) DeleteComment(ctx context.Context, id uint64, userID uint64) (bool, error) {
	comment, err := s.CommentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("评论不存在", "id", id)
			return false, nil
		}
		s.Log.Error("查询评论失败", "error", err, "id", id)
		return false, fmt.Errorf("系统内部错误")
	}

	if comment.UserID != userID {
		return false, fmt.Errorf("无权限删除此评论")
	}

	if err := s.CommentRepo.Delete(ctx, id); err != nil {
		s.Log.Error("删除评论失败", "error", err, "id", id)
		return false, fmt.Errorf("系统内部错误")
	}

	return true, nil
}

func (s *CommentService) ListComments(c *gin.Context, momentID uint64, page, size int) (*CommentListResponse, error) {
	ctx := c.Request.Context()
	comments, total, err := s.CommentRepo.FindByMomentID(ctx, momentID, page, size)
	if err != nil {
		s.Log.Error("获取评论列表失败", "error", err, "momentID", momentID)
		return nil, fmt.Errorf("系统内部错误")
	}

	replyToIDs := make([]uint64, 0)
	for _, comment := range comments {
		if comment.ReplyToID != nil {
			replyToIDs = append(replyToIDs, *comment.ReplyToID)
		}
	}

	replyToUsersMap := make(map[uint64]*model.User)
	if len(replyToIDs) > 0 {
		replyToComments, err := s.CommentRepo.FindByIDs(ctx, replyToIDs)
		if err != nil {
			s.Log.Error("查询回复目标评论失败", "error", err)
		} else {
			for _, rc := range replyToComments {
				if rc.User != nil {
					replyToUsersMap[rc.ID] = rc.User
				}
			}
		}
	}

	allComments := make([]*FrontendComment, len(comments))
	for i, comment := range comments {
		allComments[i] = s.convertToFrontendFormat(c, &comment, replyToUsersMap)
	}

	tree := s.buildCommentTree(allComments)

	return &CommentListResponse{
		Comments: tree,
		Total:    total,
		Page:     page,
		Size:     size,
	}, nil
}

func (s *CommentService) buildCommentTree(comments []*FrontendComment) []*FrontendComment {
	commentMap := make(map[uint64]*FrontendComment)
	var rootComments []*FrontendComment

	for _, comment := range comments {
		commentMap[comment.ID] = comment
	}

	for _, comment := range comments {
		if comment.ParentID == nil {
			rootComments = append(rootComments, comment)
		} else {
			if parent, ok := commentMap[*comment.ParentID]; ok {
				if parent.Children == nil {
					parent.Children = []*FrontendComment{}
				}
				parent.Children = append(parent.Children, comment)
			}
		}
	}

	return rootComments
}

func (s *CommentService) CountByMomentID(ctx context.Context, momentID uint64) (int64, error) {
	return s.CommentRepo.CountByMomentID(ctx, momentID)
}
