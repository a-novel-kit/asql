package asqlmessages

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/samber/lo"
	"github.com/uptrace/bun/migrate"

	"github.com/a-novel-kit/quicklog"
)

type migrationGroup struct {
	groupID    int64
	migrations []migrate.Migration
}

type migrationsMessage struct {
	// The list of discovered migrations.
	migrations []migrate.Migration
	// If set, the last applied migration will be highlighted.
	lastAppliedGroup int64

	quicklog.Message
}

// Return the underlying migrations, grouped and ordered by applied time.
func (migrations *migrationsMessage) getSortedMigrations() []migrationGroup {
	// The result we will return.
	var groups []migrationGroup

	// Migrations are already sorted by groupID, in reverse order. Since group IDs reflect the apply time,
	// we don't need extra sorting.
	currentGroup := migrationGroup{groupID: migrations.migrations[0].GroupID}

	for _, migration := range migrations.migrations {
		// We encountered a new migration group. Insert the current group to the output, and assign a new value
		// to currentGroup.
		if migration.GroupID != currentGroup.groupID {
			// Should not happen, but prevent adding empty groups to the output.
			if len(currentGroup.migrations) > 0 {
				groups = append(groups, currentGroup)
				currentGroup = migrationGroup{groupID: migration.GroupID}
			}
		}

		// Append the migration to the current group.
		currentGroup.migrations = append(currentGroup.migrations, migration)
	}

	// Make sure the last group is added to the output.
	if len(currentGroup.migrations) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

// Get the title of a migration group for display.
func (migrations *migrationsMessage) printGroupTitle(group migrationGroup) string {
	// Non-applied migrations have the group 0.
	if group.groupID == 0 {
		return lipgloss.NewStyle().Faint(true).Render("No group")
	}

	// With Bun, every migration in a group is either applied or un-applied.
	// This ensures a stable state of the database.
	isGroupApplied := group.migrations[0].MigratedAt != time.Time{}

	if !isGroupApplied {
		return lipgloss.NewStyle().Bold(true).Faint(true).Render(fmt.Sprintf("✗ Group %v", group.groupID))
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true).
		Faint(group.groupID != migrations.lastAppliedGroup).
		Render(fmt.Sprintf("✓ Group %v", group.groupID))
}

func (migrations *migrationsMessage) printMigrationItem(migration migrate.Migration) string {
	// The name of the migration is only its timestamp, and the comment the rest of the file name.
	// Concatenating thw 2 we get the name of the migration file, without the extension.
	migrationName := migration.Name + "_" + migration.Comment

	applied := migration.MigratedAt != time.Time{}

	// If the file has a migration date set, it has been migrated.
	if applied {
		// Show the migration date.
		migratedAt := " " + lipgloss.NewStyle().
			Faint(migration.GroupID != migrations.lastAppliedGroup).
			Render("("+migration.MigratedAt.Format(time.RFC3339)+")")

		return lipgloss.NewStyle().
			// If the migration is the last applied, highlight it with a different color.
			Foreground(lipgloss.Color("33")).
			Faint(migration.GroupID != migrations.lastAppliedGroup).
			Render(" "+migrationName) + migratedAt
	}

	return lipgloss.NewStyle().Faint(true).Render(" " + migrationName)
}

func (migrations *migrationsMessage) printGroup(group migrationGroup) *list.List {
	items := lo.Map(group.migrations, func(item migrate.Migration, _ int) string {
		return migrations.printMigrationItem(item)
	})

	// Render the list of migrations under the current group.
	return list.New(items).
		Enumerator(list.Dash).
		EnumeratorStyle(lipgloss.NewStyle().Faint(group.groupID != migrations.lastAppliedGroup))
}

func (migrations *migrationsMessage) RenderTerminal() string {
	if len(migrations.migrations) == 0 {
		return ""
	}

	// Disable enumerator for the list of groups.
	pList := list.New().
		Enumerator(func(_ list.Items, _ int) string { return "" }).
		Indenter(func(_ list.Items, _ int) string {
			return "    "
		})
	sorted := migrations.getSortedMigrations()

	// Populate the list with each group, and their underlying migrations.
	for _, group := range sorted {
		pList.Items(migrations.printGroupTitle(group), migrations.printGroup(group))
	}

	return pList.String() + "\n"
}

func (migrations *migrationsMessage) RenderJSON() map[string]interface{} {
	if len(migrations.migrations) == 0 {
		return nil
	}

	output := make(map[string]interface{})

	for _, migration := range migrations.migrations {
		mapKey := strconv.FormatInt(migration.GroupID, 10)
		if _, ok := output[mapKey]; !ok {
			output[mapKey] = make([]interface{}, 0)
		}

		elem := map[string]interface{}{
			"name":    migration.Name,
			"comment": migration.Comment,
		}

		applied := migration.MigratedAt != time.Time{}
		if applied {
			elem["migrated_at"] = migration.MigratedAt.Format(time.RFC3339)
		}

		output[mapKey] = append(output[mapKey].([]interface{}), elem)
	}

	return output
}

func NewMigrations(migrations []migrate.Migration, lastAppliedGroup int64) quicklog.Message {
	if len(migrations) > 0 {
		// Sort groups by groupID.
		slices.SortFunc(migrations, func(migrationA, migrationB migrate.Migration) int {
			// Sort by migration date, last migrated first.
			diff := int(migrationB.MigratedAt.UnixNano() - migrationA.MigratedAt.UnixNano())
			if diff == 0 {
				diff = strings.Compare(migrationB.Name+migrationB.Comment, migrationA.Name+migrationA.Comment)
			}
			if diff == 0 {
				diff = int(migrationB.GroupID - migrationA.GroupID)
			}

			return diff
		})
	}

	return &migrationsMessage{
		migrations:       migrations,
		lastAppliedGroup: lastAppliedGroup,
	}
}
