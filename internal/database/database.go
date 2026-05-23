package database

import (
	"github.com/GiovaniGama/pipefy-client-manager/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Client{}, &models.ProcessedEvent{})
}
