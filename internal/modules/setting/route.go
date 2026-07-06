package setting

import "github.com/labstack/echo/v5"

type Route struct {
	h *handler
}

func NewRoute(
	h *handler,
) *Route {
	return &Route{
		h: h,
	}
}

func (r *Route) RegisterSettingRoute(e *echo.Group) {
	group := e.Group("/setting")
	group.GET("/get", r.h.GetSettingAll())
}
