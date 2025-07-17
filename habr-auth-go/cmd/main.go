package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/qwerun/habr-auth-go/internal/auth"
	"github.com/qwerun/habr-auth-go/internal/handlers"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
	"github.com/qwerun/habr-auth-go/pkg/config"
	"github.com/qwerun/habr-auth-go/pkg/kafka"
	"github.com/qwerun/habr-auth-go/pkg/postgres"
	"github.com/qwerun/habr-auth-go/pkg/redis"
)

func main() {
	db, err := config.NewPostgresDB()
	check(err, "postgres")

	rdb, err := config.NewRedisDB()
	check(err, "redis")

	producer, err := config.NewKafkaProducer()
	check(err, "kafka")

	key, at, rt, err := auth.GetJwtInfo()
	check(err, "jwt")
	jwtMan := auth.NewJwtManager(key, at, rt)

	pExplorer := kafka.NewKafkaExplorer(producer, strings.Split(os.Getenv("KAFKA_TOPIC"), ","))
	rExplorer := redis.NewRedisExplorer(rdb)
	explorer := postgres.NewExplorer(db)

	userRepo := user_repository.New(explorer, rExplorer, pExplorer)
	handler, err := handlers.NewMux(userRepo, jwtMan)
	check(err, "mux")

	srv := &http.Server{
		Addr:    ":8081",
		Handler: handler,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Start graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}

	closeQuiet(pExplorer.Producer.Close, "kafka producer")
	closeQuiet(rdb.Close, "redis")
	closeQuiet(db.Close, "postgres")

	log.Println("Graceful shutdown completed")
}

func check(err error, what string) {
	if err != nil {
		log.Fatalf("%s: %v", what, err)
	}
}

func closeQuiet(f func() error, what string) {
	if err := f(); err != nil {
		log.Printf("%s close: %v", what, err)
	}
}
