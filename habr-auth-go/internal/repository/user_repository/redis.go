package user_repository

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func (r *Repository) SetVerificationCode(ctx context.Context, email string) (int, error) {
	code := rand.Intn(900000) + 100000
	err := r.redisExplorer.RDB.Set(ctx, email, code, 10*time.Minute).Err()
	if err != nil {
		log.Printf("SetVerificationCode: Could not set value in Redis: %v", err)
		return 0, err
	}
	return code, nil
}

func (r *Repository) CheckVerificationCode(ctx context.Context, email, code string) error {
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

func key(userID, fp string) string {
	return fmt.Sprintf("refresh:%s:%s", userID, fp)
}

func (r *Repository) SaveToken(ctx context.Context, userId, fingerPrint, token string, ttl time.Duration) error {
	redisKey := key(userId, fingerPrint)
	err := r.redisExplorer.RDB.Set(ctx, redisKey, token, ttl).Err()
	if err != nil {
		log.Printf("SaveToken: Could not save value in Redis: %v", err)
		return err
	}
	return nil
}

func (r *Repository) GetToken(ctx context.Context, userId, fingerPrint string) (string, error) {
	redisKey := key(userId, fingerPrint)
	token, err := r.redisExplorer.RDB.Get(ctx, redisKey).Result()
	if err != nil {
		log.Printf("GetToken: Could not get value in Redis: %v", err)
		return "", err
	}
	return token, nil
}
