package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codenomdev/viona/internal/database/seeders"
	"github.com/codenomdev/viona/pkg/config"
	baseGorm "github.com/codenomdev/viona/pkg/db/gorm"
	"github.com/codenomdev/viona/pkg/log"
	"go.uber.org/zap"
)

func RunSeeders(ctx context.Context, tableName string, viewAll bool) error {
	logger := log.FromContext(ctx)
	cfg := config.FromContext(ctx)
	db := baseGorm.FromContext(ctx)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.HOST.CONTEXT_TIMEOUT)*time.Second)
	defer cancel()

	// Trap SIGTERM/SIGINT for graceful shutdown
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	seeders.NewSeed(db, logger)
	logger.Info("🚀 Initialize DB Seeder...")

	if viewAll {
		seeders.GetAllIDs()
		return nil
	}

	if tableName != "" {
		logger.Info("Running db-seed with", zap.String("tablename", tableName))
		if err := seeders.ApplyWithID(tableName); err != nil {
			logger.Error("Failed to apply db seeder with tablename", zap.Error(err))

			return err
		}
	} else {
		logger.Info("Running db-seed all seeders...")
		// apply start db-seed
		if err := seeders.Apply(); err != nil {
			logger.Error("Failed to apply db seeder", zap.Error(err))

			return err
		}
	}

	return nil
}
