package controller_test

import (
	"simple-golang-tdd/dto"

	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock of the AuthService interface
type MockAuthService struct {
	mock.Mock
}

// Login mocks the Login method of AuthService
func (m *MockAuthService) Login(credentials dto.UserCredentials) (dto.AuthResponse, error) {
	args := m.Called(credentials)
	return args.Get(0).(dto.AuthResponse), args.Error(1)
}

// Logout mocks the Logout method of AuthService
func (m *MockAuthService) Logout(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

// RefreshToken mocks the RefreshToken method of AuthService
func (m *MockAuthService) RefreshToken(token dto.RefreshToken) (dto.AccessTokenResponse, error) {
	args := m.Called(token)
	return args.Get(0).(dto.AccessTokenResponse), args.Error(1)
}
