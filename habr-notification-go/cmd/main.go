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
)

func main() {
	kc, err := kafka.NewKafkaConsumerGroup()
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %v", err)
	}
	defer kc.Close()
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

	for {
		if err := kc.Consume(ctx, topics, kafkaHandler); err != nil {
			log.Printf("kafka consume error: %v", err)
			if ctx.Err() != nil {
				break
			}
		}
	}
}
