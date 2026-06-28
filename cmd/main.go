package codenomcmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/codenomdev/viona/pkg/config"
	"github.com/codenomdev/viona/pkg/db/gorm"
	"github.com/codenomdev/viona/pkg/log"
	"go.uber.org/zap"
)

var (
	// app version, please dont modify this
	Version string = "0.0.0"

	// config path
	ConfigPath string

	// Tablename for seeder
	TableName string

	// ViewAll
	ViewAll bool

	// Rollback
	Rollback bool

	// place to build new viona
	buildDir string
	// plugins needed to build in viona application
	buildWithPlugins []string
	// build output path
	buildOutput string
	// Revision is the git short commit revision number
	// If built without a Git repository, this field will be empty.
	Revision = ""
	// Time is the build time of the project
	Time = ""
)

func Execute() {
	// Root application context
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	rootCmd.SetContext(ctx)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initLogger(ctx context.Context) (context.Context, func()) {
	zl := zap.Must(zap.NewProduction())

	cleanup := func() {
		_ = zl.Sync()
	}

	l := log.NewWithZap(zl)
	l.With(ctx, zap.String("version", Version))

	ctx = log.ToContext(ctx, l)

	return ctx, cleanup
}

func initConfig(ctx context.Context) (context.Context, error) {
	logger := log.FromContext(ctx)

	conf, err := config.LoadConfig(ConfigPath, logger)
	if err != nil {
		return ctx, fmt.Errorf("failed to load config: %w", err)
	}

	ctx = config.ToContext(ctx, conf)

	return ctx, nil
}

func initDB(ctx context.Context) (context.Context, func(), error) {
	logger := log.FromContext(ctx)
	cfg := config.FromContext(ctx)

	dbConn, cleanup, err := gorm.NewPgsqlDB(cfg, ctx)
	if err != nil {
		return ctx, nil, fmt.Errorf("failed to init db: %w", err)
	}

	logger.Info("database connected")

	ctx = gorm.ToContext(ctx, dbConn)

	return ctx, cleanup, nil
}

func runApp(ctx context.Context) {
	logger := log.FromContext(ctx)
	cfg := config.FromContext(ctx)
	server, cleanup, err := initApplication(
		ctx,
		cfg,
		logger,
	)
	if err != nil {
		logger.Error("failed to initialize application", zap.Error(err))
		os.Exit(1)
	}

	defer cleanup()

	logger.Info("starting server")

	if err := server.Start(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) &&
		!errors.Is(err, context.Canceled) {

		logger.Error("failed to start server", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}
