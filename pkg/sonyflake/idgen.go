package sonyflake

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/codenomdev/viona/pkg/config"
	"github.com/codenomdev/viona/pkg/log"
	"github.com/sony/sonyflake/v2"
	"go.uber.org/zap"
)

var (
	flake *sonyflake.Sonyflake
	once  sync.Once
)

func InitFromContext(ctx context.Context) {
	once.Do(func() {
		cfg := config.FromContext(ctx)
		logger := log.FromContext(ctx)
		st := sonyflake.Settings{
			StartTime: time.Unix(cfg.SONYFLAKE.START_SERVER_TIME_UNIX, 0),
			MachineID: func() (int, error) {
				return cfg.SONYFLAKE.MACHINE_ID, nil
			},
		}
		var err error
		flake, err = sonyflake.New(st)
		if err != nil {
			logger.Error("failed to init sonyflake", zap.Error(err))
		}
	})
}

// Generate return unique ID Sonyflake.
// Panic if your not set.
func Generate(ctx context.Context) uint64 {
	logger := log.FromContext(ctx)
	if flake == nil {
		InitFromContext(ctx)
		if flake == nil {
			logger.Error("sonyflake not initialized — call sonyflake.InitFromContext() first")
			panic("sonyflake not initialized")
		}
	}
	id, err := flake.NextID()
	if err != nil {
		logger.Error("failed to generate sonyflake ID:", zap.Error(err))
	}
	return uint64(id)
}

// GenerateString return ID unique type string base36 (most short and URL-safe).
func GenerateString(ctx context.Context) string {
	id := Generate(ctx)
	return uint64ToBase36(id)
}

// GenerateWithPrefix return ID unique within prefix specified, example:
//
//	user-gx3ba27v0
func GenerateWithPrefix(ctx context.Context, prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, GenerateString(ctx))
}

const base36Charset = "0123456789abcdefghijklmnopqrstuvwxyz"

// uint64ToBase36 change number uint64 to string base36.
func uint64ToBase36(num uint64) string {
	if num == 0 {
		return "0"
	}
	n := new(big.Int).SetUint64(num)
	base := big.NewInt(36)
	result := ""
	mod := new(big.Int)

	for n.Cmp(big.NewInt(0)) > 0 {
		n, mod = new(big.Int).DivMod(n, base, mod)
		result = string(base36Charset[mod.Int64()]) + result
	}
	return result
}
