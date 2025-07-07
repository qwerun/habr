package kafka

import (
	"errors"
	"github.com/IBM/sarama"
	"os"
	"strings"
	"time"
)

func NewKafkaConsumerGroup() (sarama.ConsumerGroup, error) {
	brokersEnv := os.Getenv("KAFKA_BROKER")
	if brokersEnv == "" {
		return nil, errors.New("KAFKA_BROKER not set")
	}
	brokers := strings.Split(brokersEnv, ",")

	group := os.Getenv("KAFKA_GROUP")
	if group == "" {
		return nil, errors.New("KAFKA_GROUP not set")
	}

	config := newKafkaConfig()

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}
	return consumerGroup, nil
}

func newKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()

	config.Version = sarama.V3_7_2_0

	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second

	config.Consumer.Group.Session.Timeout = 18 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 8 * time.Second

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	return config
}
