package repository

import (
	"errors"
	"simple-golang-tdd/model"
	"simple-golang-tdd/utils"
	"sync"
)

type MerchantRepository interface {
	UpdateMerchantBalance(id string, amount float64) (model.Merchant, error)
	GetMerchantBalance(id string) (float64, error)
}

type merchantRepositoryImpl struct {
	dataSourcePath string
	merchants      []model.Merchant
	mutex          sync.RWMutex // Untuk menghindari masalah concurrency jika diperlukan
}

// NewMerchantRepository membuat repository baru dan membaca file JSON sekali saja.
func NewMerchantRepository(dataSourcePath string) (MerchantRepository, error) {
	repo := &merchantRepositoryImpl{dataSourcePath: dataSourcePath}
	err := repo.loadData()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *merchantRepositoryImpl) loadData() error {
	return utils.LoadJSONFile(r.dataSourcePath, &r.merchants)
}

func (r *merchantRepositoryImpl) saveMerchantsToFile() error {
	return utils.SaveJSONFile(r.dataSourcePath, r.merchants)
}

func (r *merchantRepositoryImpl) UpdateMerchantBalance(id string, amount float64) (model.Merchant, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, user := range r.merchants {
		if user.ID == id {
			r.merchants[i].Balance = amount
			err := r.saveMerchantsToFile()
			if err != nil {
				return model.Merchant{}, err
			}
			return r.merchants[i], nil
		}
	}

	return model.Merchant{}, errors.New("error while updating merchant balance")
}

func (r *merchantRepositoryImpl) GetMerchantBalance(id string) (float64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.merchants {
		if user.ID == id {
			return user.Balance, nil
		}
	}

	return 0, errors.New("merchant not found for balance check")
}
