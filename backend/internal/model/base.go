package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`
}

// EntityFilesGetter 定义获取关联文件的接口
type EntityFilesGetter interface {
	GetEntityFiles() []EntityFile
}

// Images 从实现了 EntityFilesGetter 的模型中获取所有图片
func Images[T EntityFilesGetter](model T) []File {
	var images []File
	for _, ef := range model.GetEntityFiles() {
		if ef.File != nil {
			images = append(images, *ef.File)
		}
	}
	return images
}
