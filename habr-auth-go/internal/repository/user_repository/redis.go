package user_repository

import (
	"context"
	"log"
	"math/rand"
	"time"
)

func (r *Repository) SetVerificationCode(email string) (int, error) {
	ctx := context.Background()
	code := rand.Intn(900000) + 100000
	err := r.redisExplorer.RDB.Set(ctx, email, code, 10*time.Minute).Err()
	if err != nil {
		log.Printf("SetVerificationCode: Could not set value in Redis: %v", err)
		return 0, err
	}
	return code, nil
}

func (r *Repository) CheckVerificationCode(email, code string) error {
	ctx := context.Background()
	validCode, err := r.redisExplorer.RDB.Get(ctx, email).Result()
	if err != nil {
		log.Printf("CheckVerificationCode: Could not get value in Redis: %v", err)
		return err
	}
	if validCode != code {
		log.Printf("CheckVerificationCode: %s: email: %v checkCode: %v validCode: %v", ErrCodeCheckNotFound, email, code, validCode)
		return ErrCodeCheckNotFound
	}
	return nil
}
