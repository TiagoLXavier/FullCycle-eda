package balances

import (
	"encoding/json"
	"fmt"
	"log"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type BalanceUpdatedEvent struct {
	Name    string `json:"Name"`
	Payload struct {
		AccountIDFrom        string  `json:"account_id_from"`
		AccountIDTo          string  `json:"account_id_to"`
		BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
		BalanceAccountIDTo   float64 `json:"balance_account_id_to"`
	} `json:"Payload"`
}

func (r *Repository) HandleMessage(msg *ckafka.Message) error {
	var event BalanceUpdatedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}
	if event.Payload.AccountIDFrom == "" || event.Payload.AccountIDTo == "" {
		return fmt.Errorf("message without account ids: %s", string(msg.Value))
	}
	if err := r.Upsert(event.Payload.AccountIDFrom, event.Payload.BalanceAccountIDFrom); err != nil {
		return err
	}
	if err := r.Upsert(event.Payload.AccountIDTo, event.Payload.BalanceAccountIDTo); err != nil {
		return err
	}
	log.Printf("balances updated: %s=%v %s=%v", event.Payload.AccountIDFrom, event.Payload.BalanceAccountIDFrom, event.Payload.AccountIDTo, event.Payload.BalanceAccountIDTo)
	return nil
}
