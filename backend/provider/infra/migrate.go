package infra

import (
	"gorm.io/gorm"

	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/model"
)

func ProvideMigrate(db *gorm.DB, logger *log.Logger) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.File{},
		&model.Album{},
		&model.Moment{},
		&model.Place{},
		&model.EntityFile{},
		&model.Anniversary{},
		&model.Setting{},
		&model.Comment{},
		&model.Notification{},
	); err != nil {
		logger.Error("Database migration failed:", "error", err)
		return err
	}

	logger.Info("Database migrated successfully")
	return nil
}
