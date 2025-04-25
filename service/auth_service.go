package service

import "simple-golang-tdd/repository"

type AuthService interface {
}

type authServiceImpl struct {
	authRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) AuthService {
	return &authServiceImpl{authRepository: authRepository}
}
