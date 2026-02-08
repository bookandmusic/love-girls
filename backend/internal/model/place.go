package model

type Place struct {
	BaseModel
	Name        string  `gorm:"size:255;not null" json:"name"`
	Latitude    float64 `gorm:"type:decimal(10,8);not null" json:"latitude"`  // 纬度
	Longitude   float64 `gorm:"type:decimal(11,8);not null" json:"longitude"` // 经度
	ImageID     *uint64 `gorm:"index" json:"image_id"`                        // 主图ID，外键关联files表
	Image       *File   `gorm:"foreignKey:ImageID" json:"image,omitempty"`
	Description string  `gorm:"type:text" json:"description"`
	Date        string  `gorm:"size:20" json:"date"` // 格式如 'YYYY-MM-DD'
}

func (Place) TableName() string {
	return "places"
}
