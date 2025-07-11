package auth

import (
	"github.com/golang-jwt/jwt"
	"time"
)

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["authorized"] = true
	claims["user"] = "username"
	return "", nil
}
