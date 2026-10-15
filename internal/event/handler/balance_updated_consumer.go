package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/update_balance"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type BalanceUpdatedMessage struct {
	Name    string                               `json:"Name"`
	Payload update_balance.UpdateBalanceInputDTO `json:"Payload"`
}

type BalanceUpdatedConsumerHandler struct {
	UpdateBalanceUseCase update_balance.UpdateBalanceUseCase
}

func NewBalanceUpdatedConsumerHandler(updateBalanceUseCase update_balance.UpdateBalanceUseCase) *BalanceUpdatedConsumerHandler {
	return &BalanceUpdatedConsumerHandler{
		UpdateBalanceUseCase: updateBalanceUseCase,
	}
}

func (h *BalanceUpdatedConsumerHandler) Handle(msg *ckafka.Message) error {
	var message BalanceUpdatedMessage
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}
	payload := message.Payload
	if payload.AccountIDFrom == "" || payload.AccountIDTo == "" {
		return fmt.Errorf("message without account ids: %s", string(msg.Value))
	}
	if err := h.UpdateBalanceUseCase.Execute(payload); err != nil {
		return err
	}
	log.Printf("balances updated: %s=%v %s=%v", payload.AccountIDFrom, payload.BalanceAccountIDFrom, payload.AccountIDTo, payload.BalanceAccountIDTo)
	return nil
}
