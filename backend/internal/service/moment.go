package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/model"
	"github.com/bookandmusic/love-girl/internal/repo"
)

// FrontendMoment 前端期望的Moment数据结构
type FrontendMoment struct {
	ID        uint64          `json:"id"`
	Content   string          `json:"content"`
	Images    []FrontendPhoto `json:"images"`
	Likes     int             `json:"likes"`
	CreatedAt string          `json:"createdAt"`
	Author    FrontendAuthor  `json:"author"`
	IsPublic  bool            `json:"isPublic"`
}

// FrontendPhoto 前端期望的Photo数据结构
type FrontendPhoto struct {
	ID       uint64        `json:"id"`
	MomentID uint64        `json:"momentId"`
	File     *FileResponse `json:"file"`
}

// FrontendAuthor 前端期望的Author数据结构
type FrontendAuthor struct {
	Name   string        `json:"name"`
	Avatar *FileResponse `json:"avatar"`
}

// 将model.Moment转换为前端期望的格式
func (s *MomentService) convertToFrontendFormat(ctx context.Context, moment *model.Moment) *FrontendMoment {
	if moment == nil {
		return nil
	}

	// 转换图片数据
	photos := make([]FrontendPhoto, 0, len(moment.EntityFiles))
	for _, ef := range moment.EntityFiles {
		if ef.File == nil {
			continue
		}
		photos = append(photos, FrontendPhoto{
			ID:       ef.File.ID,
			MomentID: ef.EntityID,
			File:     s.FileService.BuildFileResponse(ctx, ef.File),
		})
	}

	// 转换作者数据
	author := FrontendAuthor{}
	if moment.User != nil {
		author.Name = moment.User.Name
		author.Avatar = s.FileService.BuildFileResponse(ctx, moment.User.Avatar)
	}

	return &FrontendMoment{
		ID:        moment.ID,
		Content:   moment.Content,
		Images:    photos,
		Likes:     moment.Likes,
		CreatedAt: moment.CreatedAt.Format("2006-01-02 15:04:05"),
		Author:    author,
		IsPublic:  moment.IsPublic,
	}
}

type MomentService struct {
	*BaseService
	MomentRepo  *repo.MomentRepo
	FileService *FileService
}

func NewMomentService(log *log.Logger, momentRepo *repo.MomentRepo, fileService *FileService) *MomentService {
	return &MomentService{
		BaseService: &BaseService{Log: log},
		MomentRepo:  momentRepo,
		FileService: fileService,
	}
}

// MomentCreateRequest 创建动态请求
type MomentCreateRequest struct {
	Content  string   `json:"content" binding:"required"`
	ImageIds []uint64 `json:"imageIds"`
	IsPublic bool     `json:"isPublic"`
	UserID   uint64   `json:"userId" binding:"required,gt=0"`
}

// MomentUpdateRequest 更新动态请求
type MomentUpdateRequest struct {
	Content  *string  `json:"content"`
	ImageIds []uint64 `json:"imageIds"`
	IsPublic *bool    `json:"isPublic"`
}

// MomentPublicRequest 动态公开状态请求
type MomentPublicRequest struct {
	IsPublic bool `json:"isPublic"`
}

// MomentLikeResponse 点赞响应
type MomentLikeResponse struct {
	Likes int `json:"likes"`
}

// MomentListResponse 动态列表响应
type MomentListResponse struct {
	Moments    []*FrontendMoment `json:"moments"`
	Page       int               `json:"page"`
	Size       int               `json:"size"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"totalPages"`
}

// CreateMoment 创建动态
func (s *MomentService) CreateMoment(ctx context.Context, req *MomentCreateRequest) (*FrontendMoment, error) {
	moment := &model.Moment{
		Content:  req.Content,
		IsPublic: req.IsPublic,
		UserID:   req.UserID,
	}

	// 使用事务创建动态和文件关联
	if err := s.MomentRepo.CreateWithFiles(ctx, moment, req.ImageIds); err != nil {
		s.Log.Error("创建动态失败", "error", err, "content", req.Content)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 重新查询动态以获取完整信息
	createdMoment, err := s.MomentRepo.FindByID(ctx, moment.ID)
	if err != nil {
		s.Log.Error("查询刚创建的动态失败", "error", err, "momentId", moment.ID)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToFrontendFormat(ctx, createdMoment), nil
}

// ListMoments 获取动态列表（仅返回公开的动态）
func (s *MomentService) ListMoments(ctx context.Context, page, size int) (*MomentListResponse, error) {
	moments, total, err := s.MomentRepo.ListMoments(ctx, page, size,
		repo.FilterCondition{Field: "is_public", Operator: "eq", Value: true},
	)
	if err != nil {
		s.Log.Error("获取动态列表失败", "error", err, "page", page, "size", size)
		return nil, fmt.Errorf("系统内部错误")
	}

	totalPage := int((total + int64(size) - 1) / int64(size))

	// 转换为前端期望的格式
	frontendMoments := make([]*FrontendMoment, len(moments))
	for i, moment := range moments {
		momentPtr := &moment
		frontendMoments[i] = s.convertToFrontendFormat(ctx, momentPtr)
	}

	return &MomentListResponse{
		Moments:    frontendMoments,
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPage,
	}, nil
}

// ListMomentsByAuthStatus 根据认证状态获取动态列表
// 如果用户已登录，返回当前用户的所有动态；否则，只返回公开的动态
func (s *MomentService) ListMomentsByAuthStatus(ctx context.Context, page, size int, isLoggedIn bool, userID uint64) (*MomentListResponse, error) {
	var (
		moments []model.Moment
		total   int64
		err     error
	)

	if isLoggedIn {
		// 用户已登录，只返回当前用户的所有动态（包括公开和私有的）
		moments, total, err = s.MomentRepo.ListMoments(ctx, page, size,
			repo.FilterCondition{Field: "user_id", Operator: "eq", Value: userID},
		)
	} else {
		// 用户未登录，只返回公开的动态
		moments, total, err = s.MomentRepo.ListMoments(ctx, page, size,
			repo.FilterCondition{Field: "is_public", Operator: "eq", Value: true},
		)
	}

	if err != nil {
		s.Log.Error("获取动态列表失败", "error", err, "page", page, "size", size, "isLoggedIn", isLoggedIn)
		return nil, fmt.Errorf("系统内部错误")
	}

	totalPage := int((total + int64(size) - 1) / int64(size))

	// 转换为前端期望的格式
	frontendMoments := make([]*FrontendMoment, len(moments))
	for i, moment := range moments {
		momentPtr := &moment
		frontendMoments[i] = s.convertToFrontendFormat(ctx, momentPtr)
	}

	return &MomentListResponse{
		Moments:    frontendMoments,
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPage,
	}, nil
}

// UpdateMoment 更新动态
func (s *MomentService) UpdateMoment(ctx context.Context, id uint64, req *MomentUpdateRequest) (*FrontendMoment, error) {
	moment, err := s.MomentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("动态不存在", "id", id)
			return nil, nil
		}
		s.Log.Error("查询动态失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 更新动态内容
	if req.Content != nil {
		moment.Content = *req.Content
	}
	if req.IsPublic != nil {
		moment.IsPublic = *req.IsPublic
	}

	// 更新动态信息
	if err := s.MomentRepo.UpdateWithFiles(ctx, moment, req.ImageIds); err != nil {
		s.Log.Error("更新动态失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	updatedMoment, err := s.MomentRepo.FindByID(ctx, id)
	if err != nil {
		s.Log.Error("查询更新后的动态失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToFrontendFormat(ctx, updatedMoment), nil
}

// DeleteMoment 删除动态
func (s *MomentService) DeleteMoment(ctx context.Context, id uint64) (bool, error) {
	_, err := s.MomentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("动态不存在", "id", id)
			return false, nil
		}
		s.Log.Error("查询动态失败", "error", err, "id", id)
		return false, fmt.Errorf("系统内部错误")
	}

	if err := s.MomentRepo.DeleteWithFiles(ctx, id); err != nil {
		s.Log.Error("删除动态失败", "error", err, "id", id)
		return false, fmt.Errorf("系统内部错误")
	}

	return true, nil
}

// UpdatePublicStatus 更新动态公开状态
func (s *MomentService) UpdatePublicStatus(ctx context.Context, id uint64, status bool) (*FrontendMoment, error) {
	_, err := s.MomentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("动态不存在", "id", id)
			return nil, nil
		}
		s.Log.Error("查询动态失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	if err := s.MomentRepo.UpdatePublicStatus(ctx, id, status); err != nil {
		s.Log.Error("更新动态公开状态失败", "error", err, "id", id, "status", status)
		return nil, fmt.Errorf("系统内部错误")
	}

	updatedMoment, err := s.MomentRepo.FindByID(ctx, id)
	if err != nil {
		s.Log.Error("查询更新后的动态失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToFrontendFormat(ctx, updatedMoment), nil
}

// LikeMoment 点赞动态
func (s *MomentService) LikeMoment(ctx context.Context, id uint64) (*FrontendMoment, error) {
	_, err := s.MomentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("动态不存在", "id", id)
			return nil, nil
		}
		s.Log.Error("查询动态失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	if err := s.MomentRepo.UpdateLike(ctx, id); err != nil {
		s.Log.Error("更新动态点赞数失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	likedMoment, err := s.MomentRepo.FindByID(ctx, id)
	if err != nil {
		s.Log.Error("查询点赞后的动态失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToFrontendFormat(ctx, likedMoment), nil
}
