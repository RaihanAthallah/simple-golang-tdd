package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function untuk inisialisasi repository dengan file data json
func setupRepository(t *testing.T) CustomerRepository {
	repo, err := NewCustomerRepository("../../data/customers.json")
	require.NoError(t, err)
	return repo
}

// ========== SUCCESS CASES ==========

func TestGetUserByUsername_Success(t *testing.T) {
	repo := setupRepository(t)

	customer, err := repo.GetUserByUsername("johndoe")

	require.NoError(t, err)
	assert.Equal(t, "johndoe", customer.Username)
}

func TestGetUserByID_Success(t *testing.T) {
	repo := setupRepository(t)

	customer, err := repo.GetUserByID("cust-001")

	require.NoError(t, err)
	assert.Equal(t, "cust-001", customer.ID)
}

func TestGetUserBalance_Success(t *testing.T) {
	repo := setupRepository(t)

	balance, err := repo.GetUserBalance("cust-001")

	require.NoError(t, err)
	assert.Equal(t, balance, balance)
}

func TestUpdateUserBalance_Success(t *testing.T) {
	repo := setupRepository(t)

	updatedCustomer, err := repo.UpdateUserBalance("cust-001", 500.0)
	balance, err := repo.GetUserBalance("cust-001")

	require.NoError(t, err)
	assert.Equal(t, balance, updatedCustomer.Balance)
}

// ========== ERROR CASES ==========

func TestGetUserByUsername_Error(t *testing.T) {
	repo := setupRepository(t)

	customer, err := repo.GetUserByUsername("unknownuser")

	require.Error(t, err)
	assert.Empty(t, customer)
	assert.EqualError(t, err, "user not found")
}

func TestGetUserByID_Error(t *testing.T) {
	repo := setupRepository(t)

	customer, err := repo.GetUserByID("unknown_id")

	require.Error(t, err)
	assert.Empty(t, customer)
	assert.EqualError(t, err, "user not found by ID")
}

func TestGetUserBalance_Error(t *testing.T) {
	repo := setupRepository(t)

	balance, err := repo.GetUserBalance("unknown_id")

	require.Error(t, err)
	assert.Equal(t, 0.0, balance)
	assert.EqualError(t, err, "user not found for balance check")
}

func TestUpdateUserBalance_Error(t *testing.T) {
	repo := setupRepository(t)

	customer, err := repo.UpdateUserBalance("invalid_id", 500.0)

	require.Error(t, err)
	assert.Empty(t, customer)
	assert.EqualError(t, err, "error while updating user balance")
}
