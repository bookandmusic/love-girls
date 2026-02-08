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

// AlbumCoverImage 相册封面图片结构
type AlbumCoverImage struct {
	ID      uint64        `json:"id"`
	AlbumID uint64        `json:"albumId,omitempty"`
	File    *FileResponse `json:"file"`
}

// Album 相册结构
type Album struct {
	ID          uint64           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	CoverImage  *AlbumCoverImage `json:"coverImage,omitempty"`
	PhotoCount  int              `json:"photoCount"`
	CreatedAt   string           `json:"createdAt"`
}

// AlbumListResponse 相册列表响应
type AlbumListResponse struct {
	Albums     []*Album `json:"albums"`
	Page       int      `json:"page"`
	Size       int      `json:"size"`
	Total      int64    `json:"total"`
	TotalPages int      `json:"totalPages"`
	TotalCount int64    `json:"totalCount"`
}

// AlbumCreateRequest 创建相册请求
type AlbumCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// AlbumUpdateRequest 更新相册请求
type AlbumUpdateRequest struct {
	Name        *string `json:"name" binding:"required"`
	Description *string `json:"description"`
}

// AlbumPhoto 相册照片结构
type AlbumPhoto struct {
	ID        uint64        `json:"id"`
	AlbumID   uint64        `json:"albumId"`
	File      *FileResponse `json:"file"`
	Alt       string        `json:"alt,omitempty"`
	CreatedAt string        `json:"createdAt"`
}

// AlbumPhotoListResponse 相册照片列表响应
type AlbumPhotoListResponse struct {
	Photos     []*AlbumPhoto `json:"photos"`
	Page       int           `json:"page"`
	Size       int           `json:"size"`
	Total      int64         `json:"total"`
	TotalPages int           `json:"totalPages"`
	TotalCount int64         `json:"totalCount"`
}

// AlbumAddPhotosRequest 添加照片到相册请求
type AlbumAddPhotosRequest struct {
	PhotoIDs []uint64 `json:"photoIds" binding:"required"`
}

// AlbumSetCoverRequest 设置相册封面请求
type AlbumSetCoverRequest struct {
	PhotoID uint64 `json:"photoId" binding:"required"`
}

// AlbumService 相册服务
type AlbumService struct {
	*BaseService
	AlbumRepo   *repo.AlbumRepo
	FileService *FileService
}

// NewAlbumService 创建相册服务实例
func NewAlbumService(log *log.Logger, albumRepo *repo.AlbumRepo, fileService *FileService) *AlbumService {
	return &AlbumService{
		BaseService: &BaseService{Log: log},
		AlbumRepo:   albumRepo,
		FileService: fileService,
	}
}

// 将model.Album转换为前端响应格式
func (s *AlbumService) convertToAlbum(ctx context.Context, album *model.Album) *Album {
	if album == nil {
		return nil
	}

	response := &Album{
		ID:          album.ID,
		Name:        album.Name,
		Description: album.Description,
		PhotoCount:  album.PhotoCount,
		CreatedAt:   album.CreatedAt.Format("2006-01-02"),
	}

	// 转换封面图片
	if album.CoverImage != nil {
		response.CoverImage = &AlbumCoverImage{
			ID:      album.CoverImage.ID,
			AlbumID: album.ID,
			File:    s.FileService.BuildFileResponse(ctx, album.CoverImage),
		}
	}

	return response
}

// 将model.File转换为AlbumPhoto格式
func (s *AlbumService) convertToAlbumPhoto(ctx context.Context, file *model.File, albumID uint64) *AlbumPhoto {
	if file == nil {
		return nil
	}

	return &AlbumPhoto{
		ID:        file.ID,
		AlbumID:   albumID,
		File:      s.FileService.BuildFileResponse(ctx, file),
		Alt:       file.OriginalName,
		CreatedAt: file.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ListAlbums 获取相册列表
func (s *AlbumService) ListAlbums(ctx context.Context, page, size int) (*AlbumListResponse, error) {
	// 使用默认分页参数
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	} else if size > 100 {
		size = 100
	}

	// 获取相册列表
	albums, total, err := s.AlbumRepo.ListAlbums(ctx, page, size)
	if err != nil {
		s.Log.Error("获取相册列表失败", "error", err, "page", page, "size", size)
		return nil, fmt.Errorf("系统内部错误")
	}

	totalPages := int((total + int64(size) - 1) / int64(size))

	// 转换为响应格式
	responseAlbums := make([]*Album, len(albums))
	for i, album := range albums {
		albumPtr := &album
		responseAlbums[i] = s.convertToAlbum(ctx, albumPtr)
	}

	return &AlbumListResponse{
		Albums:     responseAlbums,
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
		TotalCount: total,
	}, nil
}

// CreateAlbum 创建相册
func (s *AlbumService) CreateAlbum(ctx context.Context, req *AlbumCreateRequest) (*Album, error) {
	// 创建相册模型
	album := &model.Album{
		Name:        req.Name,
		Description: req.Description,
		PhotoCount:  0,
	}

	// 保存相册
	if err := s.AlbumRepo.BaseRepo.Create(ctx, album); err != nil {
		s.Log.Error("创建相册失败", "error", err, "request", req)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToAlbum(ctx, album), nil
}

// UpdateAlbum 更新相册
func (s *AlbumService) UpdateAlbum(ctx context.Context, id uint64, req *AlbumUpdateRequest) (*Album, error) {
	// 检查相册是否存在
	album, err := s.AlbumRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("相册不存在", "id", id)
			return nil, nil
		}
		s.Log.Error("获取相册失败", "id", id, "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 更新字段
	if req.Name != nil {
		album.Name = *req.Name
	}
	if req.Description != nil {
		album.Description = *req.Description
	}

	// 保存更新
	if err := s.AlbumRepo.BaseRepo.Update(ctx, album); err != nil {
		s.Log.Error("更新相册失败", "id", id, "error", err, "request", req)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToAlbum(ctx, album), nil
}

// DeleteAlbum 删除相册
func (s *AlbumService) DeleteAlbum(ctx context.Context, id uint64) (bool, error) {
	// 检查相册是否存在
	album, err := s.AlbumRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("相册不存在", "id", id)
			return false, nil
		}
		s.Log.Error("获取相册失败", "id", id, "error", err)
		return false, fmt.Errorf("系统内部错误")
	}

	// 检查相册下是否有照片
	if album.PhotoCount > 0 {
		s.Log.Info("相册下有照片，不允许删除", "id", id, "photoCount", album.PhotoCount)
		return false, fmt.Errorf("相册下还有照片，请先删除所有照片后再删除相册")
	}

	// 删除相册
	if err := s.AlbumRepo.BaseRepo.DeleteByID(ctx, id); err != nil {
		s.Log.Error("删除相册失败", "id", id, "error", err)
		return false, fmt.Errorf("系统内部错误")
	}

	return true, nil
}

// ListAlbumPhotos 获取相册照片列表
func (s *AlbumService) ListAlbumPhotos(ctx context.Context, albumID uint64, page, size int) (*AlbumPhotoListResponse, error) {
	// 检查相册是否存在
	_, err := s.AlbumRepo.FindByID(ctx, albumID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("相册不存在", "id", albumID)
			return nil, nil
		}
		s.Log.Error("获取相册失败", "id", albumID, "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 使用默认分页参数
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	} else if size > 100 {
		size = 100
	}

	// 获取相册照片列表
	photos, total, err := s.AlbumRepo.ListAlbumPhotos(ctx, albumID, page, size)
	if err != nil {
		s.Log.Error("获取相册照片列表失败", "albumId", albumID, "error", err, "page", page, "size", size)
		return nil, fmt.Errorf("系统内部错误")
	}

	totalPages := int((total + int64(size) - 1) / int64(size))

	// 转换为响应格式
	responsePhotos := make([]*AlbumPhoto, len(photos))
	for i, photo := range photos {
		photoPtr := &photo
		responsePhotos[i] = s.convertToAlbumPhoto(ctx, photoPtr, albumID)
	}

	return &AlbumPhotoListResponse{
		Photos:     responsePhotos,
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
		TotalCount: total,
	}, nil
}

// AddPhotosToAlbum 添加照片到相册
func (s *AlbumService) AddPhotosToAlbum(ctx context.Context, albumID uint64, req *AlbumAddPhotosRequest) ([]*AlbumPhoto, error) {
	// 检查相册是否存在
	_, err := s.AlbumRepo.FindByID(ctx, albumID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("相册不存在", "id", albumID)
			return nil, nil
		}
		s.Log.Error("获取相册失败", "id", albumID, "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 追加照片到相册（只添加新的文件关联，不删除已存在的关联）
	if err := s.AlbumRepo.AppendFiles(ctx, albumID, req.PhotoIDs); err != nil {
		s.Log.Error("追加相册照片失败", "albumId", albumID, "error", err, "photoIds", req.PhotoIDs)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 查询指定的照片并返回
	var files []model.File
	for _, photoID := range req.PhotoIDs {
		file, err := s.FileService.GetFile(ctx, photoID)
		if err != nil {
			s.Log.Warn("获取照片失败", "photoId", photoID, "error", err)
			continue
		}
		files = append(files, *file)
	}

	// 转换为响应格式
	responsePhotos := make([]*AlbumPhoto, len(files))
	for i, photo := range files {
		photoPtr := &photo
		responsePhotos[i] = s.convertToAlbumPhoto(ctx, photoPtr, albumID)
	}

	return responsePhotos, nil
}

// SetAlbumCover 设置相册封面
func (s *AlbumService) SetAlbumCover(ctx context.Context, albumID uint64, req *AlbumSetCoverRequest) (*Album, error) {
	// 检查相册是否存在
	_, err := s.AlbumRepo.FindByID(ctx, albumID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("相册不存在", "id", albumID)
			return nil, nil
		}
		s.Log.Error("获取相册失败", "id", albumID, "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 设置封面图片
	if err := s.AlbumRepo.SetCoverImage(ctx, albumID, req.PhotoID); err != nil {
		s.Log.Error("设置相册封面失败", "albumId", albumID, "photoId", req.PhotoID, "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 获取更新后的相册
	updatedAlbum, err := s.AlbumRepo.FindByID(ctx, albumID)
	if err != nil {
		s.Log.Error("获取更新后的相册失败", "id", albumID, "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	return s.convertToAlbum(ctx, updatedAlbum), nil
}

// RemovePhotoFromAlbum 从相册删除照片
func (s *AlbumService) RemovePhotoFromAlbum(ctx context.Context, albumID uint64, photoID uint64) (bool, error) {
	// 检查相册是否存在
	_, err := s.AlbumRepo.FindByID(ctx, albumID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("相册不存在", "id", albumID)
			return false, nil
		}
		s.Log.Error("获取相册失败", "id", albumID, "error", err)
		return false, fmt.Errorf("系统内部错误")
	}

	// 从相册删除照片（会在事务中删除关联、更新照片数量、清除封面）
	if err := s.AlbumRepo.RemovePhoto(ctx, albumID, photoID); err != nil {
		s.Log.Error("从相册删除照片失败", "albumId", albumID, "photoId", photoID, "error", err)
		return false, fmt.Errorf("系统内部错误")
	}

	return true, nil
}
