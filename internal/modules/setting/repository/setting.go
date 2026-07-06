package repository

import (
	"context"
	"errors"

	"github.com/codenomdev/viona/internal/modules/setting/domain"
	"gorm.io/gorm"
)

var (
	ErrNoUpdates error = errors.New("no updates")
)

type (
	Repository interface {
		Create(ctx context.Context, tx *gorm.DB, model *domain.SiteSetting) error
		UpdateByKey(ctx context.Context, tx *gorm.DB, model domain.SiteSetting) error
		GetAll(ctx context.Context, tx *gorm.DB) ([]domain.SiteSetting, error)
		GetByKey(ctx context.Context, tx *gorm.DB, key string) (*domain.SiteSetting, error)
		GetByGroup(ctx context.Context, tx *gorm.DB, group string) ([]domain.SiteSetting, error)
		GetWhere(tx *gorm.DB, query any, args ...any) gorm.ChainInterface[domain.SiteSetting]
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

// create
func (r *repository) Create(ctx context.Context, tx *gorm.DB, model *domain.SiteSetting) error {
	return gorm.G[domain.SiteSetting](r.getDB(tx)).Create(ctx, model)
}

// Update by key setting
func (r *repository) UpdateByKey(ctx context.Context, tx *gorm.DB, model domain.SiteSetting) error {
	updated, err := gorm.G[domain.SiteSetting](r.getDB(tx)).Table(domain.TableName).Where("key = ?", model.Key).Updates(ctx, model)

	if err != nil || updated == 0 {
		return ErrNoUpdates
	}

	return nil
}

// Get by key setting
// Only with status == 1
func (r *repository) GetByKey(ctx context.Context, tx *gorm.DB, key string) (*domain.SiteSetting, error) {
	result, err := r.GetWhere(tx, "key = ?", key).Where("status = 1").Take(ctx)
	return &result, err
}

// Get by group settings
// Only with status == 1
func (r *repository) GetByGroup(ctx context.Context, tx *gorm.DB, group string) ([]domain.SiteSetting, error) {
	return r.GetWhere(tx, "group_name = ?", group).Where("status = 1").Find(ctx)
}

// Get all
func (r *repository) GetAll(ctx context.Context, tx *gorm.DB) ([]domain.SiteSetting, error) {
	return r.GetWhere(tx, "status = 1").Find(ctx)
}

// Get where
func (r *repository) GetWhere(tx *gorm.DB, query any, args ...any) gorm.ChainInterface[domain.SiteSetting] {
	return gorm.G[domain.SiteSetting](r.baseQuery(tx)).Where(query, args...)

}

// Get DB without default scopes
func (r *repository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

// Get DB with default scopes all rows
func (r *repository) baseQuery(tx *gorm.DB) *gorm.DB {
	return r.getDB(tx).Scopes(scopeAllRows)
}

// scope all rows
func scopeAllRows(db *gorm.DB) *gorm.DB {
	return db.Select(
		"id, group_name, key, values, sort, status, created_at, updated_at",
	)
}
