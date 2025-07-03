package kafka

import (
	"github.com/IBM/sarama"
	"log"
	"os"
	"strings"
)

func NewKafkaConsumerGroup() (sarama.ConsumerGroup, error) {
	brokersEnv := os.Getenv("KAFKA_BROKER")
	if brokersEnv == "" {
		log.Fatalf("KAFKA_BROKER not set")
	}
	brokers := strings.Split(brokersEnv, ",")

	group := os.Getenv("KAFKA_GROUP")
	if group == "" {
		log.Fatalf("KAFKA_GROUP not set")
	}

	config := sarama.NewConfig()
	config.Version = sarama.V3_7_2_0
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}
	return consumerGroup, nil
}
