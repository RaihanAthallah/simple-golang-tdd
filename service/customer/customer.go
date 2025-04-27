package service

import (
	"fmt"
	"simple-golang-tdd/model"

	"simple-golang-tdd/dto"
	customerRepo "simple-golang-tdd/repository/customer"
	merchantRepo "simple-golang-tdd/repository/merchant"
)

type CustomerService interface {
	Payment(request dto.PaymentRequest, username string) (model.Customer, error)
}

type customerServiceImpl struct {
	customerRepository customerRepo.CustomerRepository
	merchantRepository merchantRepo.MerchantRepository
}

func NewCustomerService(customerRepository customerRepo.CustomerRepository, merchantRepository merchantRepo.MerchantRepository) CustomerService {
	return &customerServiceImpl{
		customerRepository: customerRepository,
		merchantRepository: merchantRepository}
}

func (s *customerServiceImpl) Payment(request dto.PaymentRequest, username string) (model.Customer, error) {
	var customer model.Customer
	customer, err := s.customerRepository.GetUserByUsername(username)
	if err != nil {
		return model.Customer{}, fmt.Errorf("failed to get user by username: %w", err)
	}

	merchantBalance, err := s.merchantRepository.GetMerchantBalance(request.MerchantID)
	if err != nil {
		return model.Customer{}, fmt.Errorf("failed to get merchant balance: %w", err)
	}

	if customer.Balance < request.Amount {
		return model.Customer{}, fmt.Errorf("insufficient balance")
	}

	merchantBalance += request.Amount
	customer.Balance -= request.Amount

	updatedCustomer, err := s.customerRepository.UpdateUserBalance(customer.ID, customer.Balance)
	if err != nil {
		return model.Customer{}, fmt.Errorf("failed to update user balance: %w", err)
	}

	_, err = s.merchantRepository.UpdateMerchantBalance(request.MerchantID, merchantBalance)
	if err != nil {
		return model.Customer{}, fmt.Errorf("failed to update merchant balance: %w", err)
	}

	return updatedCustomer, nil
}
