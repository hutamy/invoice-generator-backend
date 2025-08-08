package repositories

import (
	"github.com/hutamy/invoice-generator-backend/dto"
	"github.com/hutamy/invoice-generator-backend/models"
	"gorm.io/gorm"
)

type ClientRepository interface {
	CreateClient(client *models.Client) error
	GetAllByUserID(userID uint) ([]models.Client, error)
	GetAllByUserIDWithPagination(req dto.GetClientsRequest) ([]models.Client, int64, error)
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
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&clients).Error
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *clientRepository) GetAllByUserIDWithPagination(req dto.GetClientsRequest) ([]models.Client, int64, error) {
	var clients []models.Client
	var totalItems int64

	query := r.db.Where("user_id = ?", req.UserID)

	// Add search functionality if search term is provided
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ? OR address ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// Count total items
	if err := query.Model(&models.Client{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&clients).Error

	if err != nil {
		return nil, 0, err
	}

	return clients, totalItems, nil
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
