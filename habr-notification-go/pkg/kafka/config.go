package kafka

import (
	"errors"
	"github.com/IBM/sarama"
	"os"
	"strings"
	"time"
)

type Settings struct {
	Brokers             []string
	GroupID             string
	Version             sarama.KafkaVersion
	InitialOffset       int64
	RebalanceStrategies []sarama.BalanceStrategy
	RebalanceTimeout    time.Duration
	SessionTimeout      time.Duration
	HeartbeatInterval   time.Duration
	AutoCommit          bool
	AutoCommitInterval  time.Duration
	ReturnErrors        bool
}

func LoadKafkaSettings() (*Settings, error) {
	brokersEnv := os.Getenv("KAFKA_BROKER")
	if brokersEnv == "" {
		return nil, errors.New("KAFKA_BROKER not set")
	}
	groupEnv := os.Getenv("KAFKA_GROUP")
	if groupEnv == "" {
		return nil, errors.New("KAFKA_GROUP not set")
	}

	return &Settings{
		Brokers:             strings.Split(brokersEnv, ","),
		GroupID:             groupEnv,
		Version:             sarama.V3_7_2_0,
		InitialOffset:       sarama.OffsetOldest,
		RebalanceStrategies: []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()},
		RebalanceTimeout:    60 * time.Second,
		SessionTimeout:      18 * time.Second,
		HeartbeatInterval:   8 * time.Second,
		AutoCommit:          true,
		AutoCommitInterval:  1 * time.Second,
		ReturnErrors:        true,
	}, nil
}

func NewKafkaConfig(s *Settings) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = s.Version

	config.Consumer.Offsets.Initial = s.InitialOffset
	config.Consumer.Group.Rebalance.GroupStrategies = s.RebalanceStrategies
	config.Consumer.Group.Rebalance.Timeout = s.RebalanceTimeout
	config.Consumer.Group.Session.Timeout = s.SessionTimeout
	config.Consumer.Group.Heartbeat.Interval = s.HeartbeatInterval

	config.Consumer.Return.Errors = s.ReturnErrors
	config.Consumer.Offsets.AutoCommit.Enable = s.AutoCommit
	config.Consumer.Offsets.AutoCommit.Interval = s.AutoCommitInterval

	return config
}
