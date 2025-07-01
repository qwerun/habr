package kafka

import (
	"github.com/IBM/sarama"
	"os"
)

type Consumer struct {
	Group  sarama.ConsumerGroup
	Topics []string
}

func NewConsumer() (*Consumer, error) {
	brokers := []string{os.Getenv("KAFKA_BROKER")}
	group := os.Getenv("KAFKA_GROUP")
	topics := []string{os.Getenv("KAFKA_TOPIC")}

	config := sarama.NewConfig()
	config.Version = sarama.V3_7_2_0

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		Group:  consumerGroup,
		Topics: topics,
	}, nil
}
