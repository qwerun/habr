package internal

import (
	"github.com/IBM/sarama"
	"github.com/qwerun/habr-notification-go/pkg/kafka"
)

type Consumer struct {
	*kafka.Consumer
}

func NewConsumer(consumer *kafka.Consumer) *Consumer {
	return &Consumer{consumer}
}

func (c *Consumer) Notify() error {

	return nil
}

func ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error {
	return nil
}

func Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
