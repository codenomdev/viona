package migrations

import (
	"github.com/codenomdev/viona/internal/modules/user_verify/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{}

	m.ID = "20190000000000_user_email_verify_table"

	m.Migrate = func(db *gorm.DB) error {
		if err := db.Migrator().CreateTable(&domain.EmailVerification{}); err != nil {
			LogError(m.ID, err)
			return err
		}

		LogSuccess(m.ID)
		return nil
	}

	m.Rollback = func(db *gorm.DB) error {
		if err := db.Migrator().DropTable(&domain.EmailVerification{}); err != nil {
			LogError(m.ID, err)
			return err
		}

		LogSuccess(m.ID, true)
		return nil
	}

	AddMigration(m)
}
