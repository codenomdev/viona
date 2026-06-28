package apps

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/codenomdev/viona/internal/apps/routes"
	"github.com/codenomdev/viona/pkg/config"
	"github.com/codenomdev/viona/pkg/log"
	"github.com/labstack/echo/v4"
)

const (
	maxHeaderBytes = 1 << 20
)

type AppServer struct {
	echo      *echo.Echo
	cfg       *config.Config
	log       log.Logger
	apiRoutes *routes.RegisterApiRoutes
	uiRoutes  *routes.UIRoutes
}

// Server
func NewApp(
	cfg *config.Config,
	logger log.Logger,
	apiRoutes *routes.RegisterApiRoutes,
	uiRoutes *routes.UIRoutes,
) *AppServer {
	return &AppServer{
		echo:      echo.New(),
		cfg:       cfg,
		log:       logger,
		apiRoutes: apiRoutes,
		uiRoutes:  uiRoutes,
	}
}

// server start...
func (s *AppServer) Start() error {
	address := net.JoinHostPort(
		s.cfg.HOST.DOMAIN,
		strconv.Itoa(s.cfg.HOST.PORT),
	)

	s.echo.Server = &http.Server{
		Addr:              address,
		ReadTimeout:       time.Duration(s.cfg.HOST.SERVER_READ_TIMEOUT) * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      time.Duration(s.cfg.HOST.SERVER_WRITE_TIMEOUT) * time.Second,
		IdleTimeout:       time.Duration(s.cfg.HOST.SERVER_IDLE_TIMEOUT) * time.Second,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	s.defaultRegisterMiddleware()
	s.apiRoutes.MapBaseApiRoute(s.echo)
	s.uiRoutes.MapBaseStatic(s.echo, s.cfg)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	errChan := make(chan error, 1)

	go func() {
		var err error

		if s.cfg.HOST.SSL.SSL_ENABLE {
			err = s.echo.StartTLS(
				address,
				s.cfg.HOST.SSL.CERT_FILE,
				s.cfg.HOST.SSL.KEY_FILE,
			)
		} else {
			err = s.echo.Start(address)
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.log.Info("Shutting down server")

	case err := <-errChan:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	return s.echo.Server.Shutdown(shutdownCtx)
}
