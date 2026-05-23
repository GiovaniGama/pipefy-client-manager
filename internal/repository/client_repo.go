package repository

import (
	"github.com/seu-usuario/pipefy-client-manager/internal/models"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) Create(client *models.Client) error {
	return r.db.Create(client).Error
}

func (r *ClientRepository) FindByEmail(email string) (*models.Client, error) {
	var client models.Client
	err := r.db.Where("email = ?", email).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) Update(client *models.Client) error {
	return r.db.Save(client).Error
}
