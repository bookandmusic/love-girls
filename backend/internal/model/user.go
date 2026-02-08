package model

type User struct {
	BaseModel
	Name        string       `gorm:"size:64;not null" json:"name"`
	Email       *string      `gorm:"size:128;uniqueIndex" json:"email"` // 可为空的邮箱，只对非空值强制唯一
	Password    string       `gorm:"size:128;not null" json:"password"`
	Role        string       `gorm:"size:32" json:"role"`    // 用户角色
	Phone       string       `gorm:"size:20" json:"phone"`   // 用户手机号
	AvatarID    *uint64      `gorm:"index" json:"avatar_id"` // 当前头像ID，外键关联files表
	Avatar      *File        `gorm:"foreignKey:AvatarID" json:"avatar,omitempty"`
	EntityFiles []EntityFile `gorm:"foreignKey:EntityID;constraint:-" json:"-"` // 头像历史（多态关联，禁止外键约束）
}

// GetEntityFiles 实现 EntityFilesGetter 接口
func (u *User) GetEntityFiles() []EntityFile {
	return u.EntityFiles
}
