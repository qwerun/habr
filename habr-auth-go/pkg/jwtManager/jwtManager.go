package jwtManager

import (
	"os"
	"time"
)

type JwtExplorer struct {
	accessTtl  time.Duration
	refreshTtl time.Duration
	secret     []byte
}

func NewJwtManager(secret []byte, accessTtl, refreshTtl time.Duration) *JwtExplorer {
	return &JwtExplorer{secret: secret, accessTtl: accessTtl, refreshTtl: refreshTtl}
}

func NewJwtInfo() (accessTime time.Duration, refreshTime time.Duration, keyWord []byte) {
	accessTime = 10 * time.Minute
	refreshTime = 24 * 7 * time.Hour
	key := os.Getenv("JWT_KEY_WORD")
	keyWord = []byte(key)
	return
}
