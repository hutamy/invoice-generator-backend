package services

import (
	"github.com/hutamy/invoice-generator/models"
	"github.com/hutamy/invoice-generator/repositories"
	"github.com/hutamy/invoice-generator/utils"
	"github.com/hutamy/invoice-generator/utils/errors"
)

type AuthService interface {
	SignUp(name, email, password string) error
	SignIn(email, password string) (models.User, error)
}

type authService struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) SignUp(name, email, password string) error {
	existingUser, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.authRepo.CreateUser(user)
}

func (s *authService) SignIn(email, password string) (models.User, error) {
	user, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		return models.User{}, err
	}
	if user == nil {
		return models.User{}, errors.ErrLoginFailed
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return models.User{}, errors.ErrLoginFailed
	}

	return *user, nil
}
