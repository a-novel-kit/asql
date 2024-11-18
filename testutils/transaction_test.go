package asqltest_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	databasemocks "github.com/a-novel-kit/asql/mocks/migrations"
	asqltest "github.com/a-novel-kit/asql/testutils"
)

func TestTransaction(t *testing.T) {
	db, cleaner, err := asqltest.OpenTestDB(&databasemocks.MigrationsAll)
	require.NoError(t, err)
	defer cleaner()

	tx := asqltest.BeginTestTX(db, []*databasemocks.Table1Model{
		{
			ID:   1,
			Name: "foo",
		},
	})

	// Fixtures should be available through tx.
	var model databasemocks.Table1Model
	require.NoError(t, tx.NewSelect().Model(&model).Where("id = ?", 1).Scan(context.Background()))
	require.Equal(t, "foo", model.Name)

	// Fixtures should not be available through db.
	require.Error(t, db.NewSelect().Model(&databasemocks.Table1Model{}).Where("id = ?", 1).Scan(context.Background()))

	// Create a new model in transaction.
	_, err = tx.NewInsert().
		Model(&databasemocks.Table1Model{
			ID:   2,
			Name: "bar",
		}).
		Exec(context.Background())
	require.NoError(t, err)

	// New model should be available through tx.
	require.NoError(t, tx.NewSelect().Model(&model).Where("id = ?", 2).Scan(context.Background()))
	require.Equal(t, "bar", model.Name)

	// New model should not be available through db.
	require.Error(t, db.NewSelect().Model(&databasemocks.Table1Model{}).Where("id = ?", 2).Scan(context.Background()))

	// Rollback transaction.
	asqltest.RollbackTestTX(tx)

	// Fixtures should not be available through db.
	require.Error(t, db.NewSelect().Model(&databasemocks.Table1Model{}).Where("id = ?", 1).Scan(context.Background()))

	// New model should not be available through db.
	require.Error(t, db.NewSelect().Model(&databasemocks.Table1Model{}).Where("id = ?", 2).Scan(context.Background()))
}
