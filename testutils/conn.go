package asqltest

import (
	"context"
	"embed"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/a-novel-kit/quicklog/loggers"

	"github.com/a-novel-kit/asql"
)

// OpenTestDB opens a connection to a test DB.
//
// The test DB must be available under the value stored in DSN.
func OpenTestDB(sqlMigrations *embed.FS) (*bun.DB, func(), error) {
	database, closer, err := asql.OpenDB(TestDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("open db: %w", err)
	}

	// Just in case something went wrong on latest run.
	ClearTestDB(database)
	if sqlMigrations == nil {
		return database, closer, nil
	}

	if err = asql.Migrate(database, *sqlMigrations, loggers.NewTerminal()); err != nil {
		closer()
		return nil, nil, fmt.Errorf("migrate: %w", err)
	}

	return database, closer, nil
}

func ClearTestDB(database *bun.DB) {
	ctx := context.Background()
	if _, err := database.ExecContext(ctx, "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"); err != nil {
		panic(err)
	}

	if _, err := database.ExecContext(ctx, "GRANT ALL ON SCHEMA public TO public;"); err != nil {
		panic(err)
	}
	if _, err := database.ExecContext(ctx, "GRANT ALL ON SCHEMA public TO test;"); err != nil {
		panic(err)
	}
}
