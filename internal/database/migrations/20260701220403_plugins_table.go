package migrations

import (
	"github.com/codenomdev/viona/internal/modules/plugin/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{}

	m.ID = "20260701220403_plugins_table"

	m.Migrate = func(db *gorm.DB) error {
		if err := db.Migrator().CreateTable(
			&domain.PluginConfig{},
		); err != nil {
			LogError(m.ID, err)
			return err
		}

		LogSuccess(m.ID)
		return nil
	}

	m.Rollback = func(db *gorm.DB) error {
		if err := db.Migrator().DropTable(
			&domain.PluginConfig{},
		); err != nil {
			LogError(m.ID, err)
			return err
		}

		LogSuccess(m.ID, true)
		return nil
	}

	AddMigration(m)
}
