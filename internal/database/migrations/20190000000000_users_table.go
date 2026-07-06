package migrations

import (
	"github.com/codenomdev/viona/internal/modules/user/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{}

	m.ID = "20260702082851_users_table"

	m.Migrate = func(db *gorm.DB) error {
		if err := db.Migrator().CreateTable(&domain.User{}); err != nil {
			LogError(m.ID, err)
			return err
		}

		LogSuccess(m.ID)
		return nil
	}

	m.Rollback = func(db *gorm.DB) error {
		if err := db.Migrator().DropTable(&domain.User{}); err != nil {
			LogError(m.ID, err)
			return err
		}

		LogSuccess(m.ID, true)
		return nil
	}

	AddMigration(m)
}
