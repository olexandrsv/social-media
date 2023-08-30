package main

import (
	"context"
	"log"
	"net/http"
	"social-media/common"
	"social-media/users"
	"social-media/users/endpoint"
	"social-media/users/transport"

	"github.com/jackc/pgx/v5/pgxpool"
)

const url = `postgresql://postgres:database@localhost:5432/media?sslmode=disable`

func main() {
	db := InitPostgreSQL(url) 
	repo := users.NewRepository(db)
	auth := common.NewAuthClient()
	s := users.NewService(repo, auth)
	endpoints := endpoint.NewEndpoints(s)
	handler := transport.NewHTTPServer(endpoints)

	err := http.ListenAndServe(":8081", handler)
	if err != nil{
		log.Println(err)
	}
}

func InitPostgreSQL(url string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Println(err)
		panic("Unable to connect to database")
	}
	return pool
}