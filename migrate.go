package asql

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"

	"github.com/a-novel-kit/quicklog"
	"github.com/a-novel-kit/quicklog/messages"

	asqlmessages "github.com/a-novel-kit/asql/messages"
)

var (
	ErrDiscoverMigrations  = errors.New("failed to discover migrations")
	ErrCreateMigrator      = errors.New("failed to create migrator")
	ErrApplyMigrations     = errors.New("failed to apply migrations")
	ErrGetMigrationsStatus = errors.New("failed to get migrations status")
)

// Migrate looks for non-applied migrations, and applies them to the database.
func Migrate(database *bun.DB, sqlMigrations embed.FS, logger quicklog.Logger) error {
	loader := messages.NewLoader("discovering migrations...", &messages.LoaderConfigDefault)
	clean := logger.LogAnimated(loader)
	defer func() { go clean() }()

	// Discover existing migrations.
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(sqlMigrations); err != nil {
		loader.Error(ErrDiscoverMigrations)
		return fmt.Errorf("discover migrations: %w", err)
	}
	loader.Update("migrations successfully discovered, applying migrations...")

	migrator := migrate.NewMigrator(database, migrations)
	if err := migrator.Init(context.Background()); err != nil {
		loader.Error(ErrCreateMigrator)
		return fmt.Errorf("create migrator: %w", err)
	}

	// Run migrations.
	migrated, err := migrator.Migrate(context.Background())
	if err != nil {
		loader.Error(ErrApplyMigrations)
		return fmt.Errorf("apply migrations: %w", err)
	}

	applied, err := migrator.MigrationsWithStatus(context.Background())
	if err != nil {
		loader.Error(ErrGetMigrationsStatus)
		return fmt.Errorf("get migrations status: %w", err)
	}

	hasNewMigrations := migrated != nil && len(migrated.Migrations) > 0
	migrationsSubTitle := lo.TernaryF(
		hasNewMigrations,
		func() string {
			return fmt.Sprintf("%v new migrations applied in group %v", len(migrated.Migrations), migrated.ID)
		},
		func() string {
			return "No new migrations applied"
		},
	)

	loader.Nest(
		messages.NewTitle(
			"Migrations applied",
			migrationsSubTitle,
			asqlmessages.NewMigrations(applied, migrated.ID),
		),
	)
	loader.Success("migrations successfully applied.")

	// Great success.
	return nil
}
