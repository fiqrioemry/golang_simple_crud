package services

import (
	"server/internal/dto"
	"server/internal/repositories"
	"server/internal/utils"
)

type UserService interface {
	GetProfile(userID string) (*dto.UserProfileResponse, error)
	UpdateProfile(userID string, req *dto.UpdateProfileRequest) error
	GetMySubscription(userID string) (*dto.UserSubscriptionResponse, error)
	UpdateUserSubscriptionToken(userID string, used int) error
	GetTransactionHistory(userID string) ([]dto.PaymentResponse, error)
	GetMyForms(userID string) ([]dto.MyFormResponse, error)
	GetMyFormDetail(formID string) (*dto.MyFormDetailResponse, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetProfile(userID string) (*dto.UserProfileResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return &dto.UserProfileResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Fullname: user.Fullname,
		Avatar:   user.Avatar,
	}, nil
}

func (s *userService) UpdateProfile(userID string, req *dto.UpdateProfileRequest) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	user.Fullname = req.Fullname

	if req.Avatar != "" && req.Avatar != user.Avatar {

		if user.Avatar != "" && !utils.IsDiceBear(user.Avatar) {
			_ = utils.DeleteFromCloudinary(user.Avatar)
		}
		user.Avatar = req.Avatar
	}

	return s.repo.Update(user)
}

func (s *userService) GetMySubscription(userID string) (*dto.UserSubscriptionResponse, error) {
	sub, err := s.repo.GetUserSubscription(userID)
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

func (s *userService) UpdateUserSubscriptionToken(userID string, used int) error {
	return s.repo.UpdateUserTokenUsage(userID, used)
}

func (s *userService) GetTransactionHistory(userID string) ([]dto.PaymentResponse, error) {
	payments, err := s.repo.GetPayments(userID)
	if err != nil {
		return nil, err
	}
	var result []dto.PaymentResponse
	for _, p := range payments {
		result = append(result, dto.PaymentResponse{
			ID:            p.ID.String(),
			UserID:        p.UserID.String(),
			UserEmail:     p.User.Email,
			Fullname:      p.User.Fullname,
			TierID:        p.TierID,
			TierName:      p.Tier.Name,
			Subtotal:      p.Subtotal,
			Tax:           p.Tax,
			Total:         p.Total,
			PaymentMethod: p.Method,
			Status:        p.Status,
			PaidAt:        p.PaidAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}

func (s *userService) GetMyForms(userID string) ([]dto.MyFormResponse, error) {
	forms, err := s.repo.GetFormsByUser(userID)
	if err != nil {
		return nil, err
	}
	var result []dto.MyFormResponse
	for _, f := range forms {
		result = append(result, dto.MyFormResponse{
			ID:    f.ID.String(),
			Title: f.Title,
			Type:  f.Type,
		})
	}
	return result, nil
}

func (s *userService) GetMyFormDetail(formID string) (*dto.MyFormDetailResponse, error) {
	form, err := s.repo.GetFormDetail(formID)
	if err != nil {
		return nil, err
	}
	return &dto.MyFormDetailResponse{
		ID:          form.ID.String(),
		Title:       form.Title,
		Description: form.Description,
		Type:        form.Type,
		CreatedAt:   form.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
