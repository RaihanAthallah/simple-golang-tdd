package service

import (
	"simple-golang-tdd/model"

	"github.com/stretchr/testify/mock"
)

// MockCustomerRepository is a mock of the CustomerRepository interface
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) GetUserByUsername(username string) (model.Customer, error) {
	args := m.Called(username)
	return args.Get(0).(model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetUserByID(id string) (model.Customer, error) {
	args := m.Called(id)
	return args.Get(0).(model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetUserBalance(id string) (float64, error) {
	args := m.Called(id)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCustomerRepository) UpdateUserBalance(id string, amount float64) (model.Customer, error) {
	args := m.Called(id, amount)
	return args.Get(0).(model.Customer), args.Error(1)
}

type MockMerchantRepository struct {
	mock.Mock
}

func (m *MockMerchantRepository) UpdateMerchantBalance(id string, amount float64) (model.Merchant, error) {
	args := m.Called(id, amount)
	return args.Get(0).(model.Merchant), args.Error(1)
}

func (m *MockMerchantRepository) GetMerchantBalance(id string) (float64, error) {
	args := m.Called(id)
	return args.Get(0).(float64), args.Error(1)
}
