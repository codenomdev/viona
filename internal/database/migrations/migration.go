package migrations

import (
	"fmt"
	"sort"

	"github.com/codenomdev/viona/pkg/log"
	"github.com/go-gormigrate/gormigrate/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// migration func
var Migration *gormigrate.Gormigrate

// migration lists
var MigrationsList = []*gormigrate.Migration{}

// Logger
var Logger log.Logger

// Initialize db migrate.
func NewMigrate(db *gorm.DB, log log.Logger) {
	Migration = gormigrate.New(db, gormigrate.DefaultOptions, MigrationsList)
	Logger = log
}

// Add register db migration.
func AddMigration(migration *gormigrate.Migration) {
	MigrationsList = append(MigrationsList, migration)
}

// Get all db migration ids.
func GetAllIDs() {
	fmt.Println("List tablename migration(s):")
	// sort asc to desc
	sortMigrations()
	for _, v := range MigrationsList {
		fmt.Printf("%s\n", v.ID)
	}
}

// Running migration with ID.
func MigrateWithID(tablename string) error {
	return Migration.MigrateTo(tablename)
}

// Running migration.
func Migrate() error {
	// sort asc to desc
	sortMigrations()
	return Migration.Migrate()
}

// Rollback db migration with ID.
func RollbackWithID(tablename string) error {
	return Migration.RollbackTo(tablename)
}

// Rollback db migration
// RollbackLast undo the last migration.
func Rollback() error {
	// sort a < b
	sortMigrations()
	return Migration.RollbackLast()
}

// Log success.
func LogSuccess(tablename string, rollback ...bool) {
	if len(rollback) > 0 && rollback[0] {
		Logger.Info("Rolled back", zap.String("migration:", tablename))
	}
	Logger.Info("Applied ", zap.String("migration:", tablename))
}

func LogError(tablename string, err error, rollback ...bool) {
	if len(rollback) > 0 && rollback[0] {
		Logger.Error(fmt.Sprintf("Failed to rollback migration: %s, error: %s", tablename, err))
	}
	Logger.Error(fmt.Sprintf("Failed to apply migration: %s, error: %s", tablename, err))
}

func GetPendingMigrations(db *gorm.DB) []*gormigrate.Migration {
	applied := map[string]bool{}
	var appliedIDs []string

	// Default migration table
	// tableName := "gorp_migrations"
	// if Migration != nil && Migration.Options.TableName != "" {
	// 	tableName = Migration.Options.TableName
	// }

	// Get applied migration IDs from DB
	if err := db.Table("migrations").Select("id").Scan(&appliedIDs).Error; err != nil {
		Logger.Error(fmt.Sprintf("Failed to read migration table '%s': %v", "migrations", err))
	}

	for _, id := range appliedIDs {
		applied[id] = true
	}

	sortMigrations()
	var pending []*gormigrate.Migration

	for _, m := range MigrationsList {
		if !applied[m.ID] {
			pending = append(pending, m)
		}
	}

	return pending
}

// Sort migration lists.
func sortMigrations() {
	// sort asc to desc
	sort.Slice(MigrationsList[:], func(i, j int) bool {
		return MigrationsList[i].ID < MigrationsList[j].ID
	})
}
