package user_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qwerun/habr-auth-go/internal/models"
	"log"
	"math/rand"
	"time"
)

var (
	ErrEmailAlreadyExists    = errors.New("user with this email already exists")
	ErrNicknameAlreadyExists = errors.New("user with this nickname already exists")
)

func (r *Repository) Create(user *models.User) (string, error) {
	query := `
		INSERT INTO users (email, password_hash, nickname)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id string
	err := r.explorer.DB.QueryRowContext(
		context.Background(),
		query,
		user.Email,
		user.PasswordHash,
		user.Nickname,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return "", ErrEmailAlreadyExists
			case "users_nickname_key":
				return "", ErrNicknameAlreadyExists
			}
		}
		log.Printf("Failed registration insert error: %v", err)
		return "", fmt.Errorf("rie")
	}
	return id, nil
}

func (r *Repository) SetVerificationCode(email string) (int, error) {
	ctx := context.Background()
	code := rand.Intn(900000) + 100000
	err := r.redisExplorer.RDB.Set(ctx, email, code, 10*time.Minute).Err()
	if err != nil {
		log.Printf("Could not set value in Redis: %v", err)
		return 0, err
	}
	return code, nil
}
