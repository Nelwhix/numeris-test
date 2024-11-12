package pkg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func CreateDbConn() (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%v:%v@%v/%v",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST")+":"+os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
