package update_balance

import (
	"github.com.br/devfullcycle/fc-ms-wallet/internal/entity"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/gateway"
)

type UpdateBalanceInputDTO struct {
	AccountIDFrom        string  `json:"account_id_from"`
	AccountIDTo          string  `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIDTo   float64 `json:"balance_account_id_to"`
}

type UpdateBalanceUseCase struct {
	BalanceGateway gateway.BalanceGateway
}

func NewUpdateBalanceUseCase(balanceGateway gateway.BalanceGateway) *UpdateBalanceUseCase {
	return &UpdateBalanceUseCase{
		BalanceGateway: balanceGateway,
	}
}

func (uc *UpdateBalanceUseCase) Execute(input UpdateBalanceInputDTO) error {
	from, err := entity.NewBalance(input.AccountIDFrom, input.BalanceAccountIDFrom)
	if err != nil {
		return err
	}
	to, err := entity.NewBalance(input.AccountIDTo, input.BalanceAccountIDTo)
	if err != nil {
		return err
	}
	err = uc.BalanceGateway.Save(from)
	if err != nil {
		return err
	}
	err = uc.BalanceGateway.Save(to)
	if err != nil {
		return err
	}
	return nil
}
