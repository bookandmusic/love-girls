package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/config"
	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/model"
	"github.com/bookandmusic/love-girl/internal/repo"
	"github.com/bookandmusic/love-girl/internal/storage"
)

type FileService struct {
	*BaseService
	Storage   storage.Storage
	FileRepo  repo.FileRepo
	serverCfg *config.ServerConfig
}

func NewFileService(log *log.Logger, storage storage.Storage, fileRepo repo.FileRepo, serverCfg *config.ServerConfig) *FileService {
	return &FileService{
		BaseService: &BaseService{Log: log},
		Storage:     storage,
		FileRepo:    fileRepo,
		serverCfg:   serverCfg,
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
	uniqueFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext) // 使用时间戳纳秒值作为唯一文件名
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

func (s *FileService) JoinFileURL(id uint64) string {
	return fmt.Sprintf("%s://%s/api/v1/file/%d", s.serverCfg.Schema, s.serverCfg.HostName, id)
}

func (s *FileService) OriginalImageURL(ctx context.Context, file *model.File) (string, error) {
	return s.Storage.URL(ctx, file.ID, file.Path, 0, 0, s.JoinFileURL)
}

func (s *FileService) ThumbnailImageURL(ctx context.Context, file *model.File, width, height int) (string, error) {
	return s.Storage.URL(ctx, file.ID, file.Path, width, height, s.JoinFileURL)
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

// FileURLs 文件URL结构体
type FileURLs struct {
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

// GetFileURLs 获取文件URL和缩略图URL
func (s *FileService) GetFileURLs(ctx context.Context, file *model.File) FileURLs {
	urls := FileURLs{}

	if file != nil {
		// 获取原图URL，失败时使用空字符串
		if originalURL, err := s.OriginalImageURL(ctx, file); err == nil && originalURL != "" {
			urls.URL = originalURL
		}

		// 获取缩略图URL，失败时使用空字符串
		if thumbnailURL, err := s.ThumbnailImageURL(ctx, file, 200, 200); err == nil && thumbnailURL != "" {
			urls.Thumbnail = thumbnailURL
		}
	}

	return urls
}

// BuildFileResponse 构建文件响应对象
func (s *FileService) BuildFileResponse(ctx context.Context, file *model.File) *FileResponse {
	if file == nil {
		return nil
	}

	urls := s.GetFileURLs(ctx, file)

	return &FileResponse{
		ID:        file.ID,
		URL:       urls.URL,
		Thumbnail: urls.Thumbnail,
		Name:      file.OriginalName,
		Size:      file.Size,
		MimeType:  file.MimeType,
	}
}
