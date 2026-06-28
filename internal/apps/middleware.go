package apps

import (
	"time"

	mainMiddleware "github.com/codenomdev/viona/internal/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4/middleware"
)

func (s *AppServer) defaultRegisterMiddleware() {
	s.echo.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.NewString()
		},
	}))

	s.echo.Use(mainMiddleware.ContextInjector(s.cfg, s.log))

	// Register middleware request logger
	s.echo.Use(mainMiddleware.RequestLoggerMiddleware)

	s.echo.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))

	s.echo.Use(middleware.Secure())

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     s.cfg.HOST.CORSConfig.ALLOW_ORIGINS,
		AllowHeaders:     s.cfg.HOST.CORSConfig.ALLOW_HEADERS,
		AllowCredentials: s.cfg.HOST.CORSConfig.WITH_CREDENTIALS,
		AllowMethods:     s.cfg.HOST.CORSConfig.ALLOW_METHODS,
	}))

	s.echo.Use(middleware.BodyLimit("2M"))

	// When SERVER_DEBUG set to true, we will dump request via middleware
	s.echo.Use(mainMiddleware.DebugMiddleware)
}
