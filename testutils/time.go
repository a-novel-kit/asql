package asqltest

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

const setFakeNowFn = `
CREATE SCHEMA IF NOT EXISTS override;

GRANT ALL ON SCHEMA override TO public;
GRANT ALL ON SCHEMA override TO test;

CREATE OR REPLACE FUNCTION override.now() 
  RETURNS timestamptz IMMUTABLE PARALLEL SAFE AS 
$$
BEGIN
    return ?0::timestamptz;
END
$$ language plpgsql;

set search_path = override,pg_temp,"$user",public,pg_catalog;
`

const unsetFakeNowFn = `
DROP FUNCTION IF EXISTS override.now();
set search_path = pg_temp,"$user",public,pg_catalog;
`

func FreezeTime(db bun.IDB, date time.Time) error {
	// https://stackoverflow.com/questions/48243934/mocking-postgresql-now-function-for-testing
	_, err := db.ExecContext(context.Background(), setFakeNowFn, date)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func RestoreTime(db bun.IDB) error {
	_, err := db.ExecContext(context.Background(), unsetFakeNowFn)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
