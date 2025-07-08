package main

import (
	"github.com/IBM/sarama"
	"github.com/qwerun/habr-auth-go/internal/handlers"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
	"github.com/qwerun/habr-auth-go/pkg/kafka"
	"github.com/qwerun/habr-auth-go/pkg/postgres"
	"github.com/qwerun/habr-auth-go/pkg/redis"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	db, err := postgres.NewPostgresDB()
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	rdb, err := redis.NewRedisDB()
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	pc, err := kafka.NewKafkaProducer()
	if err != nil {
		log.Fatalf("kafka: %v", err)
	}

	pExplorer := kafka.NewKafkaExplorer(pc, strings.Split(os.Getenv("KAFKA_TOPIC"), ","))
	defer func(Producer sarama.SyncProducer) {
		err = Producer.Close()
		if err != nil {
			log.Fatalf("kafka close producer: %v", err)
		}
	}(pExplorer.Producer)
	rExplorer := redis.NewRedisExplorer(rdb)
	explorer := postgres.NewExplorer(db)
	userRepo := user_repository.New(explorer, rExplorer, pExplorer)
	handler, err := handlers.NewMux(userRepo)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8081", handler)
	if err != nil {
		log.Fatalf("start server: %v", err)
	}

}
