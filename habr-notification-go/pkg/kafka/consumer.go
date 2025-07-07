package kafka

import (
	"github.com/IBM/sarama"
)

func NewKafkaConsumerGroup() (sarama.ConsumerGroup, error) {

	settings, err := LoadKafkaSettings()
	if err != nil {
		return nil, err
	}
	cfg := NewKafkaConfig(settings)

	consumerGroup, err := sarama.NewConsumerGroup(settings.Brokers, settings.GroupID, cfg)
	if err != nil {
		return nil, err
	}
	return consumerGroup, nil
}
