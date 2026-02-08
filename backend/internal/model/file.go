package model

// File represents a file record in the database.
type File struct {
	BaseModel
	OriginalName string `gorm:"type:varchar(255);not null" json:"original_name"`
	Storage      string `gorm:"type:varchar(32);not null" json:"storage"` // local | s3 | webdav
	Path         string `gorm:"type:varchar(512);not null" json:"path"`   // path in the storage system
	Size         int64  `gorm:"not null" json:"size"`
	MimeType     string `gorm:"type:varchar(128)" json:"mime_type,omitempty"`
	Hash         string `gorm:"type:char(64)" json:"hash,omitempty"`
}
