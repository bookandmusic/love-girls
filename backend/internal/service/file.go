package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/config"
	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/model"
	"github.com/bookandmusic/love-girl/internal/repo"
	"github.com/bookandmusic/love-girl/internal/storage"
)

type FileService struct {
	*BaseService
	Storage       storage.Storage
	FileRepo      repo.FileRepo
	serverCfg     *config.ServerConfig
	storageCfg    *config.StorageConfig
	imageProxyCfg *config.ImageProxyConfig
}

func NewFileService(log *log.Logger, storage storage.Storage, fileRepo repo.FileRepo, serverCfg *config.ServerConfig, storageCfg *config.StorageConfig, imageProxyCfg *config.ImageProxyConfig) *FileService {
	return &FileService{
		BaseService:   &BaseService{Log: log},
		Storage:       storage,
		FileRepo:      fileRepo,
		serverCfg:     serverCfg,
		storageCfg:    storageCfg,
		imageProxyCfg: imageProxyCfg,
	}
}

func (s *FileService) SaveFile(ctx context.Context, filename, path, mimeType, hash string, size int64, r io.Reader) (*model.File, error) {
	// 先根据 hash 值查询是否已经存在相同的文件
	existingFile, err := s.FileRepo.FindByHash(ctx, hash)
	if err == nil && existingFile != nil {
		// 找到相同 hash 的文件，直接返回，实现秒传
		s.Log.Info("文件已存在，返回现有文件", "hash", hash, "fileId", existingFile.ID)
		return existingFile, nil
	}

	// 不存在相同 hash 的文件，继续执行保存逻辑
	ext := getFileExtByMimeType(mimeType)
	uniqueFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	fullPath := uniqueFileName
	if path != "" {
		fullPath = fmt.Sprintf("%s/%s", path, uniqueFileName)
	}
	err = s.Storage.Save(ctx, fullPath, r)
	if err != nil {
		s.Log.Error("上传文件失败", "filename", filename, "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}
	file := &model.File{
		OriginalName: filename,
		Path:         fullPath,
		Storage:      s.Storage.Name(),
		Size:         size,
		MimeType:     mimeType,
		Hash:         hash,
	}
	err = s.FileRepo.BaseRepo.Create(ctx, file)
	if err != nil {
		s.Log.Error("保存文件到数据库失败", "filename", filename, "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}
	return file, nil
}

// getFileExtByMimeType 根据MIME类型获取文件扩展名
func getFileExtByMimeType(mimeType string) string {
	mimeParts := strings.Split(mimeType, "/")
	if len(mimeParts) >= 2 {
		subType := mimeParts[1]
		switch subType {
		case "jpeg", "jpg":
			return ".jpg"
		case "png":
			return ".png"
		case "gif":
			return ".gif"
		case "pdf":
			return ".pdf"
		case "mp4":
			return ".mp4"
		case "webm":
			return ".webm"
		case "zip":
			return ".zip"
		case "json":
			return ".json"
		default:
			return "." + subType
		}
	}
	return ""
}

func (s *FileService) ReadFile(ctx context.Context, id uint64) (io.ReadCloser, *model.File, error) {
	file, err := s.GetFile(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	fileReader, err := s.Storage.Open(ctx, file.Path)
	if err != nil {
		s.Log.Error("存储系统打开文件失败", "storage", s.Storage.Name(), "id", id, "error", err)
		return nil, nil, fmt.Errorf("系统内部错误")
	}
	return fileReader, file, nil
}

func (s *FileService) GetFile(ctx context.Context, id uint64) (*model.File, error) {
	file, err := s.FileRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("文件不存在", "id", id)
			return nil, fmt.Errorf("文件不存在")
		}
		s.Log.Error("查询文件失败", "id", id, "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}
	if s.Storage.Name() != file.Storage {
		s.Log.Warn("文件存储系统不匹配", "id", id, "expected", s.Storage.Name(), "actual", file.Storage)
		return nil, fmt.Errorf("文件存储系统不匹配")
	}
	return file, err
}

func (s *FileService) DeleteFile(ctx context.Context, id uint64) error {
	file, err := s.GetFile(ctx, id)
	if err != nil {
		return err
	}
	err = s.Storage.Delete(ctx, file.Path)
	if err != nil {
		s.Log.Error("存储系统删除文件失败", "storage", s.Storage.Name(), "id", id, "error", err)
		return fmt.Errorf("系统内部错误")
	}
	return s.FileRepo.BaseRepo.DeleteByID(ctx, id)
}

// HasImageProxyInternal 判断 ImageProxy 内网地址是否可用
func (s *FileService) HasImageProxyInternal() bool {
	return s.imageProxyCfg != nil && s.imageProxyCfg.InternalURL != ""
}

// GetImageProxyURL 获取 ImageProxy 代理 URL（用于 Gin 转发缩略图请求）
func (s *FileService) GetImageProxyURL(c *gin.Context, fileID uint64, width int) string {
	if !s.HasImageProxyInternal() {
		return ""
	}
	ginInternalURL := s.buildGinURL(c, fileID, true)
	return s.buildImageProxyURL(s.imageProxyCfg.InternalURL, ginInternalURL, width)
}

// getDynamicBaseURL 从请求中获取动态域名（协议+主机）
func (s *FileService) getDynamicBaseURL(c *gin.Context) string {
	proto := c.GetHeader("X-Forwarded-Proto")
	if proto == "" {
		if c.Request.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}
	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}
	return fmt.Sprintf("%s://%s", proto, host)
}

// buildGinURL 构建 Gin URL
// useInternal: true 使用配置的内网地址（给 ImageProxy 用），false 使用动态域名（给前端用）
func (s *FileService) buildGinURL(c *gin.Context, fileID uint64, useInternal bool) string {
	var baseURL string
	if useInternal && s.serverCfg.InternalURL != "" {
		baseURL = s.serverCfg.InternalURL
	} else {
		baseURL = s.getDynamicBaseURL(c)
	}
	return fmt.Sprintf("%s/api/v1/file/%d", baseURL, fileID)
}

// buildImageProxyURL 构建 ImageProxy URL（imgproxy 格式）
func (s *FileService) buildImageProxyURL(baseURL, sourceURL string, width int) string {
	if baseURL == "" {
		return sourceURL
	}
	if width == 0 {
		return fmt.Sprintf("%s/%s", baseURL, sourceURL)
	}
	return fmt.Sprintf("%s/%dx0/%s", baseURL, width, sourceURL)
}

// GetImageURL 获取图片 URL（核心逻辑）
func (s *FileService) GetImageURL(c *gin.Context, file *model.File, width int) string {
	// 1. 存储有公开链接 -> 直接返回存储公开链接
	if publicURL := s.getStoragePublicURL(file); publicURL != "" {
		return publicURL
	}

	// 2. ImageProxy 公开 -> 原图和缩略图都走 ImageProxy
	if s.imageProxyCfg != nil && s.imageProxyCfg.PublicURL != "" {
		ginInternalURL := s.buildGinURL(c, file.ID, true)
		return s.buildImageProxyURL(s.imageProxyCfg.PublicURL, ginInternalURL, width)
	}

	// 3. 都不公开 -> Gin 代理
	ginURL := s.buildGinURL(c, file.ID, false)
	if width == 0 {
		return ginURL
	}
	return fmt.Sprintf("%s?w=%d&h=%d", ginURL, width, width)
}

// getStoragePublicURL 获取存储公开链接
func (s *FileService) getStoragePublicURL(file *model.File) string {
	switch s.storageCfg.Backend {
	case "s3":
		if s.storageCfg.S3 != nil && s.storageCfg.S3.PublicURL != "" {
			return fmt.Sprintf("%s/%s", s.storageCfg.S3.PublicURL, file.Path)
		}
	case "webdav":
		if s.storageCfg.WebDAV != nil && s.storageCfg.WebDAV.PublicURL != "" {
			return fmt.Sprintf("%s/%s", s.storageCfg.WebDAV.PublicURL, file.Path)
		}
	}
	return ""
}

// FileResponse 文件URL响应模型
type FileResponse struct {
	ID        uint64 `json:"id"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	Name      string `json:"name,omitempty"`
	Size      int64  `json:"size,omitempty"`
	MimeType  string `json:"mime_type,omitempty"`
}

// BuildFileResponse 构建文件响应对象
func (s *FileService) BuildFileResponse(c *gin.Context, file *model.File) *FileResponse {
	if file == nil {
		return nil
	}
	return &FileResponse{
		ID:        file.ID,
		URL:       s.GetImageURL(c, file, 0),
		Thumbnail: s.GetImageURL(c, file, 200),
		Name:      file.OriginalName,
		Size:      file.Size,
		MimeType:  file.MimeType,
	}
}
