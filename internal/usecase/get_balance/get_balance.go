package get_balance

import (
	"github.com.br/devfullcycle/fc-ms-wallet/internal/gateway"
)

type GetBalanceInputDTO struct {
	AccountID string
}

type GetBalanceOutputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type GetBalanceUseCase struct {
	BalanceGateway gateway.BalanceGateway
}

func NewGetBalanceUseCase(balanceGateway gateway.BalanceGateway) *GetBalanceUseCase {
	return &GetBalanceUseCase{
		BalanceGateway: balanceGateway,
	}
}

func (uc *GetBalanceUseCase) Execute(input GetBalanceInputDTO) (*GetBalanceOutputDTO, error) {
	balance, err := uc.BalanceGateway.FindByAccountID(input.AccountID)
	if err != nil {
		return nil, err
	}
	output := &GetBalanceOutputDTO{
		AccountID: balance.AccountID,
		Balance:   balance.Balance,
	}
	return output, nil
}
