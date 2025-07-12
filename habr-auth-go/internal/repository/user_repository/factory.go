package user_repository

import (
	"github.com/qwerun/habr-auth-go/pkg/jwtManager"
	"github.com/qwerun/habr-auth-go/pkg/kafka"
	"github.com/qwerun/habr-auth-go/pkg/postgres"
	"github.com/qwerun/habr-auth-go/pkg/redis"
)

type Repository struct {
	explorer      *postgres.Explorer
	redisExplorer *redis.RedisExplorer
	kafkaExplorer *kafka.KafkaExplorer
	jwtExplorer   *jwtManager.JwtExplorer
}

func New(explorer *postgres.Explorer, redis *redis.RedisExplorer, kafka *kafka.KafkaExplorer, jwtManager *jwtManager.JwtExplorer) *Repository {
	return &Repository{explorer: explorer, redisExplorer: redis, kafkaExplorer: kafka, jwtExplorer: jwtManager}
}
