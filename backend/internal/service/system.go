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
	"github.com/bookandmusic/love-girl/internal/utils"
)

type SystemService struct {
	*BaseService
	UserRepo    repo.UserRepo
	SettingRepo repo.SettingRepo
	AlbumRepo   repo.AlbumRepo
	PlaceRepo   repo.PlaceRepo
	MomentRepo  repo.MomentRepo
	WishRepo    repo.WishRepo
	FileService *FileService
	config      *config.AppConfig
	jwt         auth.JWT
}

func NewSystemService(
	log *log.Logger,
	userRepo repo.UserRepo,
	settingRepo repo.SettingRepo,
	albumRepo repo.AlbumRepo,
	placeRepo repo.PlaceRepo,
	momentRepo repo.MomentRepo,
	wishRepo repo.WishRepo,
	fileService *FileService,
	config *config.AppConfig,
	jwt auth.JWT,
) *SystemService {
	return &SystemService{
		BaseService: &BaseService{Log: log},
		UserRepo:    userRepo,
		SettingRepo: settingRepo,
		AlbumRepo:   albumRepo,
		PlaceRepo:   placeRepo,
		MomentRepo:  momentRepo,
		WishRepo:    wishRepo,
		FileService: fileService,
		config:      config,
		jwt:         jwt,
	}
}

type InitSystemRequest struct {
	SiteName            string `json:"siteName" binding:"required"`
	SiteDescription     string `json:"siteDescription"`
	StartDate           string `json:"startDate" binding:"required"`
	UserAName           string `json:"userAName" binding:"required"`
	UserARole           string `json:"userARole" binding:"required,oneof=boy girl"`
	UserAEmail          string `json:"userAEmail"`
	UserAPhone          string `json:"userAPhone"`
	UserBName           string `json:"userBName" binding:"required"`
	UserBRole           string `json:"userBRole" binding:"required,oneof=boy girl"`
	UserBEmail          string `json:"userBEmail"`
	UserBPhone          string `json:"userBPhone"`
	SitePassword        string `json:"sitePassword" binding:"required,min=6"`
	SitePasswordConfirm string `json:"sitePasswordConfirm" binding:"required,eqfield=SitePassword"`
}

// FrontendSystemInfo 适配前端期望的数据结构
type FrontendSystemInfo struct {
	Site struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		StartDate   string `json:"startDate"`
	} `json:"site"`
	Couple struct {
		Boy struct {
			Name   string        `json:"name"`
			Avatar *FileResponse `json:"avatar"`
		} `json:"boy"`
		Girl struct {
			Name   string        `json:"name"`
			Avatar *FileResponse `json:"avatar"`
		} `json:"girl"`
	} `json:"couple"`
}

// InitSystem 初始化系统
func (s *SystemService) InitSystem(ctx context.Context, req *InitSystemRequest) error {
	// 检查系统是否已经初始化
	isInitialized, err := s.IsSystemInitialized(ctx)
	if err != nil {
		s.Log.Error("检查初始化状态失败", "error", err)
		return fmt.Errorf("系统内部错误")
	}

	if isInitialized {
		s.Log.Info("尝试重复初始化系统")
		return fmt.Errorf("系统已经初始化，无法重复初始化")
	}

	// 创建初始数据
	if err := s.createInitialData(ctx, req); err != nil {
		s.Log.Error("创建初始数据失败", "error", err)
		return fmt.Errorf("系统内部错误")
	}

	return nil
}

// createInitialData 创建初始数据
func (s *SystemService) createInitialData(ctx context.Context, req *InitSystemRequest) error {
	encryptedPassword, err := utils.EncryptPassword(req.SitePassword)
	if err != nil {
		s.Log.Error("加密密码失败", "error", err)
		return fmt.Errorf("系统内部错误")
	}

	// 创建用户A，只有在邮箱非空时才设置邮箱
	var userEmailA *string
	if req.UserAEmail != "" {
		userEmailA = &req.UserAEmail
	}
	userA := &model.User{
		Name:     req.UserAName,
		Email:    userEmailA,
		Password: encryptedPassword,
		Role:     req.UserARole,
		Phone:    req.UserAPhone,
	}

	if err := s.UserRepo.Create(ctx, userA); err != nil {
		s.Log.Error("创建用户A失败", "error", err, "username", req.UserAName)
		return fmt.Errorf("系统内部错误")
	}

	// 创建用户B，只有在邮箱非空时才设置邮箱
	var userEmailB *string
	if req.UserBEmail != "" {
		userEmailB = &req.UserBEmail
	}
	userB := &model.User{
		Name:     req.UserBName,
		Email:    userEmailB,
		Password: encryptedPassword,
		Role:     req.UserBRole,
		Phone:    req.UserBPhone,
	}

	if err := s.UserRepo.Create(ctx, userB); err != nil {
		s.Log.Error("创建用户B失败", "error", err, "username", req.UserBName)
		return fmt.Errorf("系统内部错误")
	}

	// 创建设置
	settings := []model.Setting{
		{Key: "siteTitle", Value: req.SiteName, Type: "text", Label: "站点标题", Group: "general"},
		{Key: "siteDescription", Value: req.SiteDescription, Type: "textarea", Label: "站点描述", Group: "general"},
		{Key: "startDate", Value: req.StartDate, Type: "date", Label: "故事开始日期", Group: "general"},
	}

	for _, setting := range settings {
		if err := s.SettingRepo.Create(ctx, &setting); err != nil {
			s.Log.Error("创建设置项失败", "error", err, "key", setting.Key)
			return fmt.Errorf("系统内部错误")
		}
	}

	return nil
}

// IsSystemInitialized 检查系统是否已初始化
func (s *SystemService) IsSystemInitialized(ctx context.Context) (bool, error) {
	// 检查数据库中是否存在站点配置（siteTitle）
	setting, err := s.SettingRepo.GetSettingByKey(ctx, "siteTitle")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		s.Log.Error("获取站点配置失败", "error", err)
		return false, fmt.Errorf("系统内部错误")
	}

	return setting != nil, nil
}

// GetSystemInfo 获取系统信息
func (s *SystemService) GetSystemInfo(ctx context.Context) (*FrontendSystemInfo, error) {
	settings, err := s.getSystemSettings(ctx)
	if err != nil {
		return nil, err
	}

	users, err := s.getAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return s.buildFrontendInfo(ctx, settings, users), nil
}

func (s *SystemService) getSystemSettings(ctx context.Context) (map[string]*model.Setting, error) {
	keys := []string{"siteTitle", "siteDescription", "startDate"}
	settings := make(map[string]*model.Setting)

	for _, key := range keys {
		setting, err := s.SettingRepo.GetSettingByKey(ctx, key)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Info("系统未初始化，缺少设置", "key", key)
				return nil, errors.New("系统未初始化")
			}
			s.Log.Error("获取设置失败", "key", key, "error", err)
			return nil, fmt.Errorf("系统内部错误")
		}
		settings[key] = setting
	}

	return settings, nil
}

func (s *SystemService) getAllUsers(ctx context.Context) ([]model.User, error) {
	users, _, err := s.UserRepo.ListUsers(ctx, 1, 100)
	if err != nil {
		s.Log.Error("获取用户列表失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}
	return users, nil
}

func (s *SystemService) buildFrontendInfo(ctx context.Context, settings map[string]*model.Setting, users []model.User) *FrontendSystemInfo {
	frontendInfo := &FrontendSystemInfo{}

	frontendInfo.Site.Name = settings["siteTitle"].Value
	frontendInfo.Site.Description = settings["siteDescription"].Value
	frontendInfo.Site.StartDate = settings["startDate"].Value

	for i := range users {
		user := &users[i]
		switch user.Role {
		case "boy":
			frontendInfo.Couple.Boy.Name = user.Name
			frontendInfo.Couple.Boy.Avatar = s.FileService.BuildFileResponse(ctx, user.Avatar)
		case "girl":
			frontendInfo.Couple.Girl.Name = user.Name
			frontendInfo.Couple.Girl.Avatar = s.FileService.BuildFileResponse(ctx, user.Avatar)
		}
	}

	return frontendInfo
}

// GetSettings 获取所有站点设置
func (s *SystemService) GetSettings(ctx context.Context) (map[string]string, error) {
	settings, err := s.SettingRepo.GetSettingsByGroup(ctx, "general")
	if err != nil {
		s.Log.Error("获取通用设置失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}

	return result, nil
}

// SaveSettings 保存站点设置
func (s *SystemService) SaveSettings(ctx context.Context, settings map[string]string) error {
	if settings["siteTitle"] == "" {
		return errors.New("站点标题不能为空")
	}
	for key, value := range settings {
		setting := &model.Setting{
			Value: value,
			Type:  "text",
			Label: getLabelByKey(key),
			Group: "general",
		}
		err := s.SettingRepo.UpdateOrCreateByKey(ctx, key, setting)
		if err != nil {
			s.Log.Error("更新或创建设置项失败", "error", err, "key", key)
			return fmt.Errorf("系统内部错误")
		}
	}
	return nil
}

// GetDashboardStats 获取仪表盘统计数据
func (s *SystemService) GetDashboardStats(ctx context.Context) (*model.DashboardStats, error) {
	// 获取相册统计
	albumCount, err := s.AlbumRepo.GetAlbumCount(ctx)
	if err != nil {
		s.Log.Error("获取相册数量失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 获取照片总数
	photoCount, err := s.AlbumRepo.GetPhotoCount(ctx)
	if err != nil {
		s.Log.Error("获取照片数量失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 获取地点统计
	placeCountVal, err := s.PlaceRepo.BaseRepo.Count(ctx)
	if err != nil {
		s.Log.Error("获取地点数量失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}
	placeCount := int(placeCountVal)

	// 获取动态统计
	momentCount, err := s.MomentRepo.BaseRepo.Count(ctx)
	if err != nil {
		s.Log.Error("获取动态数量失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}

	// 获取愿望统计
	wishCountVal, err := s.WishRepo.BaseRepo.Count(ctx)
	if err != nil {
		s.Log.Error("获取愿望数量失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}
	wishCount := int(wishCountVal)

	// 获取待审核愿望数量
	pendingWishCountVal, err := s.WishRepo.BaseRepo.CountWithConditions(ctx,
		repo.FilterCondition{Field: "approved", Operator: "eq", Value: false},
	)
	if err != nil {
		s.Log.Error("获取待审核愿望数量失败", "error", err)
		return nil, fmt.Errorf("系统内部错误")
	}
	pendingWishCount := int(pendingWishCountVal)

	return &model.DashboardStats{
		AlbumStats: model.AlbumStats{
			Total:       albumCount,
			TotalPhotos: photoCount,
		},
		PlaceStats: model.PlaceStats{
			Total: placeCount,
		},
		MomentStats: model.MomentStats{
			Total: int(momentCount),
		},
		WishStats: model.WishStats{
			Total:   wishCount,
			Pending: pendingWishCount,
		},
	}, nil
}

// getLabelByKey 根据键名获取标签
func getLabelByKey(key string) string {
	labels := map[string]string{
		"siteTitle":       "站点标题",
		"siteName":        "站点名称",
		"siteDescription": "站点描述",
		"startDate":       "故事开始日期",
	}
	if label, ok := labels[key]; ok {
		return label
	}
	return key
}
