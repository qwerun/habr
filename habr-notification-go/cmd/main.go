package main

import (
	"github.com/qwerun/habr-notification-go/internal"
	"github.com/qwerun/habr-notification-go/pkg/kafka"
	"log"
)

func main() {
	kc, err := kafka.NewConsumer()
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %v", err)
	}
	defer kc.Group.Close()

	ic := internal.NewConsumer(kc)
	err = ic.Notify()
	if err != nil {
		log.Fatalf("failed to use kafka consumer: %v", err)
	}
}
