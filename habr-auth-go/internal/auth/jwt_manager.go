package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type JwtManager struct {
	accessTtl  time.Duration
	RefreshTtl time.Duration
	secret     []byte
}

func NewJwtManager(secret []byte, accessTtl, refreshTtl time.Duration) *JwtManager {
	return &JwtManager{secret: secret, accessTtl: accessTtl, RefreshTtl: refreshTtl}
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

func (m *JwtManager) NewPair(uid, fingerPrint string) (string, string, error) {
	now := time.Now()
	accessClaims := jwt.RegisteredClaims{
		Subject:   uid,
		ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTtl)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        fingerPrint,
	}
	a := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	access, err := a.SignedString(m.secret)
	if err != nil {
		return "", "", fmt.Errorf("SignedString: Err:%v", err)
	}

	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("rand.Read: Err:%v", err)
	}
	refresh := base64.RawURLEncoding.EncodeToString(buf)
	return access, refresh, nil
}

func (m *JwtManager) ParseAccess(t string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(
		t,
		claims,
		func(_ *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithoutClaimsValidation(),
	)
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt == nil {
		return nil, errors.New("missing exp")
	}
	if time.Now().After(claims.ExpiresAt.Time) {
		return claims, jwt.ErrTokenExpired
	}

	return claims, nil
}
