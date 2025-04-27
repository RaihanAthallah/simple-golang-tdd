package repository

import (
	"errors"
	"fmt"
	"simple-golang-tdd/model"
	"simple-golang-tdd/utils"
	"sync"
)

type CustomerRepository interface {
	GetUserByUsername(username string) (model.Customer, error)
	GetUserByID(id string) (model.Customer, error)
	GetUserBalance(id string) (float64, error)
	UpdateUserBalance(id string, amount float64) (model.Customer, error)
}

type customerRepositoryImpl struct {
	dataSourcePath string
	customers      []model.Customer
	mutex          sync.RWMutex // Untuk menghindari masalah concurrency jika diperlukan
}

// NewCustomerRepository membuat repository baru dan membaca file JSON sekali saja.
func NewCustomerRepository(dataSourcePath string) (CustomerRepository, error) {
	repo := &customerRepositoryImpl{dataSourcePath: dataSourcePath}
	err := repo.loadData()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *customerRepositoryImpl) loadData() error {
	return utils.LoadJSONFile(r.dataSourcePath, &r.customers)
}

func (r *customerRepositoryImpl) saveCustomersToFile() error {
	return utils.SaveJSONFile(r.dataSourcePath, r.customers)
}

func (r *customerRepositoryImpl) GetUserByUsername(username string) (model.Customer, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.customers {
		if user.Username == username {
			return user, nil
		}
	}
	return model.Customer{}, errors.New("user not found")
}

func (r *customerRepositoryImpl) GetUserByID(id string) (model.Customer, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.customers {
		if user.ID == id {
			return user, nil
		}
	}
	return model.Customer{}, errors.New("user not found by ID")
}

func (r *customerRepositoryImpl) GetUserBalance(id string) (float64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.customers {
		if user.ID == id {
			return user.Balance, nil
		}
	}
	return 0.0, errors.New("user not found for balance check")
}

func (r *customerRepositoryImpl) UpdateUserBalance(id string, amount float64) (model.Customer, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, user := range r.customers {
		if user.ID == id {
			r.customers[i].Balance = amount
			err := r.saveCustomersToFile()
			if err != nil {
				return model.Customer{}, fmt.Errorf("error while updating user balance: %v", err)
			}
			return r.customers[i], nil
		}
	}

	return model.Customer{}, errors.New("error while updating user balance")
}
