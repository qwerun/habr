package kafka

import (
	"github.com/IBM/sarama"
	"os"
	"strings"
)

type KafkaExplorer struct {
	Producer sarama.SyncProducer
	Topic    []string
}

func NewKafkaExplorer(producer sarama.SyncProducer, topic []string) *KafkaExplorer {
	return &KafkaExplorer{
		Producer: producer,
		Topic:    topic,
	}
}

func NewKafkaProducer() (sarama.SyncProducer, error) {
	brokers := strings.Split(os.Getenv("KAFKA_BROKER"), ",")
	config := sarama.NewConfig()
	config.Version = sarama.V3_7_2_0
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	return sarama.NewSyncProducer(brokers, config)
}
