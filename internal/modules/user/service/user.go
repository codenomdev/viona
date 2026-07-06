package service

import (
	"context"
	"errors"
	"strings"

	"github.com/codenomdev/viona/internal/modules/user/domain"
	"github.com/codenomdev/viona/internal/modules/user/dto"
	"github.com/codenomdev/viona/internal/modules/user/repository"
	"github.com/codenomdev/viona/pkg/response"
	"gorm.io/gorm"
)

type (
	Service interface {
		GetByEmail(ctx context.Context, tx *gorm.DB, req *dto.GetUserByEmailRequest) (*dto.ResponseUser, error)
		CreateNew(ctx context.Context, tx *gorm.DB, req *dto.RequestCreateNew) (*dto.ResponseUser, error)
	}
	service struct {
		db       *gorm.DB
		userRepo repository.Repository
	}
)

func NewService(
	db *gorm.DB,
	userRepo repository.Repository,
) Service {
	return &service{
		db:       db,
		userRepo: userRepo,
	}
}

// insert new user
func (s *service) CreateNew(ctx context.Context, tx *gorm.DB, req *dto.RequestCreateNew) (*dto.ResponseUser, error) {
	if tx == nil {
		tx = s.db.WithContext(ctx)
	}

	userModel := &domain.User{
		Email:    req.Email,
		Username: req.Username,
		FullName: req.FullName,
		Password: req.Password,
	}

	if err := s.userRepo.Create(ctx, tx, userModel); err != nil {
		errMsg := strings.ToLower(err.Error())

		switch {
		case strings.Contains(errMsg, "email"):
			return nil, response.NewHttpConflict([]string{
				"email already exist",
			})

		case strings.Contains(errMsg, "username"):
			return nil, response.NewHttpConflict([]string{
				"username already exist",
			})

		default:
			return nil, response.NewHttpConflict([]string{
				"create new user failed, please try again later",
			})
		}
	}

	return buildResponse(userModel), nil
}

func (s *service) GetByEmail(ctx context.Context, tx *gorm.DB, req *dto.GetUserByEmailRequest) (*dto.ResponseUser, error) {
	if tx == nil {
		tx = s.db.WithContext(ctx)
	}
	user, err := s.userRepo.GetByEmail(ctx, tx, req.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewHttpNotFound([]string{"user not found"}, nil)
		}
		return nil, response.NewHttpUnprocessedEntity([]string{"unprocessed"})
	}

	return buildResponse(user), nil
}

func buildResponse(model *domain.User) *dto.ResponseUser {
	return &dto.ResponseUser{
		ID:              model.ID,
		FullName:        model.FullName,
		Username:        model.Username,
		Email:           model.Email,
		Password:        model.Password,
		Avatar:          model.Avatar,
		EmailVerifiedAt: model.EmailVerifiedAt,
		IsActive:        int(model.IsActive),
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}
