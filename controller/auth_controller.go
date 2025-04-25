package controller

import "simple-golang-tdd/service"

type AuthController struct {
	authService service.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{}
}
