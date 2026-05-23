package models

import "time"

type ProcessedEvent struct {
	EventID     string    `json:"event_id" gorm:"primaryKey"`
	ProcessedAt time.Time `json:"processado_em"`
}
