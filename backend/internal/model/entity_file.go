package model

// EntityFile 是用于表示任意实体与文件之间的多对多关系的中间表
type EntityFile struct {
	BaseModel
	EntityID   uint64 `gorm:"index:idx_entity_files_entity;not null" json:"entity_id"`             // 关联实体的ID
	EntityType string `gorm:"size:50;not null;index:idx_entity_files_entity" json:"entity_type"`   // 关联实体的类型，如 'album', 'moment', 'user_avatar'
	FileID     uint64 `gorm:"index;not null" json:"file_id"`                                       // 文件ID
	File       *File  `gorm:"foreignKey:FileID;constraint:OnDelete:CASCADE" json:"file,omitempty"` // 关联的文件
}

func (EntityFile) TableName() string {
	return "entity_files" // 通用实体文件关联表
}
