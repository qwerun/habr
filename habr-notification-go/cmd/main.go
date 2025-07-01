package main

import (
	"github.com/qwerun/habr-notification-go/pkg/kafka"
	"log"
)

func main() {
	consumer, err := kafka.NewConsumer()
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %v", err)
	}
}
