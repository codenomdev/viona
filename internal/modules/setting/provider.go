package setting

import (
	"github.com/codenomdev/viona/internal/modules/setting/repository"
	"github.com/codenomdev/viona/internal/modules/setting/service"
	"github.com/google/wire"
)

var Provider wire.ProviderSet = wire.NewSet(
	service.NewService,
	repository.NewRepository,
	NewHandler,
	NewRoute,
)
