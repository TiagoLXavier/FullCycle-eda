package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewBalance(t *testing.T) {
	balance, err := NewBalance("account-1", 100.5)
	assert.Nil(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, "account-1", balance.AccountID)
	assert.Equal(t, 100.5, balance.Balance)
	assert.False(t, balance.UpdatedAt.IsZero())
}

func TestCreateNewBalanceWhenArgsAreInvalid(t *testing.T) {
	balance, err := NewBalance("", 100)
	assert.NotNil(t, err)
	assert.Nil(t, balance)
	assert.EqualError(t, err, "account id is required")
}
