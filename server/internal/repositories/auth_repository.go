package repositories

import (
	"server/internal/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateNewUser(user *models.User) error
	DeleteRefreshToken(token string) error
	StoreRefreshToken(token *models.Token) error
	GetUserByID(userID string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	FindRefreshToken(token string) (*models.Token, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) CreateNewUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.Token{}).Error
}

func (r *authRepository) StoreRefreshToken(token *models.Token) error {
	return r.db.Create(token).Error
}

func (r *authRepository) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindRefreshToken(token string) (*models.Token, error) {
	var t models.Token
	if err := r.db.Where("token = ?", token).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
