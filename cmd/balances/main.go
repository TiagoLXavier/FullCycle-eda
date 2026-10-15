package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/balances"
	"github.com.br/devfullcycle/fc-ms-wallet/pkg/kafka"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"
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

	repository := balances.NewRepository(db)

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
	go func() {
		for msg := range msgChan {
			if err := repository.HandleMessage(msg); err != nil {
				log.Printf("error handling message: %v", err)
			}
		}
	}()

	handler := balances.NewHandler(repository)
	router := chi.NewRouter()
	router.Get("/balances/{account_id}", handler.GetBalance)

	log.Println("Balances is running on :3003")
	log.Fatal(http.ListenAndServe(":3003", router))
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
