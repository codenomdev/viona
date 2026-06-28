package codenomcmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codenomdev/viona/internal/database/migrations"
	"github.com/codenomdev/viona/pkg/config"
	baseGorm "github.com/codenomdev/viona/pkg/db/gorm"
	"github.com/codenomdev/viona/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// DRY-RUN
	DryRun bool
)

// runDBMigrations executes database migrations with graceful cancel + timeout
func runDBMigrations(ctx context.Context) error {
	logger := log.FromContext(ctx)
	cfg := config.FromContext(ctx)
	db := baseGorm.FromContext(ctx)

	// Create base timeout context
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.HOST.CONTEXT_TIMEOUT)*time.Second)
	defer cancel()

	// Trap SIGTERM/SIGINT for graceful shutdown
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("🚀 Initialize DB Migrations...")

	migrations.NewMigrate(db, logger)

	// Handle flag: view all
	if ViewAll {
		migrations.GetAllIDs()
		return nil
	}

	// Handle flag: dry run
	if DryRun {
		logger.Debug("Running in DRY-RUN mode — no changes will be applied.")
		return simulateMigration(ctx, db, logger)
	}

	// Specific table migration
	if TableName != "" {
		return runMigrationByTable(ctx, db, logger, TableName)
	}

	// Run all migrations
	return runMigrationAll(ctx, db, logger)
}

func runMigrationByTable(ctx context.Context, db *gorm.DB, logger log.Logger, table string) error {
	logger.Info("Running migration for", zap.String("table", table))

	if Rollback {
		if err := migrations.RollbackWithID(table); err != nil {
			return fail(logger, "rollback table failed", err)
		}
		logger.Info("Rollback successfully", zap.String("table:", table))
		return nil
	}

	if err := migrations.MigrateWithID(table); err != nil {
		return fail(logger, "migrate table failed", err)
	}

	logger.Info("Migrate successfully for", zap.String("table:", table))
	return nil
}

func runMigrationAll(ctx context.Context, db *gorm.DB, logger log.Logger) error {
	logger.Info("Running all migrations...")

	if Rollback {
		if err := migrations.Rollback(); err != nil {
			return fail(logger, "rollback all tables failed", err)
		}
		logger.Info("Rollback all tables successfully.")
		return nil
	}

	done := make(chan error, 1)
	go func() {
		done <- migrations.Migrate()
	}()

	select {
	case <-ctx.Done():
		logger.Error("Migration canceled or timed out.")
		return ctx.Err()
	case err := <-done:
		if err != nil {
			return fail(logger, "migration failed", err)
		}
	}

	logger.Info("Migrations completed successfully.")
	return nil
}

// simulateMigration is used in DRY-RUN mode (no database changes)
func simulateMigration(ctx context.Context, db *gorm.DB, logger log.Logger) error {
	pending := migrations.GetPendingMigrations(db)
	if len(pending) == 0 {
		logger.Info("✅ No pending migrations. Database is already up to date.")
		return nil
	}

	for _, m := range pending {
		logger.Info("[DRY-RUN] Would apply", zap.String("table:", m.ID))
	}
	logger.Info("[DRY-RUN] Migration simulation completed.")
	return nil
}

// fail logs an error and returns it up (no os.Exit)
func fail(logger log.Logger, msg string, err error) error {
	logger.Error(fmt.Sprintf("%s: %v", msg, err))
	return err
}
