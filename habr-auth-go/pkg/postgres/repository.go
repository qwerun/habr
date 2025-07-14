package postgres

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Explorer struct {
	DB *sql.DB
}

func NewExplorer(db *sql.DB) *Explorer {
	return &Explorer{DB: db}
}
