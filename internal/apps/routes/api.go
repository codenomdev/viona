package routes

import (
	"github.com/codenomdev/viona/internal/modules/plugin"
	"github.com/codenomdev/viona/internal/modules/setting"
	"github.com/labstack/echo/v4"
)

type RegisterApiRoutes struct {
	settingRoute *setting.Route
	pluginRoute  *plugin.Route
}

func NewApiRoutes(
	settingRoute *setting.Route,
	pluginRoute *plugin.Route,
) *RegisterApiRoutes {
	return &RegisterApiRoutes{
		settingRoute: settingRoute,
		pluginRoute:  pluginRoute,
	}
}

func (r *RegisterApiRoutes) MapBaseApiRoute(e *echo.Echo) {
	apiRoute := e.Group("/api/v1")
	apiRoute.GET("/health", func(c echo.Context) error {
		return c.NoContent(200)
	})
	// register setting route
	r.settingRoute.RegisterSettingRoute(apiRoute)
	// register plugin route
	r.pluginRoute.RegisterPluginRoute(apiRoute)
}
