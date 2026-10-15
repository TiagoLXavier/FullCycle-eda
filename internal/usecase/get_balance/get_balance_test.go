package get_balance

import (
	"database/sql"
	"testing"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/entity"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetBalanceUseCase_Execute(t *testing.T) {
	balance, _ := entity.NewBalance("account-1", 150.75)
	m := &mocks.BalanceGatewayMock{}
	m.On("FindByAccountID", "account-1").Return(balance, nil)
	uc := NewGetBalanceUseCase(m)

	output, err := uc.Execute(GetBalanceInputDTO{AccountID: "account-1"})
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "account-1", output.AccountID)
	assert.Equal(t, 150.75, output.Balance)
	m.AssertExpectations(t)
	m.AssertNumberOfCalls(t, "FindByAccountID", 1)
}

func TestGetBalanceUseCase_ExecuteWhenNotFound(t *testing.T) {
	m := &mocks.BalanceGatewayMock{}
	m.On("FindByAccountID", "missing").Return(nil, sql.ErrNoRows)
	uc := NewGetBalanceUseCase(m)

	output, err := uc.Execute(GetBalanceInputDTO{AccountID: "missing"})
	assert.Nil(t, output)
	assert.Equal(t, sql.ErrNoRows, err)
	m.AssertNumberOfCalls(t, "FindByAccountID", 1)
}
