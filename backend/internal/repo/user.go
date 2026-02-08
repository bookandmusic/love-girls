package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/auth"
	"github.com/bookandmusic/love-girl/internal/model"
)

// UserRepo 用户仓库
// 功能：
//   - 创建用户（可同时创建头像关联）
//   - 删除用户
//   - 查询用户（支持单查、列表、分页）
//   - 更新用户（支持头像历史管理）
//
// 说明：用户的头像是允许不存在的
type UserRepo struct {
	*BaseRepo[model.User]
	JWT            auth.JWT
	entityFileRepo *EntityFileRepo
}

// NewUserRepo 创建新的用户仓库实例
func NewUserRepo(dbCli *gorm.DB, jwt auth.JWT) *UserRepo {
	return &UserRepo{
		BaseRepo:       NewBaseRepo[model.User](dbCli),
		JWT:            jwt,
		entityFileRepo: NewEntityFileRepo(dbCli),
	}
}

// FindOneByKey 根据邮箱或用户名查找用户
// 参数：
//   - ctx: 上下文
//   - key: 邮箱或用户名
//
// 返回：用户实体、错误
func (s *UserRepo) FindOneByKey(ctx context.Context, key string) (*model.User, error) {
	// 由于需要 OR 条件，使用自定义查询
	var user model.User
	if err := s.BaseRepo.DB().WithContext(ctx).Where("email = ? OR name = ?", key, key).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindOneByEmail 根据邮箱查找用户
// 参数：
//   - ctx: 上下文
//   - email: 邮箱
//
// 返回：用户实体、错误
func (s *UserRepo) FindOneByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.FindOne(ctx, WithConditions(
		FilterCondition{Field: "email", Operator: "eq", Value: email},
	))
}

// FindOneByName 根据用户名查找用户
// 参数：
//   - ctx: 上下文
//   - name: 用户名
//
// 返回：用户实体、错误
func (s *UserRepo) FindOneByName(ctx context.Context, name string) (*model.User, error) {
	return s.FindOne(ctx, WithConditions(
		FilterCondition{Field: "name", Operator: "eq", Value: name},
	))
}

// CreateWithAvatar 创建用户记录，如果存在头像图片，需要同步添加一条关联关系
// 参数：
//   - ctx: 上下文
//   - user: 用户实体（包含AvatarID字段）
//
// 返回：如果创建成功返回nil，否则返回错误
func (s *UserRepo) CreateWithAvatar(ctx context.Context, user *model.User) error {
	// 创建用户记录
	if err := s.BaseRepo.Create(ctx, user); err != nil {
		return err
	}

	// 如果存在头像ID，创建关联关系，记录头像历史
	if user.AvatarID != nil {
		return s.entityFileRepo.AssociateFile(ctx, user.ID, "user_avatar", *user.AvatarID)
	}

	return nil
}

// ListUsers 分页查询所有用户信息，关联查询头像
// 参数：
//   - ctx: 上下文
//   - page: 页码，从1开始
//   - size: 每页数量
//   - opts: 查询选项（排序、过滤等）
//
// 返回：用户列表、总数、错误
func (s *UserRepo) ListUsers(ctx context.Context, page, size int, opts ...QueryOption) ([]model.User, int64, error) {
	// 确保预加载关联数据
	allOpts := append(opts, WithPreloads("Avatar"))
	return s.BaseRepo.FindWithPagination(ctx, page, size, allOpts...)
}

// GetAvatarHistory 查询用户头像历史，即关联的所有图片
// 参数：
//   - ctx: 上下文
//   - userID: 用户ID
//
// 返回：头像历史记录、错误
func (s *UserRepo) GetAvatarHistory(ctx context.Context, userID uint64) ([]model.EntityFile, error) {
	// 查询用户的头像历史记录
	return s.entityFileRepo.ListByEntity(ctx, userID, "user_avatar")
}

// GetAvatarHistoryWithPagination 分页查询用户头像历史，即关联的所有图片
// 参数：
//   - ctx: 上下文
//   - userID: 用户ID
//   - page: 页码，从1开始
//   - size: 每页数量
//   - opts: 查询选项（排序等）
//
// 返回：文件列表、总数、错误
func (s *UserRepo) GetAvatarHistoryWithPagination(ctx context.Context, userID uint64, page, size int, opts ...QueryOption) ([]model.File, int64, error) {
	// 分页查询用户的头像历史记录
	return s.entityFileRepo.GetAssociatedFilesWithPagination(ctx, userID, "user_avatar", page, size, opts...)
}

// updateAvatarAssociation 更新头像关联关系
// 参数：
//   - ctx: 上下文
//   - userID: 用户ID
//   - oldAvatarID: 旧头像ID
//   - newAvatarID: 新头像ID
//
// 返回：如果更新成功返回nil，否则返回错误
func (s *UserRepo) updateAvatarAssociation(ctx context.Context, userID uint64, oldAvatarID, newAvatarID *uint64) error {
	// 如果新旧头像ID不同，需要更新关联关系
	if (oldAvatarID == nil && newAvatarID != nil) ||
		(oldAvatarID != nil && newAvatarID == nil) ||
		(oldAvatarID != nil && newAvatarID != nil && *oldAvatarID != *newAvatarID) {
		// 如果旧头像存在，删除旧头像的关联关系
		if oldAvatarID != nil {
			_ = s.entityFileRepo.RemoveAssociation(ctx, userID, "user_avatar", *oldAvatarID)
		}

		// 如果新头像存在，添加新的关联关系
		if newAvatarID != nil {
			_ = s.entityFileRepo.AssociateFile(ctx, userID, "user_avatar", *newAvatarID)
		}
	}

	return nil
}

// UpdateUserInfo 修改用户信息，如果存在头像，会自动处理头像历史关联关系
// 参数：
//   - ctx: 上下文
//   - user: 用户实体（包含AvatarID字段）
//
// 返回：如果更新成功返回nil，否则返回错误
func (s *UserRepo) UpdateUserInfo(ctx context.Context, user *model.User) error {
	// 先查询用户信息，获取旧的头像ID
	oldUser, err := s.BaseRepo.FindByID(ctx, user.ID)
	if err != nil {
		return err
	}

	// 更新用户信息
	if err := s.BaseRepo.Update(ctx, user); err != nil {
		return err
	}

	// 更新头像关联关系
	return s.updateAvatarAssociation(ctx, user.ID, oldUser.AvatarID, user.AvatarID)
}

// UpdateAvatar 更新用户头像
// 参数：
//   - ctx: 上下文
//   - userID: 用户ID
//   - avatarID: 新头像ID（可选）
//
// 返回：如果更新成功返回nil，否则返回错误
func (s *UserRepo) UpdateAvatar(ctx context.Context, userID uint64, avatarID *uint64) error {
	// 先查询用户信息，获取旧的头像ID
	user, err := s.BaseRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 保存旧的头像ID
	oldAvatarID := user.AvatarID
	user.AvatarID = avatarID

	// 更新用户信息
	if err := s.Update(ctx, user); err != nil {
		return err
	}

	// 更新头像关联关系
	return s.updateAvatarAssociation(ctx, userID, oldAvatarID, avatarID)
}

// DeleteAvatarHistory 删除某个历史头像，即删除关联关系
// 参数：
//   - ctx: 上下文
//   - userID: 用户ID
//   - fileID: 文件ID
//
// 返回：如果删除成功返回nil，否则返回错误
func (s *UserRepo) DeleteAvatarHistory(ctx context.Context, userID uint64, fileID uint64) error {
	// 删除EntityFile记录
	return s.entityFileRepo.RemoveAssociation(ctx, userID, "user_avatar", fileID)
}
