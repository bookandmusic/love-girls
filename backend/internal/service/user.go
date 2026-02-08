package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/auth"
	"github.com/bookandmusic/love-girl/internal/config"
	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/model"
	"github.com/bookandmusic/love-girl/internal/repo"
	"github.com/bookandmusic/love-girl/internal/storage"
	"github.com/bookandmusic/love-girl/internal/utils"
)

type UserService struct {
	*BaseService
	UserRepo    repo.UserRepo
	FileRepo    repo.FileRepo
	FileService *FileService
	Storage     storage.Storage
	JWT         auth.JWT
	serverCfg   *config.ServerConfig
}

func NewUserService(log *log.Logger, userRepo repo.UserRepo, fileRepo repo.FileRepo, fileService *FileService, storage storage.Storage, serverCfg *config.ServerConfig, jwt auth.JWT) *UserService {
	return &UserService{
		BaseService: &BaseService{Log: log},
		UserRepo:    userRepo,
		FileRepo:    fileRepo,
		FileService: fileService,
		Storage:     storage,
		JWT:         jwt,
		serverCfg:   serverCfg,
	}
}

func (s *UserService) GenerateToken(ctx context.Context, username, password string) (*model.User, string, error) {
	user, err := s.UserRepo.FindOneByKey(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("用户登录失败，用户不存在", "username", username)
			return nil, "", fmt.Errorf("用户名或密码错误")
		}
		s.Log.Error("用户查询失败", "error", err, "username", username)
		return nil, "", fmt.Errorf("系统内部错误")
	}

	if !utils.VerifyPassword(user.Password, password) {
		s.Log.Info("用户登录失败，密码错误", "username", username)
		return nil, "", fmt.Errorf("用户名或密码错误")
	}
	token, err := s.JWT.Generate(&auth.Claims{
		Role:   "user",
		UserID: uint64(user.ID),
	})
	if err != nil {
		s.Log.Error("用户生成token失败", "error", err, "username", username)
		return nil, "", fmt.Errorf("系统内部错误")
	}
	return user, token, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, userID uint64) (*model.User, error) {
	user, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("查询用户信息失败，用户不存在", "userID", userID)
			return nil, fmt.Errorf("用户不存在")
		}
		s.Log.Error("用户信息查询失败", "error", err, "userID", userID)
		return nil, fmt.Errorf("系统内部错误")
	}
	return user, nil
}

// UserResponse 用户响应数据结构
type UserResponse struct {
	ID       uint64        `json:"id"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	Role     string        `json:"role"`
	JoinDate string        `json:"joinDate"`
	Avatar   *FileResponse `json:"avatar,omitempty"`
}

// GetUsers 获取所有用户列表
func (s *UserService) GetUsers(ctx context.Context) ([]UserResponse, error) {
	users, _, err := s.UserRepo.ListUsers(ctx, 1, 100) // 暂时使用固定分页参数
	if err != nil {
		s.Log.Error("获取用户列表失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 转换为响应格式
	var userResponses []UserResponse
	for _, user := range users {
		// 获取email值
		email := ""
		if user.Email != nil {
			email = *user.Email
		}

		userResponse := UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    email,
			Role:     user.Role,
			JoinDate: user.CreatedAt.Format("2006-01-02"),
			Avatar:   s.FileService.BuildFileResponse(ctx, user.Avatar),
		}
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}

// JoinFileURL 生成文件URL
func (s *UserService) JoinFileURL(id uint64) string {
	return fmt.Sprintf("%s://%s/api/v1/file/%d", s.serverCfg.Schema, s.serverCfg.HostName, id)
}

// GetAvatarURL 获取用户头像URL
func (s *UserService) GetAvatarURL(ctx context.Context, user *model.User) (string, error) {
	if user.AvatarID == nil {
		return "", nil
	}

	filePath := ""
	if user.Avatar != nil {
		filePath = user.Avatar.Path
	} else {
		// 如果Avatar关联未加载，通过FileRepo查询文件路径
		file, err := s.FileRepo.FindByID(ctx, *user.AvatarID)
		if err == nil && file != nil {
			filePath = file.Path
		}
	}

	return s.Storage.URL(ctx, *user.AvatarID, filePath, 32, 32, s.JoinFileURL)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, userID uint64, name, email string, avatarID *uint64, newPassword string) (*UserResponse, error) {
	// 先查询用户信息
	user, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.Info("更新用户信息失败，用户不存在", "userID", userID)
			return nil, fmt.Errorf("用户不存在")
		}
		s.Log.Error("查询用户信息失败", "error", err, "userID", userID)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 更新用户信息
	user.Name = name
	if email != "" {
		user.Email = &email
	} else {
		user.Email = nil
	}
	user.AvatarID = avatarID

	// 如果提供了新密码，则更新密码
	if newPassword != "" {
		hashedPassword, err := utils.EncryptPassword(newPassword)
		if err != nil {
			s.Log.Error("密码加密失败", "error", err, "userID", userID)
			return nil, fmt.Errorf("系统内部错误")
		}
		user.Password = hashedPassword
	}

	// 保存更新
	if err := s.UserRepo.UpdateUserInfo(ctx, user); err != nil {
		s.Log.Error("更新用户信息失败", "error", err, "userID", userID)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 获取email值
	userEmail := ""
	if user.Email != nil {
		userEmail = *user.Email
	}

	userResponse := &UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    userEmail,
		Role:     user.Role,
		JoinDate: user.CreatedAt.Format("2006-01-02"),
		Avatar:   s.FileService.BuildFileResponse(ctx, user.Avatar),
	}

	return userResponse, nil
}

// AvatarHistoryItem 头像历史项
type AvatarHistoryItem struct {
	ID           uint64 `json:"id"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
	CreatedAt    string `json:"createdAt"`
}

// AvatarHistoryResponse 头像历史响应
type AvatarHistoryResponse struct {
	Avatars    []AvatarHistoryItem `json:"avatars"`
	Page       int                 `json:"page"`
	Size       int                 `json:"size"`
	Total      int64               `json:"total"`
	TotalPages int                 `json:"totalPages"`
}

// GetAvatarHistory 获取用户头像历史
func (s *UserService) GetAvatarHistory(ctx context.Context, userID uint64, page, size int) (*AvatarHistoryResponse, error) {
	files, total, err := s.UserRepo.GetAvatarHistoryWithPagination(ctx, userID, page, size)
	if err != nil {
		s.Log.Error("获取头像历史失败", "error", err, "userID", userID)
		return nil, fmt.Errorf("系统内部错误")
	}

	avatars := make([]AvatarHistoryItem, 0, len(files))
	for _, file := range files {
		urls := s.FileService.GetFileURLs(ctx, &file)
		avatars = append(avatars, AvatarHistoryItem{
			ID:           file.ID,
			URL:          urls.URL,
			ThumbnailURL: urls.Thumbnail,
			CreatedAt:    file.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalPage := int((total + int64(size) - 1) / int64(size))

	return &AvatarHistoryResponse{
		Avatars:    avatars,
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPage,
	}, nil
}
