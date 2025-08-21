package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore returns a new SQLStore
func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: pool,
		Queries:  New(pool),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// отримаємо окреме з'єднання з пулу
	conn, err := store.db.(*pgxpool.Pool).Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// починаємо транзакцію
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)

	// виконуємо логіку
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
