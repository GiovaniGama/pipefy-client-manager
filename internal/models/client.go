package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	Name        string  `json:"nome" gorm:"not null"`
	Email       string  `json:"email" gorm:"uniqueIndex;not null"`
	RequestType string  `json:"tipo_solicitacao" gorm:"not null"`
	AssetValue  float64 `json:"valor_patrimonio" gorm:"not null"`
	Status      string  `json:"status" gorm:"default:'Aguardando Análise'"`
	Priority    string  `json:"prioridade"`
}
