package model

import "time"

type Wish struct {
	BaseModel
	Content    string    `gorm:"type:text;not null" json:"content"`
	AuthorName string    `gorm:"size:100;not null" json:"author_name"`
	Email      string    `gorm:"size:150" json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	Approved   bool      `gorm:"default:false" json:"approved"`
}

func (Wish) TableName() string {
	return "wishes"
}
