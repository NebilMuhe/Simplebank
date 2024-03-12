package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	query := New(tx)

	err = fn(query)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error %v, rb error %v", err, rbErr)
		}

		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)
		fmt.Println(txName, "Create transfer")

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "Create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "Create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "get account 1 for update")
		// account1, err := q.GetAccountForUpdate(context.Background(), arg.FromAccountId)

		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "update account1 balance")
		// result.FromAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
		// 	ID:      arg.FromAccountId,
		// 	Balance: account1.Balance - arg.Amount,
		// })

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMany(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
			if err != nil {
				return err
			}
			// result.FromAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			// 	ID:     arg.FromAccountId,
			// 	Amount: -arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }

			// result.ToAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			// 	ID:     arg.ToAccountId,
			// 	Amount: arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
		} else {

			result.ToAccount, result.FromAccount, err = addMany(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
			if err != nil {
				return err
			}
			// result.ToAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			// 	ID:     arg.ToAccountId,
			// 	Amount: arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }

			// result.FromAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			// 	ID:     arg.FromAccountId,
			// 	Amount: -arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
		}

		// fmt.Println(txName, "get account 2 for update")
		// account2, err := q.GetAccountForUpdate(context.Background(), arg.ToAccountId)

		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "update account2 balance")
		// result.ToAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
		// 	ID:      arg.ToAccountId,
		// 	Balance: account2.Balance + arg.Amount,
		// })

		return nil
	})

	return result, err
}

func addMany(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
