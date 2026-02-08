package model

// Anniversary 纪念日表，存储重要的纪念日信息
type Anniversary struct {
	BaseModel
	Title       string `gorm:"size:255;not null" json:"title"`
	Date        string `gorm:"size:20;not null" json:"date"` // 格式如 'YYYY-MM-DD'
	Description string `gorm:"type:text" json:"description"`
	Calendar    string `gorm:"size:10;default:'solar'" json:"calendar"` // 日历类型：solar/lunar
}

func (Anniversary) TableName() string {
	return "anniversaries"
}
