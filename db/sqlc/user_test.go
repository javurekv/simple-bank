package db

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"simple_bank.sqlc.dev/app/util"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}
func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text(sql.NullString{
			String: newFullName,
			Valid:  true,
		}),
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: pgtype.Text(sql.NullString{
			String: newEmail,
			Valid:  true,
		}),
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text(sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		}),
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text(sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		}),
		FullName: pgtype.Text(sql.NullString{
			String: newFullName,
			Valid:  true,
		}),
		Email: pgtype.Text(sql.NullString{
			String: newEmail,
			Valid:  true,
		}),
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
}

//func TestUpdateUser(t *testing.T) {
//	account1 := createRandomAccount(t)

//	arg := UpdateUserParams{
//		ID:      account1.ID,
//		Balance: util.RandomMoney(),
//	}
//
//	account2, err := testStore.UpdateAccount(context.Background(), arg)
//	require.NoError(t, err)
//	require.NotEmpty(t, account2)
//
//	require.Equal(t, account1.ID, account2.ID)
//	require.Equal(t, account1.Owner, account2.Owner)
//	require.Equal(t, arg.Balance, account2.Balance)
//	require.Equal(t, account1.Currency, account2.Currency)
//	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
//}

//
//func TestDeleteUser(t *testing.T) {
//	account1 := createRandomAccount(t)
//	err := testStore.DeleteAccount(context.Background(), account1.ID)
//	require.NoError(t, err)
//
//	account2, err := testStore.GetAccount(context.Background(), account1.ID)
//	require.Error(t, err)
//	require.ErrorIs(t, err, pgx.ErrNoRows)
//	require.Empty(t, account2)
//}
//
//func TestListUsers(t *testing.T) {
//	for i := 0; i < 10; i++ {
//		createRandomAccount(t)
//	}
//
//	arg := ListAccountsParams{
//		Limit:  5,
//		Offset: 5,
//	}
//
//	accounts, err := testStore.ListAccounts(context.Background(), arg)
//	require.NoError(t, err)
//	require.Len(t, accounts, 5)
//
//	for _, account := range accounts {
//		require.NotEmpty(t, account)
//	}
//}
