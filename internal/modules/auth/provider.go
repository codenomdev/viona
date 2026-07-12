package auth

import (
	"github.com/codenomdev/viona/internal/modules/auth/handler"
	"github.com/codenomdev/viona/internal/modules/auth/route"
	"github.com/codenomdev/viona/internal/modules/auth/service"
	"github.com/google/wire"
)

var Provider wire.ProviderSet = wire.NewSet(
	service.NewService,
	handler.NewHandler,
	route.NewAuthRoute,
)
