package services

import (
	"server/internal/dto"
	"server/internal/models"
	"server/internal/repositories"
)

type SubscriptionService interface {
	CreateTier(req *dto.CreateTierRequest) error
	UpdateTier(req *dto.UpdateTierRequest) error
	DeleteTier(id uint) error
	ResetAllUserTokens() error
	GetAllSubscriptions() ([]dto.UserSubscriptionResponse, error)
	GetSubscriptionByUserID(userID string) (*dto.UserSubscriptionResponse, error)
}

type subscriptionService struct {
	repo repositories.SubscriptionRepository
}

func NewAdminSubscriptionService(repo repositories.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo}
}

func (s *subscriptionService) CreateTier(req *dto.CreateTierRequest) error {
	tier := &models.SubscriptionTier{
		Name:        req.Name,
		TokenLimit:  req.TokenLimit,
		Price:       req.Price,
		Duration:    req.Duration,
		Description: req.Description,
	}
	return s.repo.CreateTier(tier)
}

func (s *subscriptionService) UpdateTier(req *dto.UpdateTierRequest) error {
	tier := &models.SubscriptionTier{
		ID:          req.ID,
		Name:        req.Name,
		TokenLimit:  req.TokenLimit,
		Price:       req.Price,
		Duration:    req.Duration,
		Description: req.Description,
	}
	return s.repo.UpdateTier(tier)
}

func (s *subscriptionService) DeleteTier(id uint) error {
	return s.repo.DeleteTier(id)
}

func (s *subscriptionService) ResetAllUserTokens() error {
	return s.repo.ResetAllUserToken()
}

func (s *subscriptionService) GetAllSubscriptions() ([]dto.UserSubscriptionResponse, error) {
	subs, err := s.repo.FindAllUserSubscriptions()
	if err != nil {
		return nil, err
	}

	var result []dto.UserSubscriptionResponse
	for _, s := range subs {
		result = append(result, dto.UserSubscriptionResponse{
			UserID:    s.User.ID.String(),
			Email:     s.User.Email,
			Fullname:  s.User.Fullname,
			Avatar:    s.User.Avatar,
			TierName:  s.SubscriptionTier.Name,
			IsActive:  s.IsActive,
			ExpiresAt: s.ExpiresAt.Format("2006-01-02"),
			Remaining: s.RemainingTokens,
		})
	}
	return result, nil
}

func (s *subscriptionService) GetSubscriptionByUserID(userID string) (*dto.UserSubscriptionResponse, error) {
	sub, err := s.repo.FindUserSubscriptionByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserSubscriptionResponse{
		UserID:    sub.User.ID.String(),
		Email:     sub.User.Email,
		Fullname:  sub.User.Fullname,
		Avatar:    sub.User.Avatar,
		TierName:  sub.SubscriptionTier.Name,
		IsActive:  sub.IsActive,
		ExpiresAt: sub.ExpiresAt.Format("2006-01-02"),
		Remaining: sub.RemainingTokens,
	}, nil
}
