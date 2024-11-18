package asqltest

import (
	"context"

	"github.com/uptrace/bun"
)

func BeginTestTX[T any](database bun.IDB, fixtures []T) bun.Tx {
	transaction, err := database.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	for _, fixture := range fixtures {
		_, err := transaction.NewInsert().Model(fixture).Exec(context.Background())
		if err != nil {
			panic(err)
		}
	}

	return transaction
}

func RollbackTestTX(transaction bun.Tx) {
	_ = transaction.Rollback()
}
