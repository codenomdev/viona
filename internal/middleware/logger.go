package middleware

import (
	"net/http"
	"time"

	"github.com/codenomdev/viona/pkg/log"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		start := time.Now()

		req := c.Request()

		// wrap manual
		rw := &responseWriter{ResponseWriter: c.Response()}

		// replace response writer
		c.SetResponse(rw)

		logger := log.FromContext(req.Context())

		ctx := log.ToContext(req.Context(), logger)
		c.SetRequest(req.WithContext(ctx))

		err := next(c)

		latency := time.Since(start)

		request := c.Request()

		ip := request.Header.Get("CF-Connecting-IP")
		if ip == "" {
			ip = request.Header.Get("X-Real-IP")
		}
		if ip == "" {
			ip = request.Header.Get("X-Forwarded-For")
		}
		if ip == "" {
			ip = c.RealIP()
		}

		fields := []zap.Field{
			zap.String("request_id", rw.Header().Get(echo.HeaderXRequestID)),
			zap.String("method", req.Method),
			zap.String("path", c.Path()),
			zap.Int("status", rw.Status()),
			zap.Int("size", rw.Size()),
			zap.Duration("latency", latency),
			zap.String("client_ip", ip),
		}

		if err != nil {
			logger.Error("http_request", append(fields, zap.String("error", err.Error()))...)
			return err
		}

		logger.Info("http_request", fields...)
		return nil
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Size() int {
	return w.size
}
