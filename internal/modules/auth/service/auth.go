package service

import (
	"context"
	"fmt"

	"github.com/codenomdev/viona/internal/modules/auth/dto"
	userDto "github.com/codenomdev/viona/internal/modules/user/dto"
	userSvc "github.com/codenomdev/viona/internal/modules/user/service"
	verifySvc "github.com/codenomdev/viona/internal/modules/user_verify/service"
	"github.com/codenomdev/viona/pkg/response"
	"github.com/codenomdev/viona/pkg/util"
	"github.com/gookit/goutil/strutil"
	"gorm.io/gorm"
)

type (
	Service interface {
		Register(ctx context.Context, req *dto.RegisterWithEmailRequest) error
	}
	service struct {
		db          *gorm.DB
		userService userSvc.Service
		verifySvc   verifySvc.Service
	}
)

func NewService(
	db *gorm.DB,
	userService userSvc.Service,
	verifySvc verifySvc.Service,
) Service {
	return &service{
		db:          db,
		userService: userService,
		verifySvc:   verifySvc,
	}
}

// Register user
func (s *service) Register(ctx context.Context, req *dto.RegisterWithEmailRequest) error {
	tx := s.db.WithContext(ctx).Begin()

	username, _, found := strutil.Cut(req.Email, "@")
	randomStr := strutil.RandWithTpl(4, "0123456789")

	if !found {
		tx.Rollback()
		return response.NewHttpBadRequest([]string{"invalid email address"}, nil)
	}

	user, err := s.userService.CreateNew(ctx, tx, &userDto.RequestCreateNew{
		FullName: req.Fullname,
		Username: fmt.Sprintf("%s%s", username, randomStr),
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := s.verifySvc.CreateNew(ctx, tx, user.ID); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return response.NewHttpUnprocessedEntity([]string{"register failed"})
	}

	return nil
}

// Email Login
func (s *service) EmailLogin(ctx context.Context, req *dto.LoginWithEmailRequest) (*dto.AuthResponse, error) {
	tx := s.db.WithContext(ctx)

	user, err := s.userService.GetByEmail(ctx, tx, &userDto.GetUserByEmailRequest{
		Email: req.Email,
	})

	if err != nil {
		return nil, err
	}

	if err := comparePassword(req.Password, user.Password); err != nil {
		return nil, err
	}

	return nil, nil
}

// Compare password with current password from DB
func comparePassword(plainText string, passHash string) error {
	h := util.NewCryptoHash()

	if valid := h.ValidatePassword(passHash, plainText); !valid {
		return response.NewHttpForbidden([]string{"password is wrong"})
	}

	return nil
}
