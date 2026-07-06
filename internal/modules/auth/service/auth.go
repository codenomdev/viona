package service

import (
	"context"

	"github.com/codenomdev/viona/internal/modules/auth/dto"
	userDto "github.com/codenomdev/viona/internal/modules/user/dto"
	userSvc "github.com/codenomdev/viona/internal/modules/user/service"
	"github.com/codenomdev/viona/pkg/response"
	"github.com/codenomdev/viona/pkg/util"
	"gorm.io/gorm"
)

type (
	Service interface{}
	service struct {
		db          *gorm.DB
		userService userSvc.Service
	}
)

func NewService(
	db *gorm.DB,
	userService userSvc.Service,
) Service {
	return &service{
		db:          db,
		userService: userService,
	}
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
