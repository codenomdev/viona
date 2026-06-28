package seeders

import (
	"fmt"
	"sort"

	"github.com/codenomdev/viona/pkg/log"
	"github.com/go-gormigrate/gormigrate/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Seed var
var Seed *gormigrate.Gormigrate

// Seeder lists
var SeedLists []*gormigrate.Migration

// Logger
var Logger log.Logger

// Initialize database seeder.
func NewSeed(db *gorm.DB, log log.Logger) {
	Seed = gormigrate.New(db, gormigrate.DefaultOptions, SeedLists)
	Logger = log
}

// Add register into seeder.
func AddSeed(seed *gormigrate.Migration) {
	SeedLists = append(SeedLists, seed)
}

// Get all db seed ids.
func GetAllIDs() {
	fmt.Println("List tablename seeder(s):")
	// sort asc to desc
	sortSeeders()
	for _, v := range SeedLists {
		fmt.Printf("%s\n", v.ID)
	}
}

// Running seeder with ID name.
func ApplyWithID(tablename string) error {
	return Seed.MigrateTo(tablename)
}

// Running seeder.
func Apply() error {
	// sort asc to desc
	sortSeeders()
	return Seed.Migrate()
}

// Log info success seeders.
func LogSuccess(tablename string) {
	Logger.Info("Applied seed for", zap.String("tablename", tablename))
}

// Log error seeders.
func LogError(tablename string, err error) {
	Logger.Error("Failed to apply", zap.String("seed", tablename), zap.Error(err))
}

// Sort seeder lists.
func sortSeeders() {
	// sort asc to desc
	sort.Slice(SeedLists[:], func(i, j int) bool {
		return SeedLists[i].ID < SeedLists[j].ID
	})
}
