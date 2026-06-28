package modules

import (
	"github.com/codenomdev/viona/internal/modules/plugin"
	"github.com/codenomdev/viona/internal/modules/setting"
	"github.com/codenomdev/viona/internal/modules/static"
	"github.com/codenomdev/viona/pkg/db/gorm"
	"github.com/google/wire"
)

var Provider wire.ProviderSet = wire.NewSet(
	gorm.NewPgsqlDB,
	setting.Provider,
	static.Provider,
	plugin.Provider,
)
