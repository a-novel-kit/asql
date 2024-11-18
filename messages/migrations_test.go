package asqlmessages_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun/migrate"

	asqlmessages "github.com/a-novel-kit/asql/messages"
)

func TestMigrations(t *testing.T) {
	t.Run("Render", func(t *testing.T) {
		content := asqlmessages.NewMigrations([]migrate.Migration{
			{
				ID:         1,
				Name:       "20200101120000",
				Comment:    "_migration_1",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:         2,
				Name:       "20200101120000",
				Comment:    "_migration_2",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:         3,
				Name:       "20200101120000",
				Comment:    "_migration_3",
				GroupID:    2,
				MigratedAt: time.Date(2020, 1, 2, 13, 0, 0, 0, time.UTC),
			},
		}, 0)

		expectConsole := " ✓ Group 2\n" +
			"     - 20200101120000__migration_3 (2020-01-02T13:00:00Z)\n" +
			" ✓ Group 1\n" +
			"     - 20200101120000__migration_2 (2020-01-02T12:00:00Z)\n" +
			"     - 20200101120000__migration_1 (2020-01-02T12:00:00Z)\n"
		expectJSON := map[string]interface{}{
			"1": []interface{}{
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_2",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_1",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
			},
			"2": []interface{}{
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_3",
					"migrated_at": "2020-01-02T13:00:00Z",
				},
			},
		}

		require.Equal(t, expectConsole, content.RenderTerminal())
		require.Equal(t, expectJSON, content.RenderJSON())
	})

	t.Run("NoMigrations", func(t *testing.T) {
		content := asqlmessages.NewMigrations([]migrate.Migration{}, 0)

		require.Equal(t, "", content.RenderTerminal())
		require.Nil(t, content.RenderJSON())
	})

	t.Run("NonAppliedGroup", func(t *testing.T) {
		content := asqlmessages.NewMigrations([]migrate.Migration{
			{
				ID:         1,
				Name:       "20200101120000",
				Comment:    "_migration_1",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:         2,
				Name:       "20200101120000",
				Comment:    "_migration_2",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:      3,
				Name:    "20200101120000",
				Comment: "_migration_3",
				GroupID: 2,
			},
		}, 0)

		expectConsole := " ✓ Group 1\n" +
			"     - 20200101120000__migration_2 (2020-01-02T12:00:00Z)\n" +
			"     - 20200101120000__migration_1 (2020-01-02T12:00:00Z)\n" +
			" ✗ Group 2\n" +
			"     - 20200101120000__migration_3\n"
		expectJSON := map[string]interface{}{
			"1": []interface{}{
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_2",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_1",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
			},
			"2": []interface{}{
				map[string]interface{}{
					"name":    "20200101120000",
					"comment": "_migration_3",
				},
			},
		}

		require.Equal(t, expectConsole, content.RenderTerminal())
		require.Equal(t, expectJSON, content.RenderJSON())
	})

	t.Run("SpecialGroup0", func(t *testing.T) {
		content := asqlmessages.NewMigrations([]migrate.Migration{
			{
				ID:      1,
				Name:    "20200101120000",
				Comment: "_migration_1",
				GroupID: 0,
			},
			{
				ID:      2,
				Name:    "20200101120000",
				Comment: "_migration_2",
				GroupID: 0,
			},
			{
				ID:         3,
				Name:       "20200101120000",
				Comment:    "_migration_3",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
		}, 0)

		expectConsole := " ✓ Group 1\n" +
			"     - 20200101120000__migration_3 (2020-01-02T12:00:00Z)\n" +
			" No group\n" +
			"     - 20200101120000__migration_2\n" +
			"     - 20200101120000__migration_1\n"
		expectJSON := map[string]interface{}{
			"0": []interface{}{
				map[string]interface{}{
					"name":    "20200101120000",
					"comment": "_migration_2",
				},
				map[string]interface{}{
					"name":    "20200101120000",
					"comment": "_migration_1",
				},
			},
			"1": []interface{}{
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_3",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
			},
		}

		require.Equal(t, expectConsole, content.RenderTerminal())
		require.Equal(t, expectJSON, content.RenderJSON())
	})

	t.Run("LastApplied", func(t *testing.T) {
		content := asqlmessages.NewMigrations([]migrate.Migration{
			{
				ID:         1,
				Name:       "20200101120000",
				Comment:    "_migration_1",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:         2,
				Name:       "20200101120000",
				Comment:    "_migration_2",
				GroupID:    1,
				MigratedAt: time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:         3,
				Name:       "20200101120000",
				Comment:    "_migration_3",
				GroupID:    2,
				MigratedAt: time.Date(2020, 1, 2, 13, 0, 0, 0, time.UTC),
			},
		}, 2)

		expectConsole := " ✓ Group 2\n" +
			"     - 20200101120000__migration_3 (2020-01-02T13:00:00Z)\n" +
			" ✓ Group 1\n" +
			"     - 20200101120000__migration_2 (2020-01-02T12:00:00Z)\n" +
			"     - 20200101120000__migration_1 (2020-01-02T12:00:00Z)\n"
		expectJSON := map[string]interface{}{
			"1": []interface{}{
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_2",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_1",
					"migrated_at": "2020-01-02T12:00:00Z",
				},
			},
			"2": []interface{}{
				map[string]interface{}{
					"name":        "20200101120000",
					"comment":     "_migration_3",
					"migrated_at": "2020-01-02T13:00:00Z",
				},
			},
		}

		require.Equal(t, expectConsole, content.RenderTerminal())
		require.Equal(t, expectJSON, content.RenderJSON())
	})
}
