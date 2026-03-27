package model

type Comment struct {
	BaseModel
	Content   string  `gorm:"type:text;not null" json:"content"`
	MomentID  uint64  `gorm:"not null;index" json:"moment_id"`
	ParentID  *uint64 `gorm:"index" json:"parent_id"`
	ReplyToID *uint64 `gorm:"index" json:"reply_to_id"`
	UserID    uint64  `gorm:"not null;index" json:"user_id"`
	User      *User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Path      string  `gorm:"type:varchar(512);index" json:"path"`
	Depth     int     `gorm:"default:0" json:"depth"`
}

func (Comment) TableName() string {
	return "comments"
}
