package middleware

import (
	"time"

	"github.com/codenomdev/viona/pkg/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		req := c.Request()
		res := c.Response()

		logger := log.FromContext(req.Context())

		// inject logger ke context (biar downstream bisa pakai)
		ctx := log.ToContext(req.Context(), logger)
		c.SetRequest(req.WithContext(ctx))

		err := next(c)

		stop := time.Now()
		latency := stop.Sub(start)

		fields := []zap.Field{
			zap.String("request_id", res.Header().Get(echo.HeaderXRequestID)),
			zap.String("method", req.Method),
			zap.String("path", req.URL.Path),
			zap.Int("status", res.Status),
			zap.Duration("latency", latency),
			zap.Int64("size", res.Size),
			zap.String("client_ip", c.RealIP()),
		}

		if err != nil {
			if he, ok := err.(*echo.HTTPError); ok {

				fields = append(fields,
					zap.String("error", he.Error()),
				)

				switch {
				case he.Code >= 500:
					logger.Error("http_request", fields...)

				case he.Code >= 400:
					logger.Warn("http_request", fields...)

				default:
					logger.Info("http_request", fields...)
				}

				return err
			}

			fields = append(fields,
				zap.String("error", err.Error()),
			)

			logger.Error("http_request", fields...)
			return err
		}

		logger.Info("http_request", fields...)
		return nil
	}
}
