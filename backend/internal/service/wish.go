package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/model"
	"github.com/bookandmusic/love-girl/internal/repo"
)

type WishQueryParams struct {
	Page    int
	Size    int
	SortBy  string
	Order   string
	Filters []repo.FilterCondition
}

// Wish 祝福结构（与前端保持一致的命名）
type Wish struct {
	ID         uint64    `json:"id"`
	Content    string    `json:"content"`
	AuthorName string    `json:"authorName"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"createdAt"`
	Approved   bool      `json:"approved"`
}

// WishListResponse 祝福列表响应
type WishListResponse struct {
	Wishs      []*Wish `json:"wishes"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	Total      int64   `json:"total"`
	TotalPages int     `json:"totalPages"`
	TotalCount int64   `json:"totalCount"`
}

// WishCreateRequest 创建祝福请求
type WishCreateRequest struct {
	Content    string `json:"content" binding:"required,max=500"`
	AuthorName string `json:"authorName" binding:"required,max=100"`
	Email      string `json:"email" binding:"omitempty,email,max=150"`
}

// WishService 祝福服务
type WishService struct {
	*BaseService
	WishRepo *repo.WishRepo
}

// NewWishService 创建祝福服务实例
func NewWishService(log *log.Logger, wishRepo *repo.WishRepo) *WishService {
	return &WishService{
		BaseService: &BaseService{Log: log},
		WishRepo:    wishRepo,
	}
}

// 将model.Wish转换为前端期望的格式
func (s *WishService) convertToWish(wish *model.Wish) *Wish {
	if wish == nil {
		return nil
	}

	return &Wish{
		ID:         wish.ID,
		Content:    wish.Content,
		AuthorName: wish.AuthorName,
		Email:      wish.Email,
		CreatedAt:  wish.CreatedAt,
		Approved:   wish.Approved,
	}
}

// ListWishs 获取祝福列表
func (s *WishService) ListWishs(ctx context.Context, page, size int, approved *bool) (*WishListResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	} else if size > 100 {
		size = 100
	}

	var opts []repo.QueryOption
	if approved != nil {
		opts = append(opts, repo.WithConditions(repo.FilterCondition{Field: "approved", Operator: "eq", Value: *approved}))
	}

	wishesData, total, err := s.WishRepo.FindWithPagination(ctx, page, size, opts...)
	if err != nil {
		s.Log.Error("获取祝福列表失败", "error", err, "page", page, "size", size, "approved", approved)
		return nil, fmt.Errorf("系统内部错误")
	}

	totalPages := int((total + int64(size) - 1) / int64(size))

	wish := make([]*Wish, len(wishesData))
	for i, wishItem := range wishesData {
		wish[i] = s.convertToWish(&wishItem)
	}

	return &WishListResponse{
		Wishs:      wish,
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
		TotalCount: total,
	}, nil
}

// ListWishsWithQuery 根据查询参数获取祝福列表
func (s *WishService) ListWishsWithQuery(ctx context.Context, params *WishQueryParams, approved *bool) (*WishListResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Size < 1 {
		params.Size = 10
	} else if params.Size > 100 {
		params.Size = 100
	}

	var opts []repo.QueryOption
	if approved != nil {
		opts = append(opts, repo.WithConditions(repo.FilterCondition{Field: "approved", Operator: "eq", Value: *approved}))
	}
	if len(params.Filters) > 0 {
		opts = append(opts, repo.WithConditions(params.Filters...))
	}
	if params.SortBy != "" {
		opts = append(opts, repo.WithOrder(params.SortBy, params.Order == "desc"))
	}

	wishesData, total, err := s.WishRepo.FindWithPagination(ctx, params.Page, params.Size, opts...)
	if err != nil {
		s.Log.Error("获取祝福列表失败", "error", err, "params", params, "approved", approved)
		return nil, fmt.Errorf("系统内部错误")
	}

	totalPages := int((total + int64(params.Size) - 1) / int64(params.Size))

	wish := make([]*Wish, len(wishesData))
	for i, wishItem := range wishesData {
		wish[i] = s.convertToWish(&wishItem)
	}

	return &WishListResponse{
		Wishs:      wish,
		Page:       params.Page,
		Size:       params.Size,
		Total:      total,
		TotalPages: totalPages,
		TotalCount: total,
	}, nil
}

// CreateWish 创建祝福
func (s *WishService) CreateWish(ctx context.Context, req *WishCreateRequest) (*Wish, error) {
	// 创建祝福模型
	wish := &model.Wish{
		Content:    req.Content,
		AuthorName: req.AuthorName,
		Email:      req.Email,
		CreatedAt:  time.Now(),
		Approved:   false, // 默认未审核
	}

	// 保存祝福
	if err := s.WishRepo.Create(ctx, wish); err != nil {
		s.Log.Error("创建祝福失败", "error", err, "request", req)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToWish(wish), nil
}

// DeleteWish 删除祝福
func (s *WishService) DeleteWish(ctx context.Context, id uint64) (bool, error) {
	// 检查祝福是否存在
	_, err := s.WishRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("祝福不存在", "id", id)
			return false, nil
		}
		s.Log.Error("查询祝福失败", "error", err, "id", id)
		return false, fmt.Errorf("系统内部错误")
	}

	// 删除祝福
	if err := s.WishRepo.DeleteByID(ctx, id); err != nil {
		s.Log.Error("删除祝福失败", "error", err, "id", id)
		return false, fmt.Errorf("系统内部错误")
	}

	return true, nil
}

// ApproveWish 批准祝福
func (s *WishService) ApproveWish(ctx context.Context, id uint64) error {
	// 检查祝福是否存在
	_, err := s.WishRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("祝福不存在", "id", id)
			return nil
		}
		s.Log.Error("查询祝福失败", "error", err, "id", id)
		return fmt.Errorf("系统内部错误")
	}

	// 批准祝福
	if err := s.WishRepo.UpdateApprovalStatus(ctx, id, true); err != nil {
		s.Log.Error("批准祝福失败", "error", err, "id", id)
		return fmt.Errorf("系统内部错误")
	}

	return nil
}
