package plugin

import "github.com/labstack/echo/v5"

type Route struct {
	h *Handler
}

func NewRoute(
	h *Handler,
) *Route {
	return &Route{
		h: h,
	}
}

func (r *Route) RegisterPluginRoute(e *echo.Group) {
	group := e.Group("/plugin")
	group.GET("/status", r.h.GetAllPluginStatus())
}
