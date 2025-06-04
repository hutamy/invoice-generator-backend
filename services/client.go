package services

import (
	"github.com/hutamy/invoice-generator/models"
	"github.com/hutamy/invoice-generator/repositories"
)

type ClientService interface {
	CreateClient(client *models.Client) error
	GetAllClientsByUserID(userID uint) ([]models.Client, error)
	GetClientByID(id, userID uint) (*models.Client, error)
	UpdateClient(client *models.Client) error
	DeleteClient(id, userID uint) error
}

type clientService struct {
	clientRepo repositories.ClientRepository
}

func NewClientService(clientRepo repositories.ClientRepository) ClientService {
	return &clientService{clientRepo: clientRepo}
}

func (s *clientService) CreateClient(client *models.Client) error {
	return s.clientRepo.CreateClient(client)
}

func (s *clientService) GetAllClientsByUserID(userID uint) ([]models.Client, error) {
	return s.clientRepo.GetAllByUserID(userID)
}

func (s *clientService) GetClientByID(id, userID uint) (*models.Client, error) {
	return s.clientRepo.GetByID(id, userID)
}

func (s *clientService) UpdateClient(client *models.Client) error {
	return s.clientRepo.UpdateClient(client)
}

func (s *clientService) DeleteClient(id, userID uint) error {
	return s.clientRepo.DeleteClient(id, userID)
}
