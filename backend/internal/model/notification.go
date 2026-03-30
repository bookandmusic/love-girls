package model

type NotificationType string

const (
	NotificationTypeComment NotificationType = "comment"
	NotificationTypeReply   NotificationType = "reply"
)

type Notification struct {
	BaseModel
	UserID    uint64           `gorm:"not null;index" json:"user_id"`
	SenderID  uint64           `gorm:"not null;index" json:"sender_id"`
	MomentID  uint64           `gorm:"not null;index" json:"moment_id"`
	CommentID uint64           `gorm:"not null;index" json:"comment_id"`
	Type      NotificationType `gorm:"type:varchar(20);not null" json:"type"`
	Content   string           `gorm:"type:text" json:"content"`
	IsRead    bool             `gorm:"default:false" json:"is_read"`
	User      *User            `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Sender    *User            `gorm:"foreignKey:SenderID;references:ID;constraint:OnDelete:CASCADE" json:"sender,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}
