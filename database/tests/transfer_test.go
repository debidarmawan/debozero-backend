package database_test

import (
	"context"
	"testing"
	"time"

	database "github.com/debidarmawan/debozero-backend/database/sqlc"
	"github.com/stretchr/testify/assert"
)

func createRandomAccount(userID int32, t *testing.T) database.Account {
	arg := database.CreateAccountParams{
		UserID:   userID,
		Balance:  200,
		Currency: "USD",
	}

	account, err := testQuery.CreateAccount(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, account)
	assert.Equal(t, account.Balance, arg.Balance)
	assert.Equal(t, account.Currency, arg.Currency)
	assert.Equal(t, account.UserID, arg.UserID)
	assert.WithinDuration(t, account.CreatedAt, time.Now(), 2*time.Second)

	return account
}

func TestTransfer(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	account1 := createRandomAccount(int32(user1.ID), t)
	account2 := createRandomAccount(int32(user2.ID), t)

	arg := database.CreateTransferParams{
		FromAccountID: int32(account1.ID),
		ToAccountID:   int32(account2.ID),
		Amount:        10,
	}

	txResponseChan := make(chan database.TransferTxResponse)
	errorChan := make(chan error)
	count := 10

	for i := 0; i < 3; i++ {
		go func() {
			tx, err := testQuery.TransferTx(context.Background(), arg)
			errorChan <- err
			txResponseChan <- tx
		}()
	}

	for x := 0; x < count; x++ {
		err := <-errorChan
		tx := <-txResponseChan

		assert.NoError(t, err)
		assert.NotEmpty(t, tx)

		// test transfer
		assert.Equal(t, tx.Transfer.FromAccountID, arg.FromAccountID)
		assert.Equal(t, tx.Transfer.ToAccountID, arg.ToAccountID)
		assert.Equal(t, tx.Transfer.Amount, arg.Amount)

		// test entry
		// entry in
		assert.Equal(t, tx.EntryIn.AccountID, arg.ToAccountID)
		assert.Equal(t, tx.EntryIn.Amount, arg.Amount)

		// entry out
		assert.Equal(t, tx.EntryOut.AccountID, arg.FromAccountID)
		assert.Equal(t, tx.EntryOut.Amount, -1*arg.Amount)
	}

	newAccount1, err := testQuery.GetAccountByID(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccount1)

	newAccount2, err := testQuery.GetAccountByID(context.Background(), account2.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccount2)

	newAmount := float64(count * int(arg.Amount))
	assert.Equal(t, newAccount1.Balance, account1.Balance-newAmount)
	assert.Equal(t, newAccount2.Balance, account1.Balance+newAmount)
}
