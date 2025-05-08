package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"simple_bank.sqlc.dev/app/util"
	"testing"
)

var testPool *pgxpool.Pool
var testStore Store

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testPool, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to test database:", err)
	}

	testStore = NewStore(testPool)

	// запустити тести
	code := m.Run()

	// закрити пул після завершення
	testPool.Close()

	os.Exit(code)
}
