package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"simple_bank.sqlc.dev/app/api"
	db "simple_bank.sqlc.dev/app/db/sqlc"
)

const (
	dbSource      = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	ctx := context.Background()

	conn, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
