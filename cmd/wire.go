//go:build wireinject
// +build wireinject

package cmd

import (
	"context"

	"github.com/codenomdev/viona/internal/apps"
	"github.com/codenomdev/viona/internal/apps/routes"
	"github.com/codenomdev/viona/internal/modules"
	"github.com/codenomdev/viona/pkg/config"
	"github.com/codenomdev/viona/pkg/log"
	"github.com/codenomdev/viona/pkg/translator"
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
		translator.Provider,
		apps.NewApp,
	)
	return nil, nil, nil
}
