package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/codenomdev/viona/pkg/config"
	"github.com/codenomdev/viona/pkg/sonyflake"
	"github.com/codenomdev/viona/pkg/sonyflake/plugin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPgsqlDB(cfg *config.Config, ctx context.Context) (*gorm.DB, func(), error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.POSTGRES.Host,
		cfg.POSTGRES.Username,
		cfg.POSTGRES.Password,
		cfg.POSTGRES.Name,
		cfg.POSTGRES.Port,
		cfg.POSTGRES.SSLMode,
		cfg.POSTGRES.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, nil, err
	}

	config, err := db.DB()

	if err != nil {
		return nil, nil, err
	}

	config.SetMaxOpenConns(cfg.POSTGRES.MaxOpenConns)
	config.SetConnMaxLifetime(time.Duration(cfg.POSTGRES.ConnMaxLifetime) * time.Second)
	config.SetMaxIdleConns(cfg.POSTGRES.MaxIdleConns)
	config.SetConnMaxIdleTime(time.Duration(cfg.POSTGRES.ConnMaxIdleTime) * time.Second)

	if err := config.Ping(); err != nil {
		return nil, nil, err
	}

	sonyflake.InitFromContext(ctx)

	if err := db.Use(plugin.SonyflakePlugin{}); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		config.Close()
	}

	return db, cleanup, nil
}
