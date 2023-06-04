package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var PostgreConn *pgxpool.Pool

func InitPostgreSQL(url string) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Println(err)
		panic("Unable to connect to database")
	}
	PostgreConn = pool
}
