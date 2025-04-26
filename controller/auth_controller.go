package controller

import (
	dto "simple-golang-tdd/dto"
	"simple-golang-tdd/service"
	utils "simple-golang-tdd/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthController handles authentication-related operations
type AuthController struct {
	authService  service.AuthService
	authvalidate *validator.Validate
}

func NewAuthController(service service.AuthService) *AuthController {
	validate := validator.New()
	return &AuthController{authService: service, authvalidate: validate}
}

// Login godoc
// @Summary      User Login
// @Description  Logs in a user and returns access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  dto.UserCredentials  true  "User Credentials"
// @Success      200  {object} dto.SuccessResponse  "login successful"
// @Failure      400  {object} dto.ErrorResponse
// @Failure      401  {object} dto.ErrorResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /customer/login [post]
func (ac *AuthController) Login(c *gin.Context) {
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

// Logout godoc
// @Summary      User Logout
// @Description  Logs out a user by invalidating the token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "Authorization Bearer Token"
// @Success      200  {object} dto.SuccessResponse
// @Failure      401  {object} dto.ErrorResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /customer/logout [post]
func (ac *AuthController) Logout(c *gin.Context) {
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

// RefreshToken godoc
// @Summary      Refresh Access Token
// @Description  Generates a new access token using a refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  dto.RefreshToken  true  "Refresh Token"
// @Success      200  {object} dto.AccessTokenResponse  "Token refreshed successfully"
// @Failure      400  {object} dto.ErrorResponse  "Invalid request body"
// @Failure      500  {object} dto.ErrorResponse "Internal server error"
// @Router       /customer/refresh-token [post]
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

	utils.SuccessResponse(c, 200, "token refreshed successfully", newAccessToken)
}
