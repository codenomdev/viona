package plugin

import (
	"context"

	"github.com/codenomdev/viona/internal/modules/plugin/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		GetAll(ctx context.Context, tx *gorm.DB) ([]domain.PluginConfig, error)
		GetWhere(ctx context.Context, tx *gorm.DB, query any, args ...any) gorm.ChainInterface[domain.PluginConfig]
	}
	repository struct {
		db *gorm.DB
	}
)

func NewRepository(
	db *gorm.DB,
) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context, tx *gorm.DB) ([]domain.PluginConfig, error) {
	return gorm.G[domain.PluginConfig](r.getDB(tx)).Find(ctx)
}

// Get where
func (r *repository) GetWhere(ctx context.Context, tx *gorm.DB, query any, args ...any) gorm.ChainInterface[domain.PluginConfig] {
	return gorm.G[domain.PluginConfig](r.getDB(tx)).Scopes(func(db *gorm.Statement) {
		db.Scopes(scopeAllRows)
	}).Where(query, args...)
}

func (r *repository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *repository) baseQuery(tx *gorm.DB) *gorm.DB {
	return r.getDB(tx).Scopes(scopeAllRows)
}

// scope all rows
func scopeAllRows(db *gorm.DB) *gorm.DB {
	return db.Select("id, slug_name, value")
}
