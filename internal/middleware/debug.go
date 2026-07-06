package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/codenomdev/viona/pkg/config"
	"github.com/codenomdev/viona/pkg/log"
	"github.com/labstack/echo/v5"
)

// Debug dump request middleware
func DebugMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		cfg := config.FromContext(c.Request().Context())
		logger := log.FromContext(c.Request().Context())
		if cfg.HOST.SERVER_DEBUG {
			dump, err := httputil.DumpRequest(c.Request(), true)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			logger.Info(fmt.Sprintf("\nRequest dump begin :--------------\n\n%s\n\nRequest dump end :--------------", dump))
		}
		return next(c)
	}
}
