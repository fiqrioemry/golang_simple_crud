package repositories

import (
	"server/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(userID string) (*models.User, error)
	Update(user *models.User) error
	GetUserSubscription(userID string) (*models.UserSubscription, error)
	UpdateUserTokenUsage(userID string, used int) error
	GetPayments(userID string) ([]models.Payment, error)
	GetFormsByUser(userID string) ([]models.Form, error)
	GetFormDetail(formID string) (*models.Form, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetByID(userID string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", userID).Error
	return &user, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Model(&models.User{}).Where("id = ?", user.ID).Updates(user).Error
}

func (r *userRepository) GetUserSubscription(userID string) (*models.UserSubscription, error) {
	var sub models.UserSubscription
	err := r.db.Preload("SubscriptionTier").First(&sub, "user_id = ?", userID).Error
	return &sub, err
}

func (r *userRepository) UpdateUserTokenUsage(userID string, used int) error {
	return r.db.Model(&models.UserSubscription{}).
		Where("user_id = ?", userID).
		Update("remaining_tokens", gorm.Expr("remaining_tokens - ?", used)).Error
}

func (r *userRepository) GetPayments(userID string) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Preload("Tier").Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

func (r *userRepository) GetFormsByUser(userID string) ([]models.Form, error) {
	var forms []models.Form
	err := r.db.Where("user_id = ?", userID).Find(&forms).Error
	return forms, err
}

func (r *userRepository) GetFormDetail(formID string) (*models.Form, error) {
	var form models.Form
	err := r.db.First(&form, "id = ?", formID).Error
	return &form, err
}
