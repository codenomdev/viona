package repository

import (
	"context"
	"errors"

	"github.com/codenomdev/viona/internal/modules/user_verify/domain"
	"gorm.io/gorm"
)

var (
	ErrNoUpdates error = errors.New("no updates")
)

type (
	Repository interface {
		GetWhere(ctx context.Context, tx *gorm.DB, query any, args ...any) gorm.ChainInterface[domain.EmailVerification]
		Create(ctx context.Context, tx *gorm.DB, model *domain.EmailVerification) error
		GetByToken(ctx context.Context, tx *gorm.DB, tokenHash string) (*domain.EmailVerification, error)
		UpdateByUserID(ctx context.Context, tx *gorm.DB, userID int64, model domain.EmailVerification) error
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

// Create new
func (r *repository) Create(ctx context.Context, tx *gorm.DB, model *domain.EmailVerification) error {
	return gorm.G[domain.EmailVerification](r.getDB(tx)).Create(ctx, model)
}

// Update by user ID
func (r *repository) UpdateByUserID(ctx context.Context, tx *gorm.DB, userID int64, model domain.EmailVerification) error {
	updated, err := gorm.G[domain.EmailVerification](r.getDB(tx)).Table(domain.TableName).Where("user_id = ?", userID).Updates(ctx, model)

	if err != nil || updated == 0 {
		return ErrNoUpdates
	}

	return nil
}

// Get email verify by token
func (r *repository) GetByToken(ctx context.Context, tx *gorm.DB, tokenHash string) (*domain.EmailVerification, error) {
	result, err := r.GetWhere(ctx, tx, "token_hash = ?", tokenHash).Where("verified_at IS NULL").Take(ctx)
	return &result, err
}

// Get where
func (r *repository) GetWhere(ctx context.Context, tx *gorm.DB, query any, args ...any) gorm.ChainInterface[domain.EmailVerification] {
	return gorm.G[domain.EmailVerification](r.getDB(tx)).Scopes(func(db *gorm.Statement) {
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
	return db.Select("id, user_id, token_hash, expires_at, verified_at, ip_address, user_agent, created_at")
}
