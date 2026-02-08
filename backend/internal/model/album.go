package model

type Album struct {
	BaseModel
	Name         string       `gorm:"size:255;not null" json:"name"`
	Description  string       `gorm:"size:512" json:"description"`
	CoverImageID *uint64      `gorm:"index" json:"cover_image_id,omitempty"` // 封面图片ID，外键引用files表
	CoverImage   *File        `gorm:"foreignKey:CoverImageID" json:"cover_image,omitempty"`
	PhotoCount   int          `gorm:"default:0" json:"photo_count"`
	EntityFiles  []EntityFile `gorm:"foreignKey:EntityID;constraint:-" json:"-"` // 相册照片（多态关联，禁止外键约束）
}

func (Album) TableName() string {
	return "albums"
}

// GetEntityFiles 实现 EntityFilesGetter 接口
func (a *Album) GetEntityFiles() []EntityFile {
	return a.EntityFiles
}
