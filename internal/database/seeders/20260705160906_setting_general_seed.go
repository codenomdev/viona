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

	s.ID = "20260705160906_setting_general_seed"

	s.Migrate = func(db *gorm.DB) error {
		generalCfg := []domain.SiteSetting{
			// language
			{
				Group:  constant.GroupGeneral,
				Key:    "language",
				Values: util.ToJSON("en_US"),
			},
			// timezone
			{
				Group:  constant.GroupGeneral,
				Key:    "timezone",
				Values: util.ToJSON("UTC"),
			},
		}

		if err := db.Debug().Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&generalCfg).Error; err != nil {
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
