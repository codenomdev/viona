package service

import (
	"context"

	"github.com/codenomdev/viona/internal/modules/user_verify/domain"
	"github.com/codenomdev/viona/internal/modules/user_verify/repository"
	"github.com/codenomdev/viona/pkg/response"
	"gorm.io/gorm"
)

type (
	Service interface {
		CreateNew(ctx context.Context, tx *gorm.DB, userID int64) error
	}
	service struct {
		db   *gorm.DB
		repo repository.Repository
	}
)

func NewService(
	db *gorm.DB,
	repo repository.Repository,
) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateNew(ctx context.Context, tx *gorm.DB, userID int64) error {
	if tx == nil {
		tx = s.db.WithContext(ctx)
	}

	if err := s.repo.Create(ctx, tx, &domain.EmailVerification{
		UserID: userID,
	}); err != nil {
		return response.NewHttpUnprocessedEntity([]string{"cannot unprocessed"})
	}
	return nil
}
