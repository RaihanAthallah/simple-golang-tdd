package controller

import (
	dto "simple-golang-tdd/dto"
	"simple-golang-tdd/service"
	utils "simple-golang-tdd/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService  service.AuthService
	authvalidate *validator.Validate
}

func NewAuthController(service service.AuthService) *AuthController {
	validate := validator.New()
	return &AuthController{authService: service, authvalidate: validate}
}

func (ac *AuthController) Login(c *gin.Context) {
	// fetch username and password from request body
	var userCredentials dto.UserCredentials
	err := c.BindJSON(&userCredentials)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid request body")
		return
	}

	if err := ac.authvalidate.Struct(userCredentials); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	tokens, err := ac.authService.Login(userCredentials)
	if err != nil {
		utils.ErrorResponse(c, 401, "invalid credentials")
		return
	}

	utils.SuccessResponse(c, 200, "login successful", tokens)
}

func (ac *AuthController) Logout(c *gin.Context) {
	// fetch token from request header
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		utils.ErrorResponse(c, 401, "please provide a token")
		return
	}

	err := ac.authService.Logout(token)
	if err != nil {
		utils.ErrorResponse(c, 500, "internal server error")
		return
	}

	utils.SuccessResponse(c, 200, "logout successful", nil)
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	var refreshToken dto.RefreshToken
	err := c.BindJSON(&refreshToken)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid request body")
		return
	}

	if err := ac.authvalidate.Struct(refreshToken); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	newAccessToken, err := ac.authService.RefreshToken(refreshToken)
	if err != nil {
		utils.ErrorResponse(c, 500, "internal server error")
		return
	}

	// newAccessToken := dto.AccessTokenResponse{AccessToken: "new_dummy_token"}
	utils.SuccessResponse(c, 200, "token refreshed successfully", newAccessToken)

}
