package models

import "github.com/jackc/pgx/v5"

type Model struct {
	Conn *pgx.Conn
}
