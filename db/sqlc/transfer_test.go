package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simple_bank.sqlc.dev/app/util"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T, from *int64, to *int64) Transfer {
	var fromId int64
	var toId int64
	if from != nil && to != nil {
		fromId = *from
		toId = *to
	} else {
		fromId = util.RandomInt(5, 11)
		toId = util.RandomInt(1, 6)
	}

	arg := CreateTransferParams{
		FromAccountID: fromId,
		ToAccountID:   toId,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testStore.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.CreatedAt)
	require.NotZero(t, transfer.ID)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t, nil, nil)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t, nil, nil)
	transfer2, err := testStore.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
	require.WithinDuration(t, transfer1.CreatedAt.Time, transfer2.CreatedAt.Time, time.Second)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, &account1.ID, &account2.ID)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testStore.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfers)
		require.Equal(t, account1.ID, transfer.FromAccountID)
	}
}
