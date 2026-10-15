package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/database"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/event/handler"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/get_balance"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/update_balance"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/web"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/web/webserver"
	"github.com.br/devfullcycle/fc-ms-wallet/pkg/kafka"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", "root"),
		getEnv("DB_HOST", "mysql-balances"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "balances"),
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for {
		if err := db.Ping(); err == nil {
			break
		} else {
			log.Printf("waiting for database: %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	balanceDb := database.NewBalanceDB(db)

	updateBalanceUseCase := update_balance.NewUpdateBalanceUseCase(balanceDb)
	getBalanceUseCase := get_balance.NewGetBalanceUseCase(balanceDb)

	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:29092"),
		"group.id":          "balances",
		"auto.offset.reset": "earliest",
	}
	msgChan := make(chan *ckafka.Message)
	consumer := kafka.NewConsumer(configMap, []string{"balances"})
	go func() {
		if err := consumer.Consume(msgChan); err != nil {
			log.Printf("consumer stopped: %v", err)
		}
	}()

	balanceUpdatedConsumerHandler := handler.NewBalanceUpdatedConsumerHandler(*updateBalanceUseCase)
	go func() {
		for msg := range msgChan {
			if err := balanceUpdatedConsumerHandler.Handle(msg); err != nil {
				log.Printf("error handling message: %v", err)
			}
		}
	}()

	balanceHandler := web.NewWebBalanceHandler(*getBalanceUseCase)

	webserver := webserver.NewWebServer(":3003")
	webserver.AddGetHandler("/balances/{account_id}", balanceHandler.GetBalance)

	log.Println("Balances is running on :3003")
	webserver.Start()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
