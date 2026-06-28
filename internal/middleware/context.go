package middleware

import (
	"context"
	"time"

	"github.com/codenomdev/viona/pkg/config"
	"github.com/codenomdev/viona/pkg/log"
	"github.com/labstack/echo/v4"
)

func ContextInjector(cfg *config.Config, logger log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := config.ToContext(req.Context(), cfg)
			ctx = log.ToContext(ctx, logger)

			timeout := time.Duration(cfg.HOST.SERVER_REQUEST_TIMEOUT) * time.Second
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			c.SetRequest(req.WithContext(ctx))
			return next(c)
		}
	}
}
