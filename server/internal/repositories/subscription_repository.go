// REPOSITORY - repositories/admin_subscription.go
package repositories

import (
	"server/internal/models"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	CreateTier(tier *models.SubscriptionTier) error
	UpdateTier(tier *models.SubscriptionTier) error
	GetTierByID(id uint) (*models.SubscriptionTier, error)
	DeleteTier(id uint) error
	ResetAllUserToken() error
	FindAllUserSubscriptions() ([]models.UserSubscription, error)
	FindUserSubscriptionByID(userID string) (*models.UserSubscription, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewAdminSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db}
}

func (r *subscriptionRepository) CreateTier(tier *models.SubscriptionTier) error {
	return r.db.Create(tier).Error
}

func (r *subscriptionRepository) UpdateTier(tier *models.SubscriptionTier) error {
	return r.db.Model(&models.SubscriptionTier{}).Where("id = ?", tier.ID).Updates(tier).Error
}

func (r *subscriptionRepository) DeleteTier(id uint) error {
	return r.db.Delete(&models.SubscriptionTier{}, id).Error
}

func (r *subscriptionRepository) ResetAllUserToken() error {
	return r.db.Model(&models.UserSubscription{}).
		Where("is_active = ?", true).
		Update("remaining_tokens", gorm.Expr("token_limit")).Error
}

func (r *subscriptionRepository) FindAllUserSubscriptions() ([]models.UserSubscription, error) {
	var subs []models.UserSubscription
	err := r.db.Preload("User").Preload("SubscriptionTier").Find(&subs).Error
	return subs, err
}

func (r *subscriptionRepository) FindUserSubscriptionByID(userID string) (*models.UserSubscription, error) {
	var sub models.UserSubscription
	err := r.db.Preload("User").Preload("SubscriptionTier").
		Where("user_id = ?", userID).
		First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) GetTierByID(id uint) (*models.SubscriptionTier, error) {
	var tier models.SubscriptionTier
	if err := r.db.First(&tier, id).Error; err != nil {
		return nil, err
	}
	return &tier, nil
}
