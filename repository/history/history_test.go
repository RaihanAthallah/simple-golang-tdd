package repository

import (
	"simple-golang-tdd/dto"
	"simple-golang-tdd/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function untuk inisialisasi repository dengan file data json
func setupRepository(t *testing.T) HistoryRepository {
	repo, err := NewHistoryRepository("../../data/histories.json")
	require.NoError(t, err)
	return repo
}

func TestUpdateHistory_Success(t *testing.T) {
	repo := setupRepository(t)

	fakeRequest := dto.PaymentRequest{
		MerchantID: "merchant-002",
		Amount:     500.0,
	}

	fakeModel := model.History{
		ID:         "invalid_id",
		CustomerID: "cust-001",
		Action:     "update",
		Details:    fakeRequest,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	history, err := repo.UpdateHistory(fakeModel)

	require.NoError(t, err)
	assert.Equal(t, history, fakeModel)
}

func TestUpdateHistory_Error(t *testing.T) {
	repo := setupRepository(t)

	fakeRequest := dto.PaymentRequest{
		MerchantID: "merchant-002",
		Amount:     500.0,
	}

	fakeModel := model.History{
		ID:         "invalid_id",
		CustomerID: "cust-001",
		Action:     "update",
		Details:    fakeRequest,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	customer, err := repo.UpdateHistory(fakeModel)

	require.Error(t, err)
	assert.Empty(t, customer)
	assert.EqualError(t, err, "error while updating logging history")
}

func TestUpdateHistory_ErrorSavingFile(t *testing.T) {
	repo := setupRepository(t)

	fakeRequest := dto.PaymentRequest{
		MerchantID: "merchant-002",
		Amount:     500.0,
	}

	fakeModel := model.History{
		ID:         "invalid_id",
		CustomerID: "cust-001",
		Action:     "update",
		Details:    fakeRequest,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	customer, err := repo.UpdateHistory(fakeModel)

	require.Error(t, err)
	assert.Empty(t, customer)
	assert.EqualError(t, err, "error while saving history to file")
}
