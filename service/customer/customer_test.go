package service

import (
	"errors"
	"simple-golang-tdd/dto"
	"simple-golang-tdd/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerService_Payment_Success(t *testing.T) {
	mockCustomerRepository := new(MockCustomerRepository)
	mockMerchantRepository := new(MockMerchantRepository)
	customerService := NewCustomerService(mockCustomerRepository, mockMerchantRepository)

	fakePayment := dto.PaymentRequest{
		MerchantID: "merchant123",
		Amount:     100.0,
	}

	fakeUsername := "testuser"

	expectedCustomer := model.Customer{
		ID:       "1",
		Name:     "testuser",
		Username: "testuser",
		Password: "password", // password cocok
		Balance:  1000.0,     // saldo awal
	}

	expectedMerchant := model.Merchant{
		ID:      "merchant123",
		Name:    "Merchant 123",
		Balance: 600.0, // saldo awal merchant
	}

	mockCustomerRepository.On("GetUserByUsername", fakeUsername).Return(expectedCustomer, nil)
	mockMerchantRepository.On("GetMerchantBalance", fakePayment.MerchantID).Return(500.0, nil)
	mockCustomerRepository.On("UpdateUserBalance", expectedCustomer.ID, expectedCustomer.Balance-100.0).Return(expectedCustomer, nil)
	mockMerchantRepository.On("UpdateMerchantBalance", fakePayment.MerchantID, 500.0+100.0).Return(expectedMerchant, nil)
	resp, err := customerService.Payment(fakePayment, fakeUsername)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	mockCustomerRepository.AssertExpectations(t)
}

func TestCustomerService_Payment_UserNotFound(t *testing.T) {
	mockCustomerRepository := new(MockCustomerRepository)
	mockMerchantRepository := new(MockMerchantRepository)
	customerService := NewCustomerService(mockCustomerRepository, mockMerchantRepository)

	fakePayment := dto.PaymentRequest{
		MerchantID: "merchant123",
		Amount:     100.0,
	}

	fakeUsername := "testuser"

	mockCustomerRepository.On("GetUserByUsername", fakeUsername).Return(model.Customer{}, errors.New("user not found"))

	resp, err := customerService.Payment(fakePayment, fakeUsername)

	assert.Error(t, err)
	assert.Empty(t, resp.ID)
	assert.Empty(t, resp.Username)
	mockCustomerRepository.AssertExpectations(t)
}

func TestCustomerService_Payment_MerchantNotFound(t *testing.T) {
	mockCustomerRepository := new(MockCustomerRepository)
	mockMerchantRepository := new(MockMerchantRepository)
	customerService := NewCustomerService(mockCustomerRepository, mockMerchantRepository)

	fakePayment := dto.PaymentRequest{
		MerchantID: "merchant123",
		Amount:     100.0,
	}

	fakeUsername := "testuser"

	expectedCustomer := model.Customer{
		ID:       "1",
		Name:     "testuser",
		Username: "testuser",
		Password: "password",
		Balance:  1000.0,
	}

	mockCustomerRepository.On("GetUserByUsername", fakeUsername).Return(expectedCustomer, nil)
	mockMerchantRepository.On("GetMerchantBalance", fakePayment.MerchantID).Return(0.0, errors.New("merchant not found"))

	resp, err := customerService.Payment(fakePayment, fakeUsername)

	assert.Error(t, err)
	assert.Empty(t, resp.ID)
	assert.Empty(t, resp.Username)
	mockCustomerRepository.AssertExpectations(t)
	mockMerchantRepository.AssertExpectations(t)
}

func TestCustomerService_Payment_InsufficientBalance(t *testing.T) {
	mockCustomerRepository := new(MockCustomerRepository)
	mockMerchantRepository := new(MockMerchantRepository)
	customerService := NewCustomerService(mockCustomerRepository, mockMerchantRepository)

	fakePayment := dto.PaymentRequest{
		MerchantID: "merchant123",
		Amount:     2000.0, // lebih besar dari saldo
	}

	fakeUsername := "testuser"

	expectedCustomer := model.Customer{
		ID:       "1",
		Name:     "testuser",
		Username: "testuser",
		Password: "password", // password cocok
		Balance:  1000.0,     // saldo awal
	}

	mockCustomerRepository.On("GetUserByUsername", fakeUsername).Return(expectedCustomer, nil)
	mockMerchantRepository.On("GetMerchantBalance", fakePayment.MerchantID).Return(500.0, nil)

	resp, err := customerService.Payment(fakePayment, fakeUsername)

	assert.Error(t, err)
	assert.Empty(t, resp.ID)
	assert.Empty(t, resp.Username)
	mockCustomerRepository.AssertExpectations(t)
}
