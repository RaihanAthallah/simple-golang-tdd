package repository

import (
	"simple-golang-tdd/model"
	"simple-golang-tdd/utils"
	"sync"
)

type HistoryRepository interface {
	UpdateHistory(history model.History) (model.History, error)
}

type historyRepositoryImpl struct {
	dataSourcePath string
	histories      []model.History
	mutex          sync.RWMutex // Untuk menghindari masalah concurrency jika diperlukan
}

// NewHistoryRepository membuat repository baru dan membaca file JSON sekali saja.
func NewHistoryRepository(dataSourcePath string) (HistoryRepository, error) {
	repo := &historyRepositoryImpl{dataSourcePath: dataSourcePath}
	err := repo.loadData()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *historyRepositoryImpl) loadData() error {
	return utils.LoadJSONFile(r.dataSourcePath, &r.histories)
}

func (r *historyRepositoryImpl) saveHistoriesToFile() error {
	return utils.SaveJSONFile(r.dataSourcePath, r.histories)
}

func (r *historyRepositoryImpl) UpdateHistory(history model.History) (model.History, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Append the new history to the in-memory slice
	r.histories = append(r.histories, history)

	// Now, save the updated slice to the JSON file
	err := r.saveHistoriesToFile()
	if err != nil {
		return model.History{}, err
	}

	return history, nil
}
