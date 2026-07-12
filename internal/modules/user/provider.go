package user

import (
	"github.com/codenomdev/viona/internal/modules/user/repository"
	"github.com/codenomdev/viona/internal/modules/user/service"
	"github.com/google/wire"
)

var Provider wire.ProviderSet = wire.NewSet(
	repository.NewRepository,
	service.NewService,
)
