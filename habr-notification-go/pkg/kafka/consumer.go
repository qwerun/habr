package kafka

import (
	"github.com/IBM/sarama"
	"os"
	"strings"
)

type KafkaConsumerExplorer struct {
	ConsumerGroup sarama.ConsumerGroup
	Topics        []string
}

func NewKafkaConsumerExplorer(consumerGroup sarama.ConsumerGroup, topics []string) *KafkaConsumerExplorer {
	return &KafkaConsumerExplorer{
		ConsumerGroup: consumerGroup,
		Topics:        topics,
	}
}

func NewKafkaConsumerGroup() (sarama.ConsumerGroup, error) {
	brokers := strings.Split(os.Getenv("KAFKA_BROKER"), ",")
	group := os.Getenv("KAFKA_GROUP")

	config := sarama.NewConfig()
	config.Version = sarama.V3_7_2_0
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}
	return consumerGroup, nil
}
