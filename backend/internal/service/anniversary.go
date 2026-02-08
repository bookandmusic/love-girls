package service

import (
	"context"
	"fmt"

	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/model"
	"github.com/bookandmusic/love-girl/internal/repo"
)

// FrontendAnniversary 前端期望的Anniversary数据结构
type FrontendAnniversary struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Date        string `json:"date"` // MM-DD
	Description string `json:"description"`
	Calendar    string `json:"calendar"` // solar/lunar
}

// 将model.Anniversary转换为前端期望的格式
func (s *AnniversaryService) convertToFrontendFormat(anniversary *model.Anniversary) *FrontendAnniversary {
	if anniversary == nil {
		return nil
	}

	// 保持完整日期格式YYYY-MM-DD
	date := anniversary.Date

	return &FrontendAnniversary{
		ID:          anniversary.ID,
		Title:       anniversary.Title,
		Date:        date,
		Description: anniversary.Description,
		Calendar:    anniversary.Calendar,
	}
}

type AnniversaryService struct {
	*BaseService
	AnniversaryRepo *repo.AnniversaryRepo
}

func NewAnniversaryService(log *log.Logger, anniversaryRepo *repo.AnniversaryRepo) *AnniversaryService {
	return &AnniversaryService{
		BaseService:     &BaseService{Log: log},
		AnniversaryRepo: anniversaryRepo,
	}
}

// AnniversaryCreateRequest 创建纪念日请求
type AnniversaryCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Date        string `json:"date" binding:"required"` // YYYY-MM-DD格式
	Description string `json:"description"`
	Calendar    string `json:"calendar" binding:"required,oneof=solar lunar"`
}

// AnniversaryUpdateRequest 更新纪念日请求
type AnniversaryUpdateRequest struct {
	Title       *string `json:"title"`
	Date        *string `json:"date"` // YYYY-MM-DD格式
	Description *string `json:"description"`
	Calendar    *string `json:"calendar" binding:"omitempty,oneof=solar lunar"`
}

// AnniversaryListResponse 纪念日列表响应
type AnniversaryListResponse struct {
	Anniversaries []*FrontendAnniversary `json:"anniversaries"`
	Page          int                    `json:"page"`
	Size          int                    `json:"size"`
	Total         int64                  `json:"total"`
	TotalPages    int                    `json:"totalPages"`
	TotalCount    int64                  `json:"totalCount"`
}

// ListAnniversaries 获取纪念日列表
func (s *AnniversaryService) ListAnniversaries(ctx context.Context, page, size int) (*AnniversaryListResponse, error) {
	// 获取所有纪念日数据，不进行分页
	anniversaries, err := s.AnniversaryRepo.List(ctx)
	if err != nil {
		s.Log.Error("获取纪念日列表失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 转换为前端格式
	frontendAnniversaries := make([]*FrontendAnniversary, len(anniversaries))
	for i, anniversary := range anniversaries {
		frontendAnniversaries[i] = s.convertToFrontendFormat(&anniversary)
	}

	// 计算总数和总页数（总页数为1，因为返回所有数据）
	totalCount := int64(len(anniversaries))
	totalPages := 1

	return &AnniversaryListResponse{
		Anniversaries: frontendAnniversaries,
		Page:          page, // 保持传入的页码不变
		Size:          size, // 保持传入的每页数量不变
		Total:         totalCount,
		TotalPages:    totalPages, // 总页数为1，因为返回所有数据
		TotalCount:    totalCount,
	}, nil
}

// CreateAnniversary 创建纪念日
func (s *AnniversaryService) CreateAnniversary(ctx context.Context, req *AnniversaryCreateRequest) (*FrontendAnniversary, error) {
	// 直接使用YYYY-MM-DD格式的日期
	anniversary := &model.Anniversary{
		Title:       req.Title,
		Date:        req.Date,
		Description: req.Description,
		Calendar:    req.Calendar,
	}

	if err := s.AnniversaryRepo.BaseRepo.Create(ctx, anniversary); err != nil {
		s.Log.Error("创建纪念日失败", "error", err, "title", req.Title)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToFrontendFormat(anniversary), nil
}

// UpdateAnniversary 更新纪念日
func (s *AnniversaryService) UpdateAnniversary(ctx context.Context, id uint64, req *AnniversaryUpdateRequest) (*FrontendAnniversary, error) {
	// 首先查询纪念日是否存在
	anniversary, err := s.AnniversaryRepo.FindByID(ctx, id)
	if err != nil {
		s.Log.Error("查询纪念日失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	if anniversary == nil {
		return nil, fmt.Errorf("纪念日不存在")
	}

	// 更新字段
	if req.Title != nil {
		anniversary.Title = *req.Title
	}

	if req.Date != nil {
		// 直接使用YYYY-MM-DD格式的日期
		anniversary.Date = *req.Date
	}

	if req.Description != nil {
		anniversary.Description = *req.Description
	}

	if req.Calendar != nil {
		anniversary.Calendar = *req.Calendar
	}

	// 保存更新
	if err := s.AnniversaryRepo.BaseRepo.Update(ctx, anniversary); err != nil {
		s.Log.Error("更新纪念日失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToFrontendFormat(anniversary), nil
}

// DeleteAnniversary 删除纪念日
func (s *AnniversaryService) DeleteAnniversary(ctx context.Context, id uint64) error {
	// 首先查询纪念日是否存在
	anniversary, err := s.AnniversaryRepo.FindByID(ctx, id)
	if err != nil {
		s.Log.Error("查询纪念日失败", "error", err, "id", id)
		return fmt.Errorf("系统内部错误")
	}

	if anniversary == nil {
		return fmt.Errorf("纪念日不存在")
	}

	// 删除纪念日
	if err := s.AnniversaryRepo.BaseRepo.DeleteByID(ctx, id); err != nil {
		s.Log.Error("删除纪念日失败", "error", err, "id", id)
		return fmt.Errorf("系统内部错误")
	}

	return nil
}
