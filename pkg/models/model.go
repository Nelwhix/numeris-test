package models

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Model struct {
	Conn *pgxpool.Pool
}
