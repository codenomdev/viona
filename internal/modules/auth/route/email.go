package route

import (
	"github.com/codenomdev/viona/internal/modules/auth/handler"
	"github.com/labstack/echo/v5"
)

type Route struct {
	h handler.Handler
}

func NewAuthRoute(
	h handler.Handler,
) *Route {
	return &Route{
		h: h,
	}
}

func (r *Route) RegisterAuthRoute(e *echo.Group) {
	group := e.Group("/auth")
	group.POST("/register", r.h.RegisterHandler())
}
