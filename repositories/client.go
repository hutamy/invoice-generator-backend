package repositories

import (
	"github.com/hutamy/invoice-generator/models"
	"gorm.io/gorm"
)

type ClientRepository interface {
	CreateClient(client *models.Client) error
	GetAllByUserID(userID uint) ([]models.Client, error)
	GetClientByID(id, userID uint) (*models.Client, error)
	UpdateClient(client *models.Client) error
	DeleteClient(id, userID uint) error
}

type clientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) CreateClient(client *models.Client) error {
	return r.db.Create(client).Error
}

func (r *clientRepository) GetAllByUserID(userID uint) ([]models.Client, error) {
	var clients []models.Client
	err := r.db.Where("user_id = ?", userID).Find(&clients).Error
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *clientRepository) GetClientByID(id, userID uint) (*models.Client, error) {
	var client models.Client
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&client).Error
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *clientRepository) UpdateClient(client *models.Client) error {
	return r.db.Save(client).Error
}

func (r *clientRepository) DeleteClient(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Client{}).Error
}
