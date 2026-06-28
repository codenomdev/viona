package gorm

import (
	"context"

	"gorm.io/gorm"
)

type dbCtxKey struct{}

func ToContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbCtxKey{}, db)
}

func FromContext(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return nil
	}
	if db, ok := ctx.Value(dbCtxKey{}).(*gorm.DB); ok {
		return db
	}
	return nil
}
