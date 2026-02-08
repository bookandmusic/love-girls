package model

// Moment 动态表
type Moment struct {
	BaseModel
	Content     string       `gorm:"column:content;type:text;not null" json:"content"`
	Likes       int          `gorm:"column:likes;type:int;default:0;not null" json:"likes"`
	IsPublic    bool         `gorm:"column:is_public;type:boolean;not null" json:"is_public"`
	UserID      uint64       `gorm:"column:user_id;type:bigint;not null;index" json:"user_id"` // 关联用户
	User        *User        `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	EntityFiles []EntityFile `gorm:"foreignKey:EntityID;references:ID;constraint:-" json:"entity_files"` // 关联的文件记录（多态关联，禁止外键约束）
}

func (Moment) TableName() string {
	return "moments"
}

// GetEntityFiles 实现 EntityFilesGetter 接口
func (m *Moment) GetEntityFiles() []EntityFile {
	return m.EntityFiles
}
