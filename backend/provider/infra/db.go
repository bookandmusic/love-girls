package infra

import (
	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/config"
	"github.com/bookandmusic/love-girl/internal/db"
	"github.com/bookandmusic/love-girl/internal/log"
)

func ProvideDB(
	cfg *config.AppConfig,
	gormLogger *log.GormLogger,
) (*gorm.DB, error) {
	// 直接创建数据库连接（默认使用 SQLite）
	dbCli, err := db.NewDBConn(cfg.DataSource.Database, gormLogger)
	if err != nil {
		return nil, err
	}

	return dbCli, nil
}
