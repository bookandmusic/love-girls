package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/model"
	"github.com/bookandmusic/love-girl/internal/repo"
)

type PlaceQueryParams struct {
	Page    int
	Size    int
	SortBy  string
	Order   string
	Filters []repo.FilterCondition
}

// PlaceImage 地点图片结构
type PlaceImage struct {
	ID      uint64        `json:"id"`
	PlaceID uint64        `json:"placeId,omitempty"`
	File    *FileResponse `json:"file"`
}

// PlaceResponse 单个地点响应
type PlaceResponse struct {
	ID          uint64      `json:"id"`
	Name        string      `json:"name"`
	Latitude    float64     `json:"latitude"`
	Longitude   float64     `json:"longitude"`
	Image       *PlaceImage `json:"image,omitempty"`
	Description string      `json:"description"`
	Date        string      `json:"date"`
}

// PlaceListResponse 地点列表响应
type PlaceListResponse struct {
	Places     []*PlaceResponse `json:"places"`
	Page       int              `json:"page"`
	Size       int              `json:"size"`
	Total      int64            `json:"total"`
	TotalPages int              `json:"totalPages"`
	TotalCount int64            `json:"totalCount"`
}

// PlaceCreateRequest 创建地点请求
type PlaceCreateRequest struct {
	Name        string      `json:"name" binding:"required"`
	Latitude    float64     `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude   float64     `json:"longitude" binding:"required,min=-180,max=180"`
	Description string      `json:"description"`
	Date        string      `json:"date" binding:"required"`
	Image       *PlaceImage `json:"image,omitempty"`
}

// PlaceUpdateRequest 更新地点请求
type PlaceUpdateRequest struct {
	Name        *string     `json:"name"`
	Latitude    *float64    `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude   *float64    `json:"longitude" binding:"omitempty,min=-180,max=180"`
	Description *string     `json:"description"`
	Date        *string     `json:"date"`
	Image       *PlaceImage `json:"image,omitempty"`
}

// PlaceService 地点服务
type PlaceService struct {
	*BaseService
	PlaceRepo   *repo.PlaceRepo
	FileService *FileService
}

// NewPlaceService 创建地点服务实例
func NewPlaceService(log *log.Logger, placeRepo *repo.PlaceRepo, fileService *FileService) *PlaceService {
	return &PlaceService{
		BaseService: &BaseService{Log: log},
		PlaceRepo:   placeRepo,
		FileService: fileService,
	}
}

// 将model.Place转换为前端响应格式
func (s *PlaceService) convertToResponse(c *gin.Context, place *model.Place) *PlaceResponse {
	if place == nil {
		return nil
	}

	response := &PlaceResponse{
		ID:          place.ID,
		Name:        place.Name,
		Latitude:    place.Latitude,
		Longitude:   place.Longitude,
		Description: place.Description,
		Date:        place.Date,
	}

	// 转换图片信息
	if place.Image != nil {
		response.Image = &PlaceImage{
			ID:      place.Image.ID,
			PlaceID: place.ID,
			File:    s.FileService.BuildFileResponse(c, place.Image),
		}
	}

	return response
}

// ListPlaces 获取所有地点（保持分页结构）
func (s *PlaceService) ListPlaces(c *gin.Context) (*PlaceListResponse, error) {
	ctx := c.Request.Context()
	places, total, err := s.PlaceRepo.ListPlaces(ctx, 1, 10000)
	if err != nil {
		s.Log.Error("获取地点列表失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	responsePlaces := make([]*PlaceResponse, len(places))
	for i, place := range places {
		placePtr := &place
		responsePlaces[i] = s.convertToResponse(c, placePtr)
	}

	return &PlaceListResponse{
		Places:     responsePlaces,
		Page:       1,
		Size:       int(total),
		Total:      total,
		TotalPages: 1,
		TotalCount: total,
	}, nil
}

// ListPlacesWithQuery 根据查询参数获取地点列表
func (s *PlaceService) ListPlacesWithQuery(c *gin.Context, params *PlaceQueryParams) (*PlaceListResponse, error) {
	ctx := c.Request.Context()
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Size < 1 {
		params.Size = 10
	} else if params.Size > 100 {
		params.Size = 100
	}

	var opts []repo.QueryOption
	if len(params.Filters) > 0 {
		opts = append(opts, repo.WithConditions(params.Filters...))
	}
	if params.SortBy != "" {
		opts = append(opts, repo.WithOrder(params.SortBy, params.Order == "desc"))
	}

	places, total, err := s.PlaceRepo.ListPlacesWithOpts(ctx, params.Page, params.Size, opts...)
	if err != nil {
		s.Log.Error("获取地点列表失败", "error", err, "params", params)
		return nil, fmt.Errorf("系统内部错误")
	}

	totalPages := int((total + int64(params.Size) - 1) / int64(params.Size))

	responsePlaces := make([]*PlaceResponse, len(places))
	for i, place := range places {
		placePtr := &place
		responsePlaces[i] = s.convertToResponse(c, placePtr)
	}

	return &PlaceListResponse{
		Places:     responsePlaces,
		Page:       params.Page,
		Size:       params.Size,
		Total:      total,
		TotalPages: totalPages,
		TotalCount: total,
	}, nil
}

// CreatePlace 创建地点
func (s *PlaceService) CreatePlace(c *gin.Context, req *PlaceCreateRequest) (*PlaceResponse, error) {
	ctx := c.Request.Context()
	// 创建地点模型
	place := &model.Place{
		Name:        req.Name,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Description: req.Description,
		Date:        req.Date,
	}

	// 如果有图片，设置图片ID
	if req.Image != nil {
		imageID := req.Image.ID
		place.ImageID = &imageID
	}

	// 保存地点
	if err := s.PlaceRepo.CreateWithImage(ctx, place); err != nil {
		s.Log.Error("创建地点失败", "error", err, "request", req)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 获取完整信息
	createdPlace, err := s.PlaceRepo.FindByID(ctx, place.ID, repo.WithPreloads("Image"))
	if err != nil {
		s.Log.Error("查询创建的地点失败", "error", err, "id", place.ID)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToResponse(c, createdPlace), nil
}

// UpdatePlace 更新地点
func (s *PlaceService) UpdatePlace(c *gin.Context, id uint64, req *PlaceUpdateRequest) (*PlaceResponse, error) {
	ctx := c.Request.Context()
	// 获取现有地点
	place, err := s.PlaceRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("地点不存在", "id", id)
			return nil, nil
		}
		s.Log.Error("查询地点失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 更新字段
	if req.Name != nil {
		place.Name = *req.Name
	}
	if req.Latitude != nil {
		place.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		place.Longitude = *req.Longitude
	}
	if req.Description != nil {
		place.Description = *req.Description
	}
	if req.Date != nil {
		place.Date = *req.Date
	}

	// 更新图片
	if req.Image != nil {
		imageID := req.Image.ID
		place.ImageID = &imageID
	} else if req.Image == nil {
		// 如果明确传递null，则删除图片关联
		place.ImageID = nil
	}

	// 保存更新
	if err := s.PlaceRepo.UpdateWithImage(ctx, place); err != nil {
		s.Log.Error("更新地点失败", "error", err, "id", id, "request", req)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 获取更新后的信息
	updatedPlace, err := s.PlaceRepo.FindByID(ctx, id, repo.WithPreloads("Image"))
	if err != nil {
		s.Log.Error("查询更新后的地点失败", "error", err, "id", id)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToResponse(c, updatedPlace), nil
}

// DeletePlace 删除地点
func (s *PlaceService) DeletePlace(ctx context.Context, id uint64) (bool, error) {
	// 检查地点是否存在
	_, err := s.PlaceRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("地点不存在", "id", id)
			return false, nil
		}
		s.Log.Error("查询地点失败", "error", err, "id", id)
		return false, fmt.Errorf("系统内部错误")
	}

	// 删除地点
	if err := s.PlaceRepo.DeleteWithImage(ctx, id); err != nil {
		s.Log.Error("删除地点失败", "error", err, "id", id)
		return false, fmt.Errorf("系统内部错误")
	}

	return true, nil
}
