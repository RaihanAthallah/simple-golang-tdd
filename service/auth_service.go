package service

import (
	"simple-golang-tdd/dto"
	"simple-golang-tdd/repository"
)

type AuthService interface {
	Login(dto.UserCredentials) (dto.AuthResponse, error)
	Logout(string) error
	RefreshToken(token dto.RefreshToken) (dto.AccessTokenResponse, error)
}

type authServiceImpl struct {
	// authRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) AuthService {
	return &authServiceImpl{}
}

func (s *authServiceImpl) Login(credentials dto.UserCredentials) (dto.AuthResponse, error) {
	var tokens dto.AuthResponse
	tokens.AccessToken = "dummy_access_token"
	tokens.RefreshToken = "dummy_refresh_token"

	return tokens, nil
}

func (s *authServiceImpl) Logout(token string) error {
	return nil
}

func (s *authServiceImpl) RefreshToken(token dto.RefreshToken) (dto.AccessTokenResponse, error) {
	var newToken dto.AccessTokenResponse

	newToken.AccessToken = "new_dummy_token"

	return newToken, nil
}
