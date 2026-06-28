package cmd

import (
	"context"

	"github.com/codenomdev/viona/pkg/log"
)

type BootstrapResult struct {
	Context context.Context
	Logger  log.Logger
	Cleanup func()
}

func bootstrap(ctx context.Context, withDB bool) (*BootstrapResult, error) {
	var cleanups []func()

	cleanup := func() {
		for i := len(cleanups) - 1; i >= 0; i-- {
			if cleanups[i] != nil {
				cleanups[i]()
			}
		}
	}

	var loggerCleanup func()

	ctx, loggerCleanup = initLogger(ctx)
	cleanups = append(cleanups, loggerCleanup)

	logger := log.FromContext(ctx)

	var err error

	ctx, err = initConfig(ctx)
	if err != nil {
		cleanup()
		return nil, err
	}

	if withDB {
		var dbCleanup func()

		ctx, dbCleanup, err = initDB(ctx)
		if err != nil {
			cleanup()
			return nil, err
		}

		cleanups = append(cleanups, dbCleanup)
	}

	return &BootstrapResult{
		Context: ctx,
		Logger:  logger,
		Cleanup: cleanup,
	}, nil
}
