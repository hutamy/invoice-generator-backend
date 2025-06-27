package repositories

import (
	"errors"

	"github.com/hutamy/invoice-generator-backend/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(user *models.User) error {
	return r.db.Create(&user).Error
}

func (r *authRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *authRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	res := r.db.First(&user, id)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (r *authRepository) UpdateUser(user *models.User) error {
	res := r.db.Save(user)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
