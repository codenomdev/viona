package static

import "github.com/labstack/echo/v4"

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

func (r *Routes) RegisterStatic(e *echo.Echo) {
	e.GET("/favicon.ico", r.h.GetFaviconIco())
	e.GET("/manifest.json", r.h.GetManifestJson())
	e.GET("/*", r.h.SPAHandler())
}
