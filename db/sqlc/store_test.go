package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)

		var fromAccountId int64
		var toAccountId int64

		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		} else {
			fromAccountId = account1.ID
			toAccountId = account2.ID
		}
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	/*
		// check results
		existed := make(map[int]bool)

		for i := 0; i < n; i++ {
			err := <-errs
			require.NoError(t, err)

			result := <-results
			require.NotEmpty(t, result)

			// check transfer
			transfer := result.Transfer
			require.NotEmpty(t, transfer)
			require.Equal(t, account1.ID, transfer.FromAccountID)
			require.Equal(t, account2.ID, transfer.ToAccountID)
			require.Equal(t, amount, transfer.Amount)
			require.NotZero(t, transfer.ID)
			require.NotZero(t, transfer.CreatedAt)

			_, err = store.GetTransfer(context.Background(), transfer.ID)
			require.NoError(t, err)

			// check entries
			fromEntry := result.FromEntry
			require.NotEmpty(t, fromEntry)
			require.Equal(t, account1.ID, fromEntry.AccountID)
			require.Equal(t, -amount, fromEntry.Amount)
			require.NotZero(t, fromEntry.ID)
			require.NotZero(t, fromEntry.CreatedAt)

			_, err = store.GetEntry(context.Background(), fromEntry.ID)
			require.NoError(t, err)

			toEntry := result.ToEntry
			require.NotEmpty(t, toEntry)
			require.Equal(t, account2.ID, toEntry.AccountID)
			require.Equal(t, amount, toEntry.Amount)
			require.NotZero(t, toEntry.ID)
			require.NotZero(t, toEntry.CreatedAt)

			_, err = store.GetEntry(context.Background(), toEntry.ID)
			require.NoError(t, err)

			// check accounts
			fromAccount := result.FromAccount
			require.NotEmpty(t, fromAccount)
			require.Equal(t, account1.ID, fromAccount.ID)

			toAccount := result.ToAccount
			require.NotEmpty(t, toAccount)
			require.Equal(t, account2.ID, toAccount.ID)

			// check balances
			fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

			diff1 := account1.Balance - fromAccount.Balance
			diff2 := toAccount.Balance - account2.Balance
			require.Equal(t, diff1, diff2)
			require.True(t, diff1 > 0)
			require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

			k := int(diff1 / amount)
			require.True(t, k >= 1 && k <= n)
			require.NotContains(t, existed, k)
			existed[k] = true
		}
	*/

	// check the final updated acccount
	updateAccount1, err := testQueries.GetAccountForUpdate(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccountForUpdate(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAccount1.Balance, updateAccount2.Balance)
	//require.Equal(t, int(account1.Balance)-n*int(amount), int(updateAccount1.Balance))
	//require.Equal(t, int(account2.Balance)+n*int(amount), int(updateAccount2.Balance))
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
