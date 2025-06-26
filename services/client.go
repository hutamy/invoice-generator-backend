package services

import (
	"github.com/hutamy/invoice-generator-backend/dto"
	"github.com/hutamy/invoice-generator-backend/models"
	"github.com/hutamy/invoice-generator-backend/repositories"
)

type ClientService interface {
	CreateClient(req dto.CreateClientRequest) error
	GetAllClientsByUserID(userID uint) ([]models.Client, error)
	GetClientByID(id, userID uint) (*models.Client, error)
	UpdateClient(req dto.UpdateClientRequest) error
	DeleteClient(id, userID uint) error
}

type clientService struct {
	clientRepo repositories.ClientRepository
}

func NewClientService(clientRepo repositories.ClientRepository) ClientService {
	return &clientService{clientRepo: clientRepo}
}

func (s *clientService) CreateClient(req dto.CreateClientRequest) error {
	client := &models.Client{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
		UserID:  req.UserID,
	}
	return s.clientRepo.CreateClient(client)
}

func (s *clientService) GetAllClientsByUserID(userID uint) ([]models.Client, error) {
	return s.clientRepo.GetAllByUserID(userID)
}

func (s *clientService) GetClientByID(id, userID uint) (*models.Client, error) {
	return s.clientRepo.GetClientByID(id, userID)
}

func (s *clientService) UpdateClient(req dto.UpdateClientRequest) error {
	client, err := s.GetClientByID(req.ID, req.UserID)
	if err != nil {
		return err
	}

	if req.Name != nil {
		client.Name = *req.Name
	}

	if req.Email != nil {
		client.Email = *req.Email
	}

	if req.Address != nil {
		client.Address = *req.Address
	}

	if req.Phone != nil {
		client.Phone = *req.Phone
	}

	return s.clientRepo.UpdateClient(client)
}

func (s *clientService) DeleteClient(id, userID uint) error {
	return s.clientRepo.DeleteClient(id, userID)
}
