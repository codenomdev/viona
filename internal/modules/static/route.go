package static

import (
	"io/fs"

	"github.com/labstack/echo/v5"
)

type Routes struct {
	h Handler
}

func NewStaticRoute(
	h Handler,
) *Routes {
	return &Routes{
		h: h,
	}
}

func (r *Routes) RegisterStatic(e *echo.Echo, buildFS fs.FS) {
	e.GET("/favicon.ico", r.h.GetFaviconIco(buildFS))
	e.GET("/manifest.json", r.h.GetManifestJson())
	e.GET("/*", r.h.SPAHandler(buildFS))
}
