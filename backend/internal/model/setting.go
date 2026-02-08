package model

type Setting struct {
	BaseModel
	Key    string `gorm:"size:100;not null;uniqueIndex" json:"key"` // 配置键名
	Value  string `gorm:"type:text" json:"value"`                   // 配置值
	Type   string `gorm:"size:50;default:'text'" json:"type"`       // 配置类型 (text, password, json, boolean)
	Label  string `gorm:"size:200" json:"label"`                    // 配置标签
	Group  string `gorm:"size:50;index" json:"group"`               // 配置组 (general, ai, storage, notification)
	Remark string `gorm:"type:text" json:"remark"`                  // 配置备注
}

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	AlbumStats  AlbumStats  `json:"albumStats"`
	PlaceStats  PlaceStats  `json:"placeStats"`
	MomentStats MomentStats `json:"momentStats"`
	WishStats   WishStats   `json:"wishStats"`
}

// AlbumStats 相册统计数据
type AlbumStats struct {
	Total       int `json:"total"`
	TotalPhotos int `json:"totalPhotos"`
}

// PlaceStats 地点统计数据
type PlaceStats struct {
	Total int `json:"total"`
}

// MomentStats 动态统计数据
type MomentStats struct {
	Total int `json:"total"`
}

// WishStats 愿望统计数据
type WishStats struct {
	Total   int `json:"total"`
	Pending int `json:"pending"`
}
