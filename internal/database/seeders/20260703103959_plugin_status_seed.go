package seeders

import (
	pluginConst "github.com/codenomdev/viona/internal/modules/plugin/constant"
	"github.com/codenomdev/viona/internal/modules/setting/constant"
	"github.com/codenomdev/viona/internal/modules/setting/domain"
	"github.com/codenomdev/viona/pkg/util"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	s := &gormigrate.Migration{}

	s.ID = "20260703103959_plugin_status_seed"

	s.Migrate = func(db *gorm.DB) error {
		pluginCfg := &domain.SiteSetting{
			Group:  constant.GroupSystem,
			Key:    pluginConst.PluginStatus,
			Values: util.ToJSON([]map[string]any{}),
		}

		if err := db.Debug().Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&pluginCfg).Error; err != nil {
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
