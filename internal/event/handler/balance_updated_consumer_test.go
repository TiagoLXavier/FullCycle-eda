package handler

import (
	"errors"
	"testing"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/entity"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/mocks"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/update_balance"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBalanceUpdatedConsumerHandler_Handle(t *testing.T) {
	m := &mocks.BalanceGatewayMock{}
	m.On("Save", mock.Anything).Return(nil)
	h := NewBalanceUpdatedConsumerHandler(*update_balance.NewUpdateBalanceUseCase(m))

	msg := &ckafka.Message{Value: []byte(`{"Name":"BalanceUpdated","Payload":{"account_id_from":"a1","account_id_to":"a2","balance_account_id_from":900,"balance_account_id_to":1100}}`)}
	err := h.Handle(msg)
	assert.Nil(t, err)
	m.AssertNumberOfCalls(t, "Save", 2)
	m.AssertCalled(t, "Save", mock.MatchedBy(func(b *entity.Balance) bool {
		return b.AccountID == "a1" && b.Balance == 900
	}))
	m.AssertCalled(t, "Save", mock.MatchedBy(func(b *entity.Balance) bool {
		return b.AccountID == "a2" && b.Balance == 1100
	}))
}

func TestBalanceUpdatedConsumerHandler_HandleInvalidJSON(t *testing.T) {
	m := &mocks.BalanceGatewayMock{}
	h := NewBalanceUpdatedConsumerHandler(*update_balance.NewUpdateBalanceUseCase(m))

	err := h.Handle(&ckafka.Message{Value: []byte(`not json`)})
	assert.NotNil(t, err)
	m.AssertNumberOfCalls(t, "Save", 0)
}

func TestBalanceUpdatedConsumerHandler_HandleWithoutAccountIDs(t *testing.T) {
	m := &mocks.BalanceGatewayMock{}
	h := NewBalanceUpdatedConsumerHandler(*update_balance.NewUpdateBalanceUseCase(m))

	err := h.Handle(&ckafka.Message{Value: []byte(`{"Name":"BalanceUpdated","Payload":{"account_id_from":"a1"}}`)})
	assert.NotNil(t, err)
	m.AssertNumberOfCalls(t, "Save", 0)
}

func TestBalanceUpdatedConsumerHandler_HandleWhenUseCaseFails(t *testing.T) {
	m := &mocks.BalanceGatewayMock{}
	m.On("Save", mock.Anything).Return(errors.New("save failed"))
	h := NewBalanceUpdatedConsumerHandler(*update_balance.NewUpdateBalanceUseCase(m))

	err := h.Handle(&ckafka.Message{Value: []byte(`{"Name":"BalanceUpdated","Payload":{"account_id_from":"a1","account_id_to":"a2"}}`)})
	assert.EqualError(t, err, "save failed")
}
