package modules

import (
	"github.com/codenomdev/viona/internal/modules/auth"
	"github.com/codenomdev/viona/internal/modules/plugin"
	"github.com/codenomdev/viona/internal/modules/setting"
	"github.com/codenomdev/viona/internal/modules/static"
	"github.com/codenomdev/viona/internal/modules/translator"
	"github.com/codenomdev/viona/internal/modules/user"
	"github.com/codenomdev/viona/internal/modules/user_verify"
	"github.com/codenomdev/viona/pkg/db/gorm"
	"github.com/google/wire"
)

var Provider wire.ProviderSet = wire.NewSet(
	gorm.NewPgsqlDB,
	setting.Provider,
	auth.Provider,
	user.Provider,
	user_verify.Provider,
	translator.Provider,
	static.Provider,
	plugin.Provider,
)
