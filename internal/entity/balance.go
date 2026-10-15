package entity

import (
	"errors"
	"time"
)

type Balance struct {
	AccountID string
	Balance   float64
	UpdatedAt time.Time
}

func NewBalance(accountID string, balance float64) (*Balance, error) {
	b := &Balance{
		AccountID: accountID,
		Balance:   balance,
		UpdatedAt: time.Now(),
	}
	err := b.Validate()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Balance) Validate() error {
	if b.AccountID == "" {
		return errors.New("account id is required")
	}
	return nil
}
