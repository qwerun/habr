package main

import (
	"context"
	"github.com/qwerun/habr-notification-go/internal"
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
	kExplorer := kafka.NewKafkaConsumerExplorer(kc, strings.Split(os.Getenv("KAFKA_TOPIC"), ","))
	defer kExplorer.ConsumerGroup.Close()

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

	ic := internal.NewConsumer(kExplorer)
	err = ic.Notify(ctx)
	if err != nil {
		log.Fatalf("failed to use kafka consumer: %v", err)
	}
}
