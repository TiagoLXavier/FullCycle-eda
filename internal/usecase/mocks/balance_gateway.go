package mocks

import (
	"github.com.br/devfullcycle/fc-ms-wallet/internal/entity"
	"github.com/stretchr/testify/mock"
)

type BalanceGatewayMock struct {
	mock.Mock
}

func (m *BalanceGatewayMock) Save(balance *entity.Balance) error {
	args := m.Called(balance)
	return args.Error(0)
}

func (m *BalanceGatewayMock) FindByAccountID(accountID string) (*entity.Balance, error) {
	args := m.Called(accountID)
	balance, _ := args.Get(0).(*entity.Balance)
	return balance, args.Error(1)
}
