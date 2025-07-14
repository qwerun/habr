package config

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/IBM/sarama"
	rds "github.com/redis/go-redis/v9"
	"os"
	"strings"
)

func NewRedisDB() (*rds.Client, error) {
	dbPassword := os.Getenv("REDIS_PASSWORD")
	dbAddr := os.Getenv("REDIS_ADDR")

	rdb := rds.NewClient(&rds.Options{
		Addr:     dbAddr,
		Password: dbPassword,
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func NewKafkaProducer() (sarama.SyncProducer, error) {
	brokers := strings.Split(os.Getenv("KAFKA_BROKER"), ",")
	config := sarama.NewConfig()
	config.Version = sarama.V3_7_2_0
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	return sarama.NewSyncProducer(brokers, config)
}

func NewPostgresDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
