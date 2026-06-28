//go:build wireinject
// +build wireinject

package codenomcmd

import (
	"context"

	"github.com/codenomdev/viona/internal/apps"
	"github.com/codenomdev/viona/internal/apps/routes"
	"github.com/codenomdev/viona/internal/modules"
	"github.com/codenomdev/viona/pkg/config"
	"github.com/codenomdev/viona/pkg/log"
	"github.com/google/wire"
)

func initApplication(
	ctx context.Context,
	cfg *config.Config,
	log log.Logger,
) (*apps.AppServer, func(), error) {
	wire.Build(
		routes.NewApiRoutes,
		routes.NewUIRoutes,
		modules.Provider,
		apps.NewApp,
	)
	return nil, nil, nil
}
