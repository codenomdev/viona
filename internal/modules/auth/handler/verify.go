package handler

import (
	"github.com/codenomdev/viona/internal/modules/auth/dto"
	"github.com/codenomdev/viona/internal/modules/auth/service"
	"github.com/codenomdev/viona/pkg/response"
	"github.com/codenomdev/viona/pkg/translator"
	"github.com/codenomdev/viona/pkg/validator"
	"github.com/labstack/echo/v5"
)

type (
	Handler interface {
		RegisterHandler() echo.HandlerFunc
	}
	handler struct {
		authSvc service.Service
	}
)

func NewHandler(
	authSvc service.Service,
) Handler {
	return &handler{
		authSvc: authSvc,
	}
}

func (h *handler) RegisterHandler() echo.HandlerFunc {
	return func(c *echo.Context) error {
		trans := translator.GetLangByCtx(c.Request().Context())
		schema := &dto.RegisterWithEmailRequest{}

		if err := c.Bind(schema); err != nil {
			return c.JSON(
				response.ParseHttpResponse(
					response.NewHttpBadRequest([]string{err.Error()}, nil),
				),
			)
		}

		if err := validator.GetValidatorByLang(trans).Check(schema); err != nil {
			return c.JSON(response.ParseHttpResponse(err))
		}

		if err := h.authSvc.Register(c.Request().Context(), schema); err != nil {
			return c.JSON(response.ParseHttpResponse(err))
		}

		return c.JSON(response.ParseHttpResponse(response.NewHttpOK(map[string]string{
			"message": "OK",
		})))
	}
}
