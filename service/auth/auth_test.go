package service

import (
	"errors"
	"simple-golang-tdd/dto"
	"simple-golang-tdd/model"
	"simple-golang-tdd/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthService_Login_Success(t *testing.T) {
	mockDependencies := new(MockCustomerRepository)
	authService := NewAuthService(mockDependencies)

	credentials := dto.UserCredentials{
		Username: "testuser",
		Password: "password",
	}

	expectedCustomer := model.Customer{
		ID:       "1",
		Username: "testuser",
		Password: "password", // password cocok
	}

	mockDependencies.On("GetUserByUsername", credentials.Username).Return(expectedCustomer, nil)

	resp, err := authService.Login(credentials)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	mockDependencies.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockDependencies := new(MockCustomerRepository)
	authService := NewAuthService(mockDependencies)

	credentials := dto.UserCredentials{
		Username: "unknownuser",
		Password: "password",
	}

	mockDependencies.On("GetUserByUsername", credentials.Username).Return(model.Customer{}, errors.New("user not found"))

	resp, err := authService.Login(credentials)

	assert.Error(t, err)
	assert.Empty(t, resp.AccessToken)
	assert.Empty(t, resp.RefreshToken)
	mockDependencies.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockDependencies := new(MockCustomerRepository)
	authService := NewAuthService(mockDependencies)

	credentials := dto.UserCredentials{
		Username: "johndoe",
		Password: "wrongpassword",
	}

	expectedCustomer := model.Customer{
		ID:       "cust-001",
		Username: "johndoe",
		Password: "password123", // password mismatch
	}

	mockDependencies.On("GetUserByUsername", credentials.Username).Return(expectedCustomer, nil)

	resp, err := authService.Login(credentials)

	assert.Error(t, err)
	assert.Empty(t, resp.AccessToken)
	assert.Empty(t, resp.RefreshToken)
	mockDependencies.AssertExpectations(t)
}

func TestAuthService_Logout_Success(t *testing.T) {
	mockDependencies := new(MockCustomerRepository)
	authService := NewAuthService(mockDependencies)

	userID := "1"
	fakeAccessToken, _ := utils.GenerateAccessToken(userID)
	// Logout biasanya cukup validasi, karena di mock ini tidak ada repo untuk logout,
	// kita anggap selalu berhasil.
	err := authService.Logout(fakeAccessToken)

	assert.NoError(t, err)
}

func TestAuthService_Logout_EmptyToken(t *testing.T) {
	mockDependencies := new(MockCustomerRepository)
	authService := NewAuthService(mockDependencies)

	fakeAccessToken := "" // Empty token assumed invalid

	// Misal kalau userID kosong dianggap error
	err := authService.Logout(fakeAccessToken)

	assert.Error(t, err)
}

func TestAuthService_Logout_InvalidTokenFormat(t *testing.T) {
	mockDependencies := new(MockCustomerRepository)
	authService := NewAuthService(mockDependencies)

	invalidToken := "invalid-token-format" // Invalid token

	err := authService.Logout(invalidToken)

	assert.Error(t, err)
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	mockDependencies := new(MockCustomerRepository)
	authService := NewAuthService(mockDependencies)
	fakeRefreshToken, _ := utils.GenerateRefreshToken("1")

	refreshReq := dto.RefreshToken{
		RefreshToken: fakeRefreshToken,
	}

	// Dalam test ini asumsi refresh token valid -> dapat access token baru
	resp, err := authService.RefreshToken(refreshReq)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	mockDependencies := new(MockCustomerRepository)
	authService := NewAuthService(mockDependencies)

	refreshReq := dto.RefreshToken{
		RefreshToken: "", // Empty token assumed invalid
	}

	resp, err := authService.RefreshToken(refreshReq)

	assert.Error(t, err)
	assert.Empty(t, resp.AccessToken)
}

