package repo

import (
	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/model"
)

// AnniversaryRepo 纪念日仓库
type AnniversaryRepo struct {
	*BaseRepo[model.Anniversary]
}

// NewAnniversaryRepo 创建新的纪念日仓库实例
func NewAnniversaryRepo(dbCli *gorm.DB) *AnniversaryRepo {
	return &AnniversaryRepo{
		BaseRepo: NewBaseRepo[model.Anniversary](dbCli),
	}
}
