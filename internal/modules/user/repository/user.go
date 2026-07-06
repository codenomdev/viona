package repository

import (
	"context"
	"errors"

	"github.com/codenomdev/viona/internal/modules/user/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(ctx context.Context, tx *gorm.DB, model *domain.User) error
		GetByEmail(ctx context.Context, tx *gorm.DB, email string) (*domain.User, error)
		GetByUsername(ctx context.Context, tx *gorm.DB, username string) (*domain.User, error)
		CountByUsername(ctx context.Context, tx *gorm.DB, username string) (totalRows int64)
		CountByEmail(ctx context.Context, tx *gorm.DB, email string) (totalRows int64)
		UpdateByID(ctx context.Context, tx *gorm.DB, id int64, model domain.User) error
		GetWhere(ctx context.Context, tx *gorm.DB, query any, args ...any) gorm.ChainInterface[domain.User]
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

// Create new user
func (r *repository) Create(ctx context.Context, tx *gorm.DB, model *domain.User) error {
	return gorm.G[domain.User](r.getDB(tx)).Create(ctx, model)
}

// Get by email
func (r *repository) GetByEmail(ctx context.Context, tx *gorm.DB, email string) (*domain.User, error) {
	u, err := r.GetWhere(ctx, tx, "email = ?", email).Take(ctx)
	return &u, err
}

// Get by username
func (r *repository) GetByUsername(ctx context.Context, tx *gorm.DB, username string) (*domain.User, error) {
	u, err := r.GetWhere(ctx, tx, "username = ?", username).Take(ctx)
	return &u, err
}

// Count by Username
func (r *repository) CountByUsername(ctx context.Context, tx *gorm.DB, username string) (totalRows int64) {
	r.baseQuery(tx).
		WithContext(ctx).
		Where("username = ?", username).
		Model(&domain.User{}).
		Distinct("id").
		Count(&totalRows)
	return totalRows
}

// Count by Email
func (r *repository) CountByEmail(ctx context.Context, tx *gorm.DB, email string) (totalRows int64) {
	r.baseQuery(tx).
		WithContext(ctx).
		Where("email = ?", email).
		Model(&domain.User{}).
		Distinct("id").
		Count(&totalRows)
	return totalRows
}

// Update user by ID
func (r *repository) UpdateByID(ctx context.Context, tx *gorm.DB, id int64, model domain.User) error {
	updated, err := gorm.G[domain.User](r.getDB(tx)).Table(domain.TableName).Where("id = ?", id).Updates(ctx, model)

	if err != nil {
		return err
	}

	if updated == 0 {
		return errors.New("no updates")
	}

	return nil
}

// Get where
func (r *repository) GetWhere(ctx context.Context, tx *gorm.DB, query any, args ...any) gorm.ChainInterface[domain.User] {
	return gorm.G[domain.User](r.getDB(tx)).Scopes(func(db *gorm.Statement) {
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
	return db.Select("id, fullname, username, email, password, avatar, is_active, email_verified_at, last_login_date, deleted_at, created_at, updated_at")
}
