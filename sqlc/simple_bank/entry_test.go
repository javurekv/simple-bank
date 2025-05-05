package simple_bank

import (
	"context"
	"github.com/stretchr/testify/require"
	"simple_bank.sqlc.dev/app/util"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T, accountID *int64) Entry {
	var id int64
	if accountID != nil {
		id = *accountID
	} else {
		id = util.RandomInt(1, 100) // fallback if nil
	}

	arg := CreateEntryParams{
		AccountID: id,
		Amount:    util.RandomMoney(),
	}

	entry, err := testStore.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)

	return entry
}
func TestCreateEntries(t *testing.T) {
	createRandomEntry(t, nil)
}

func TestGetEntries(t *testing.T) {
	entry1 := createRandomEntry(t, nil)
	entry2, err := testStore.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, &account.ID)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testStore.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}
