package gateway

import "github.com.br/devfullcycle/fc-ms-wallet/internal/entity"

type BalanceGateway interface {
	Save(balance *entity.Balance) error
	FindByAccountID(accountID string) (*entity.Balance, error)
}
