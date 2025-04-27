package service

import (
	"fmt"
	"simple-golang-tdd/dto"
	customerRepo "simple-golang-tdd/repository/customer"
	"simple-golang-tdd/utils"
)

type AuthService interface {
	Login(dto.UserCredentials) (dto.AuthResponse, error)
	Logout(string) error
	RefreshToken(token dto.RefreshToken) (dto.AccessTokenResponse, error)
}

type authServiceImpl struct {
	customerRepository customerRepo.CustomerRepository
}

func NewAuthService(customerRepository customerRepo.CustomerRepository) AuthService {
	return &authServiceImpl{customerRepository: customerRepository}
}

func (s *authServiceImpl) Login(credentials dto.UserCredentials) (dto.AuthResponse, error) {
	var tokens dto.AuthResponse

	customer, err := s.customerRepository.GetUserByUsername(credentials.Username)
	if err != nil {
		return tokens, fmt.Errorf("failed to get data user: %w", err)
	}

	if customer.Password != credentials.Password {
		return tokens, fmt.Errorf("username or password is incorrect")
	}

	// Generate tokens (dummy implementation)
	accessToken, err := utils.GenerateAccessToken(customer.Username)
	if err != nil {
		return tokens, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(customer.Username)
	if err != nil {
		return tokens, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return tokens, nil
}

func (s *authServiceImpl) Logout(token string) error {
	if _, err := utils.ValidateToken(token, "access"); err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	return nil
}

func (s *authServiceImpl) RefreshToken(token dto.RefreshToken) (dto.AccessTokenResponse, error) {
	var newToken dto.AccessTokenResponse

	tokenString, err := utils.NewAccessToken(token.RefreshToken)
	if err != nil {
		return newToken, fmt.Errorf("failed to generate new access token: %w", err)
	}

	newToken.AccessToken = tokenString

	return newToken, nil
}
