package main

import (
	"context"
	"github.com/qwerun/habr-notification-go/internal/notify"
	"github.com/qwerun/habr-notification-go/pkg/kafka"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	initialBackoff = 1 * time.Second
	maxBackoff     = 30 * time.Second
)

func main() {
	kc, err := kafka.NewKafkaConsumerGroup()
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %v", err)
	}
	defer func() {
		if err := kc.Close(); err != nil {
			log.Printf("Ошибка при закрытии Kafka consumer: %v", err)
		}
	}()
	topicsEnv := os.Getenv("KAFKA_TOPIC")
	if topicsEnv == "" {
		log.Fatalf("KAFKA_TOPIC not set")
	}
	topics := strings.Split(topicsEnv, ",")

	handler := &notify.EmailHandler{}

	kafkaHandler := &kafka.ConsumerHandler{Handler: handler}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigchan
		cancel()
		log.Println("Start graceful shutdown")
		signal.Stop(sigchan)
	}()

	backoff := initialBackoff

	for {
		err := kc.Consume(ctx, topics, kafkaHandler)
		if ctx.Err() != nil {
			break
		}

		if err != nil {
			log.Printf("kafka consume error: %v", err)
			log.Printf("Waiting %s before next consume attempt", backoff)
			time.Sleep(backoff)

			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		backoff = initialBackoff
	}
}
