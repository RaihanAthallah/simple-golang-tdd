package controller

import (
	"simple-golang-tdd/dto"
	"simple-golang-tdd/model"

	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock of the AuthService interface
type MockCustomerService struct {
	mock.Mock
}

// Payment mocks the Payment method of AuthService
func (m *MockCustomerService) Payment(request dto.PaymentRequest, username string) (model.Customer, error) {
    args := m.Called(request, username) // expects two arguments
    return args.Get(0).(model.Customer), args.Error(1)
}
