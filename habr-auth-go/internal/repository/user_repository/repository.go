package user_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qwerun/habr-auth-go/internal/models"
	"log"
)

var (
	ErrEmailAlreadyExists    = errors.New("user with this email already exists")
	ErrNicknameAlreadyExists = errors.New("user with this nickname already exists")
	ErrCodeCheckNotFound     = errors.New("verification code not found")
	ErrVerifyAccount         = errors.New("verify account error")
	ErrBadRequest            = errors.New("Bad request")
)

func (r *Repository) Create(ctx context.Context, user *models.User) (string, error) {
	query := `
		INSERT INTO users (email, password_hash, nickname)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id string
	err := r.explorer.DB.QueryRowContext(
		ctx,
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
		log.Printf("Create: Failed registration insert error: %v", err)
		return "", fmt.Errorf("rie")
	}
	return id, nil
}

func (r *Repository) VerifiedAccount(ctx context.Context, email string) error {
	query := `
        UPDATE users
        SET is_verified = true
        WHERE email = $1
    `
	res, err := r.explorer.DB.ExecContext(ctx, query, email)
	if err != nil {
		log.Printf("VerifiedAccount: unable to mark user %q as verified: %w", email, err)
		return ErrVerifyAccount
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("VerifiedAccount: could not get RowsAffected for user %q: %w", email, err)
		return ErrVerifyAccount
	}
	if rows == 0 {
		log.Printf("VerifiedAccount: user not found %q: %w", email, err)
		return ErrVerifyAccount
	}

	return nil
}

func (r *Repository) GetPassHash(ctx context.Context, email string) (string, error) {
	query := `
        SELECT password_hash
          FROM users
         WHERE email = $1
    `
	var hash string

	err := r.explorer.DB.QueryRowContext(ctx, query, email).Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("GetPassHash: no rows: %w", err)
			return "", ErrBadRequest
		}
		log.Printf("GetPassHash: scan row: %w", err)
		return "", ErrBadRequest
	}

	return hash, nil
}

func (r *Repository) GetUserId(ctx context.Context, email string) (string, error) {
	query := `
        SELECT id
          FROM users
         WHERE email = $1 and is_verified = true
    `
	var id string

	err := r.explorer.DB.QueryRowContext(ctx, query, email).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("GetUserId: no rows: %w", err)
			return "", ErrBadRequest
		}
		log.Printf("GetUserId: scan row: %w", err)
		return "", ErrBadRequest
	}

	return id, nil
}

func (r *Repository) SetNewHash(ctx context.Context, email, hash string) error {
	query := `
        UPDATE users
        SET password_hash = $1
        WHERE email = $2
    `
	res, err := r.explorer.DB.ExecContext(ctx, query, hash, email)
	if err != nil {
		log.Printf("SetNewHash: unable to change hash password %q as verified: %w", email, err)
		return ErrBadRequest
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("SetNewHash: could not get RowsAffected for user %q: %w", email, err)
		return ErrBadRequest
	}
	if rows == 0 {
		log.Printf("SetNewHash: user not found %q: %w", email, err)
		return ErrBadRequest
	}

	return nil
}
