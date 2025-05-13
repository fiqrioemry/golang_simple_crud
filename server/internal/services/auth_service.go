package services

import (
	"errors"
	"fmt"
	"server/internal/config"
	"server/internal/dto"
	"server/internal/models"
	"server/internal/utils"
	"time"

	"server/internal/repositories"
)

type AuthService interface {
	SendOTP(req *dto.SendOTPRequest) error
	VerifyOTPCode(req *dto.VerifyOTPRequest) error
	GetUserInfo(userID string) (*dto.UserInfoResponse, error)
	UserLogin(req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshUserToken(refreshToken string) (*dto.AuthResponse, error)
	UserRegister(req *dto.RegisterRequest) (*dto.AuthResponse, error)
}

type authService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) AuthService {
	return &authService{repo}
}

func (s *authService) UserRegister(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Fullname: req.Fullname,
		Avatar:   utils.RandomUserAvatar(req.Fullname),
	}

	if err := s.repo.CreateNewUser(user); err != nil {
		return nil, err
	}
	userID := user.ID.String()

	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}
	newTokenModel := &models.Token{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.repo.StoreRefreshToken(newTokenModel); err != nil {
		return nil, err
	}
	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) UserLogin(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid Email")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid Password")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	tokenModel := &models.Token{
		UserID:    user.ID,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.repo.StoreRefreshToken(tokenModel); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *authService) SendOTP(req *dto.SendOTPRequest) error {
	if _, err := s.repo.GetUserByEmail(req.Email); err == nil {
		return errors.New("email already registered")
	}

	otp := utils.GenerateOTP(6)
	subject := "one time password (OTP) XXXXXX"
	body := fmt.Sprintf("your OTP code is %s", otp)

	if err := utils.SendEmail(subject, req.Email, otp, body); err != nil {
		return errors.New("failed to send email")
	}

	err := config.RedisClient.Set(config.Ctx, "otp:"+req.Email, otp, 5*time.Minute).Err()
	return err

}

func (s *authService) VerifyOTPCode(req *dto.VerifyOTPRequest) error {
	savedOTP, err := config.RedisClient.Get(config.Ctx, "otp:"+req.Email).Result()
	if err != nil {
		return errors.New("otp expired or invalid")
	}

	if savedOTP != req.OTP {
		return errors.New("invalid OTP code")
	}

	config.RedisClient.Del(config.Ctx, "otp:"+req.Email)
	return nil
}

func (s *authService) GetUserInfo(userID string) (*dto.UserInfoResponse, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserInfoResponse{
		UserID:   user.ID.String(),
		Email:    user.Email,
		Fullname: user.Fullname,
		Avatar:   user.Avatar,
	}, nil
}

func (s *authService) RefreshUserToken(refreshToken string) (*dto.AuthResponse, error) {

	_, err := utils.DecodeRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("unauthorized ! invalid credential")
	}

	token, err := s.repo.FindRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("unauthorized ! invalid refresh token")
	}

	if token.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	user, err := s.repo.GetUserByID(token.UserID.String())
	if err != nil {
		return nil, errors.New("user not found")
	}

	newAccessToken, err := utils.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	newToken := &models.Token{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.repo.StoreRefreshToken(newToken); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
