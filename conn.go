package asql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// OpenDB automatically configures a bun.DB instance with postgresSQL drivers.
// It returns the database, along with a cleaning function, whose execution can be deferred for a graceful shutdown.
func OpenDB(dsn string) (*bun.DB, func(), error) {
	driver := pgdriver.WithDSN(dsn)
	connector := pgdriver.NewConnector(driver)
	sqldb := sql.OpenDB(connector)
	database := bun.NewDB(sqldb, pgdialect.New())

	// Closing function to be deferred. Errors are ignored, because they are not relevant anymore when the server
	// shuts down.
	cleaner := func() {
		_ = database.Close()
		_ = sqldb.Close()
	}

	// Wait for the database to be fully operational before allowing interactions.
	err := database.Ping()
	for i := 0; i < 3 && err != nil; i++ {
		time.Sleep(1 * time.Second)
		err = database.Ping()
	}

	if err != nil {
		cleaner()
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}

	return database, cleaner, nil
}
