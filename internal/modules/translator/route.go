package translator

import "github.com/labstack/echo/v5"

type Route struct {
	h *Handler
}

func NewRoute(h *Handler) *Route {
	return &Route{
		h: h,
	}
}

func (r *Route) RegisterTransRoute(e *echo.Group) {
	group := e.Group("/language")
	group.GET("/get", r.h.GetLanguage())
	group.GET("/options", r.h.GetLanguageOptions())
}
