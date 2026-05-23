package repository

import (
	"time"

	"github.com/GiovaniGama/pipefy-client-manager/internal/models"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Exists(eventID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProcessedEvent{}).Where("event_id = ?", eventID).Count(&count).Error
	return count > 0, err
}

func (r *EventRepository) Save(eventID string) error {
	event := models.ProcessedEvent{
		EventID:     eventID,
		ProcessedAt: time.Now(),
	}
	return r.db.Create(&event).Error
}
