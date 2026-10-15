package update_balance

import (
	"errors"
	"testing"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/entity"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateBalanceUseCase_Execute(t *testing.T) {
	m := &mocks.BalanceGatewayMock{}
	m.On("Save", mock.Anything).Return(nil)
	uc := NewUpdateBalanceUseCase(m)

	err := uc.Execute(UpdateBalanceInputDTO{
		AccountIDFrom:        "account-1",
		AccountIDTo:          "account-2",
		BalanceAccountIDFrom: 900,
		BalanceAccountIDTo:   1100,
	})
	assert.Nil(t, err)
	m.AssertExpectations(t)
	m.AssertNumberOfCalls(t, "Save", 2)
	m.AssertCalled(t, "Save", mock.MatchedBy(func(b *entity.Balance) bool {
		return b.AccountID == "account-1" && b.Balance == 900
	}))
	m.AssertCalled(t, "Save", mock.MatchedBy(func(b *entity.Balance) bool {
		return b.AccountID == "account-2" && b.Balance == 1100
	}))
}

func TestUpdateBalanceUseCase_ExecuteWithInvalidAccountID(t *testing.T) {
	m := &mocks.BalanceGatewayMock{}
	uc := NewUpdateBalanceUseCase(m)

	err := uc.Execute(UpdateBalanceInputDTO{
		AccountIDFrom: "account-1",
		AccountIDTo:   "",
	})
	assert.NotNil(t, err)
	m.AssertNumberOfCalls(t, "Save", 0)
}

func TestUpdateBalanceUseCase_ExecuteWhenSaveFails(t *testing.T) {
	m := &mocks.BalanceGatewayMock{}
	m.On("Save", mock.Anything).Return(errors.New("save failed"))
	uc := NewUpdateBalanceUseCase(m)

	err := uc.Execute(UpdateBalanceInputDTO{
		AccountIDFrom: "account-1",
		AccountIDTo:   "account-2",
	})
	assert.EqualError(t, err, "save failed")
	m.AssertNumberOfCalls(t, "Save", 1)
}
