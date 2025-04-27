package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function untuk inisialisasi repository dengan file data json
func setupRepository(t *testing.T) MerchantRepository {
	repo, err := NewMerchantRepository("../../data/merchants.json")
	require.NoError(t, err)
	return repo
}

func TestGetMerchantBalance_Success(t *testing.T) {
	repo := setupRepository(t)

	balance, err := repo.GetMerchantBalance("merchant-001")

	require.NoError(t, err)
	assert.Equal(t, balance, balance)
}

func TestUpdateMerchantBalance_Success(t *testing.T) {
	repo := setupRepository(t)

	updatedCustomer, err := repo.UpdateMerchantBalance("merchant-001", 500.0)
	balance, err := repo.GetMerchantBalance("merchant-001")

	require.NoError(t, err)
	assert.Equal(t, balance, updatedCustomer.Balance)
}

func TestGetMerchantBalance_Error(t *testing.T) {
	repo := setupRepository(t)

	balance, err := repo.GetMerchantBalance("unknown_id")

	require.Error(t, err)
	assert.Equal(t, 0.0, balance)
	assert.EqualError(t, err, "merchant not found for balance check")
}

func TestUpdateMerchantBalance_Error(t *testing.T) {
	repo := setupRepository(t)

	customer, err := repo.UpdateMerchantBalance("invalid_id", 500.0)

	require.Error(t, err)
	assert.Empty(t, customer)
	assert.EqualError(t, err, "error while updating merchant balance")
}
