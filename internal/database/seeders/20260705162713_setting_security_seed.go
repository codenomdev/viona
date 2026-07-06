package seeders

import (
	"github.com/codenomdev/viona/internal/modules/setting/constant"
	"github.com/codenomdev/viona/internal/modules/setting/domain"
	"github.com/codenomdev/viona/pkg/util"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	s := &gormigrate.Migration{}

	s.ID = "20260705162713_setting_security_seed"

	s.Migrate = func(db *gorm.DB) error {
		securityCfg := []domain.SiteSetting{
			{
				Group:  constant.GroupSecurity,
				Key:    "allow_new_registrations",
				Values: util.ToJSON(true),
			},
			{
				Group:  constant.GroupSecurity,
				Key:    "allow_email_registrations",
				Values: util.ToJSON(true),
			},
			{
				Group:  constant.GroupSecurity,
				Key:    "allow_user_recover",
				Values: util.ToJSON(true),
			},
			{
				Group:  constant.GroupSecurity,
				Key:    "allow_password_login",
				Values: util.ToJSON(true),
			},
		}
		if err := db.Debug().Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&securityCfg).Error; err != nil {
			LogError(s.ID, err)
		}

		LogSuccess(s.ID)
		return nil
	}

	s.Rollback = func(db *gorm.DB) error {
		// TODO: implement rollback

		LogSuccess(s.ID)
		return nil
	}

	AddSeed(s)
}
