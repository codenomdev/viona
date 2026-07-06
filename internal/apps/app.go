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
	"github.com/labstack/echo/v5"
	"github.com/segmentfault/pacman/i18n"
	"go.uber.org/zap"
)

const (
	maxHeaderBytes = 1 << 20
)

type AppServer struct {
	echo       *echo.Echo
	cfg        *config.Config
	log        log.Logger
	apiRoutes  *routes.RegisterApiRoutes
	uiRoutes   *routes.UIRoutes
	translator i18n.Translator
}

// Server
func NewApp(
	cfg *config.Config,
	logger log.Logger,
	apiRoutes *routes.RegisterApiRoutes,
	uiRoutes *routes.UIRoutes,
	translator i18n.Translator,
) *AppServer {
	return &AppServer{
		echo:       echo.New(),
		cfg:        cfg,
		log:        logger,
		apiRoutes:  apiRoutes,
		uiRoutes:   uiRoutes,
		translator: translator,
	}
}

// server start...
func (s *AppServer) Start() error {
	address := net.JoinHostPort(
		s.cfg.HOST.DOMAIN,
		strconv.Itoa(s.cfg.HOST.PORT),
	)

	sc := echo.StartConfig{
		Address: address,
		BeforeServeFunc: func(h *http.Server) error {
			// server = h
			h.ReadTimeout = time.Duration(s.cfg.HOST.SERVER_READ_TIMEOUT) * time.Second
			h.ReadHeaderTimeout = 5 * time.Second
			h.WriteTimeout = time.Duration(s.cfg.HOST.SERVER_WRITE_TIMEOUT) * time.Second
			h.IdleTimeout = time.Duration(s.cfg.HOST.SERVER_IDLE_TIMEOUT) * time.Second
			h.MaxHeaderBytes = maxHeaderBytes
			return nil
		},
	}

	s.defaultRegisterMiddleware()
	s.apiRoutes.MapBaseApiRoute(s.echo)
	s.uiRoutes.MapBaseStatic(s.echo, s.cfg)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	errChan := make(chan error, 1)

	go func() {
		var err error
		if s.cfg.HOST.SSL.SSL_ENABLE {
			err = sc.StartTLS(
				ctx,
				s.echo,
				s.cfg.HOST.SSL.CERT_FILE,
				s.cfg.HOST.SSL.KEY_FILE,
			)
		} else {
			err = sc.Start(ctx, s.echo)
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("server error", zap.Error(err))
			select {
			case errChan <- err:
			default:
			}
		}
	}()

	select {
	case <-ctx.Done():
		s.log.Info("server stopped")
		return nil

	case err := <-errChan:
		return err
	}
}
