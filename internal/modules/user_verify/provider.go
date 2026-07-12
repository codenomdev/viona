package user_verify

import (
	"github.com/codenomdev/viona/internal/modules/user_verify/repository"
	"github.com/codenomdev/viona/internal/modules/user_verify/service"
	"github.com/google/wire"
)

var Provider wire.ProviderSet = wire.NewSet(
	repository.NewRepository,
	service.NewService,
)
