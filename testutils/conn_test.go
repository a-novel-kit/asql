package asqltest_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"

	databasemocks "github.com/a-novel-kit/asql/mocks/migrations"
	"github.com/a-novel-kit/asql/testutils"
)

func TestConn(t *testing.T) {
	t.Run("NoMigrations", func(t *testing.T) {
		db, cleaner, err := asqltest.OpenTestDB(nil)
		require.NoError(t, err)
		defer cleaner()

		_, err = db.Exec("SELECT 1")
		require.NoError(t, err)
	})

	t.Run("WithMigrations", func(t *testing.T) {
		db, cleaner, err := asqltest.OpenTestDB(&databasemocks.MigrationsAll)
		require.NoError(t, err)
		defer cleaner()

		_, err = db.NewInsert().Model(&databasemocks.Table3Model{
			BaseModel: bun.BaseModel{},
			ID:        42,
			Name:      "foo",
		}).Exec(context.Background())
		require.NoError(t, err)

		res := databasemocks.Table3Model{}
		require.NoError(t, db.NewSelect().Model(&res).Where("id = ?", 42).Scan(context.Background()))
		require.Equal(t, "foo", res.Name)

		// Clear database should not fail.
		asqltest.ClearTestDB(db)

		// Database should be empty.
		require.Error(t, db.NewSelect().Model(&databasemocks.Table3Model{}).Where("id = ?", 42).Scan(context.Background()))
	})
}
