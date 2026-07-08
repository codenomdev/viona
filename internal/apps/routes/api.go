package routes

import (
	authRoute "github.com/codenomdev/viona/internal/modules/auth/route"
	"github.com/codenomdev/viona/internal/modules/plugin"
	"github.com/codenomdev/viona/internal/modules/setting"
	"github.com/codenomdev/viona/internal/modules/translator"
	"github.com/labstack/echo/v5"
)

type RegisterApiRoutes struct {
	settingRoute *setting.Route
	pluginRoute  *plugin.Route
	transRoute   *translator.Route
	authRoute    *authRoute.Route
}

func NewApiRoutes(
	authRoute *authRoute.Route,
	settingRoute *setting.Route,
	pluginRoute *plugin.Route,
	transRoute *translator.Route,
) *RegisterApiRoutes {
	return &RegisterApiRoutes{
		authRoute:    authRoute,
		settingRoute: settingRoute,
		pluginRoute:  pluginRoute,
		transRoute:   transRoute,
	}
}

func (r *RegisterApiRoutes) MapBaseApiRoute(e *echo.Echo) {
	apiRoute := e.Group("/api/v1")
	apiRoute.GET("/health", func(c *echo.Context) error {
		return c.NoContent(200)
	})
	// register auth route
	r.authRoute.RegisterAuthRoute(apiRoute)
	// register setting route
	r.settingRoute.RegisterSettingRoute(apiRoute)
	// register plugin route
	r.pluginRoute.RegisterPluginRoute(apiRoute)
	// language route
	r.transRoute.RegisterTransRoute(apiRoute)
}
