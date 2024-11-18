package databasemocks

import (
	"embed"

	"github.com/uptrace/bun"
)

//go:embed 20200101120000_migration_1.down.sql 20200101120000_migration_1.up.sql 20200101130000_migration_2.down.sql 20200101130000_migration_2.up.sql 20200101140000_migration_3.down.sql 20200101140000_migration_3.up.sql
var MigrationsAll embed.FS

//go:embed 20200101120000_migration_1.down.sql 20200101120000_migration_1.up.sql 20200101130000_migration_2.down.sql 20200101130000_migration_2.up.sql
var MigrationsGroup1 embed.FS

type Table1Model struct {
	bun.BaseModel `bun:"table1"`

	ID   int64  `bun:"id,pk"`
	Name string `bun:"name"`
}

type Table2Model struct {
	bun.BaseModel `bun:"table2"`

	ID   int64  `bun:"id,pk"`
	Name string `bun:"name"`
}

type Table3Model struct {
	bun.BaseModel `bun:"table3"`

	ID   int64  `bun:"id,pk"`
	Name string `bun:"name"`
}
