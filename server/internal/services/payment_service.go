package services

import (
	"errors"
	"server/internal/config"
	"server/internal/dto"
	"server/internal/models"
	"server/internal/repositories"
	"server/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentService interface {
	CreatePayment(userID string, req dto.CreatePaymentRequest) (*dto.CreatePaymentResponse, error)
	HandlePaymentNotification(req dto.MidtransNotificationRequest) error
	GetAllUserPayments(query string, page, limit int) (*dto.PaymentListResponse, error)
	GetPaymentByID(id string) (*dto.PaymentDetailResponse, error)
}

type paymentService struct {
	repo     repositories.PaymentRepository
	tierRepo repositories.SubscriptionRepository
	authRepo repositories.AuthRepository
}

func NewPaymentService(
	repo repositories.PaymentRepository,
	tierRepo repositories.SubscriptionRepository,
	authRepo repositories.AuthRepository,
) PaymentService {
	return &paymentService{repo, tierRepo, authRepo}
}

func (s *paymentService) CreatePayment(userID string, req dto.CreatePaymentRequest) (*dto.CreatePaymentResponse, error) {
	tier, err := s.tierRepo.GetTierByID(req.TierID)
	if err != nil {
		return nil, errors.New("subscription tier not found")
	}

	tax := tier.Price * utils.GetTaxRate()
	total := tier.Price + tax
	paymentID := uuid.New()

	payment := &models.Payment{
		ID:       paymentID,
		UserID:   uuid.MustParse(userID),
		TierID:   tier.ID,
		Subtotal: tier.Price,
		Tax:      tax,
		Total:    total,
		Method:   req.Method,
		Status:   "pending",
	}

	if err := s.repo.CreatePayment(payment); err != nil {
		return nil, err
	}

	user, err := s.authRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  paymentID.String(),
			GrossAmt: int64(total),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: user.Email,
		},
	}

	snapResp, err := config.SnapClient.CreateTransaction(snapReq)

	return &dto.CreatePaymentResponse{
		PaymentID: paymentID.String(),
		SnapToken: snapResp.Token,
		SnapURL:   snapResp.RedirectURL,
	}, nil
}

func (s *paymentService) HandlePaymentNotification(req dto.MidtransNotificationRequest) error {
	payment, err := s.repo.GetPaymentByOrderID(req.OrderID)
	if err != nil {
		return err
	}

	if payment.Status == "success" {
		// avoid duplicate processing
		return nil
	}

	payment.Method = req.PaymentType

	switch req.TransactionStatus {
	case "settlement", "capture":
		if req.FraudStatus == "accept" || req.FraudStatus == "" {
			payment.Status = "success"
			now := time.Now()
			payment.PaidAt = &now

			// Buat user subscription setelah pembayaran sukses
			tier, err := s.tierRepo.GetTierByID(payment.TierID)
			if err != nil {
				return err
			}

			subscription := &models.UserSubscription{
				UserID:             payment.UserID,
				SubscriptionTierID: tier.ID,
				StartedAt:          now,
				ExpiresAt:          now.AddDate(0, 0, tier.Duration), // durasi dalam hari
				IsActive:           true,
				RemainingTokens:    tier.TokenLimit,
			}

			if err := s.repo.CreateUserSubscription(subscription); err != nil {
				return err
			}
		}
	case "pending":
		payment.Status = "pending"
	default:
		payment.Status = "failed"
	}

	return s.repo.UpdatePayment(payment)
}

func (s *paymentService) GetPaymentByID(id string) (*dto.PaymentDetailResponse, error) {
	p, err := s.repo.GetPaymentByID(id)
	if err != nil {
		return nil, err
	}
	return &dto.PaymentDetailResponse{
		ID:       p.ID.String(),
		UserID:   p.UserID.String(),
		TierID:   p.TierID,
		Subtotal: p.Subtotal,
		Tax:      p.Tax,
		Total:    p.Total,
		Method:   p.Method,
		Status:   p.Status,
		PaidAt:   p.PaidAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *paymentService) GetAllUserPayments(query string, page, limit int) (*dto.PaymentListResponse, error) {
	offset := (page - 1) * limit

	payments, total, err := s.repo.GetAllUserPayments(query, limit, offset)
	if err != nil {
		return nil, err
	}

	var results []dto.PaymentResponse
	for _, p := range payments {
		results = append(results, dto.PaymentResponse{
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

	return &dto.PaymentListResponse{
		Payments: results,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}
