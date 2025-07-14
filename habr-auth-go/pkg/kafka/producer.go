package kafka

import (
	"github.com/IBM/sarama"
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
