package simple_bank

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

const dbSource = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"

var testPool *pgxpool.Pool
var testStore *Store

func TestMain(m *testing.M) {
	var err error

	testPool, err = pgxpool.New(context.Background(), dbSource)
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
