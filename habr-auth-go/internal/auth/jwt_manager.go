package auth

import (
	"errors"
	"os"
	"time"
)

type JwtManager struct {
	accessTtl  time.Duration
	refreshTtl time.Duration
	secret     []byte
}

func NewJwtManager(secret []byte, accessTtl, refreshTtl time.Duration) *JwtManager {
	return &JwtManager{secret: secret, accessTtl: accessTtl, refreshTtl: refreshTtl}
}

func GetJwtInfo() ([]byte, time.Duration, time.Duration, error) {
	accessTime := 10 * time.Minute
	refreshTime := 24 * 7 * time.Hour
	key := os.Getenv("JWT_KEY_WORD")
	if key == "" {
		return nil, 0, 0, errors.New("JWT_KEY_WORD not set")
	}
	keyWord := []byte(key)
	return keyWord, accessTime, refreshTime, nil
}
