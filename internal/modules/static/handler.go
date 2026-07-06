package static

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/codenomdev/viona/pkg/response"
	"github.com/labstack/echo/v5"
)

type (
	Handler interface {
		GetFaviconIco(buildFS fs.FS) echo.HandlerFunc
		GetManifestJson() echo.HandlerFunc
		SPAHandler(buildFS fs.FS) echo.HandlerFunc
	}
	handler struct{}
)

func NewStaticHandler() Handler {
	return &handler{}
}

func (h *handler) GetFaviconIco(buildFS fs.FS) echo.HandlerFunc {
	return func(c *echo.Context) error {
		return c.FileFS(
			"favicon.ico",
			buildFS,
		)
	}
}

func (h *handler) SPAHandler(buildFS fs.FS) echo.HandlerFunc {
	return func(c *echo.Context) error {
		path := c.Request().URL.Path

		if strings.HasPrefix(path, "/api/") {
			return c.JSON(
				response.ParseHttpResponse(
					response.NewHttpNotFound([]string{"not found"}, nil),
				),
			)
		}

		// optional install redirect
		if path == "/install" {
			return c.Redirect(
				http.StatusFound,
				"/",
			)
		}

		cleanPath := strings.TrimPrefix(path, "/")

		if cleanPath != "" {
			if _, err := fs.Stat(buildFS, cleanPath); err == nil {
				return c.FileFS(cleanPath, buildFS)
			}
		}

		file, err := fs.ReadFile(buildFS, "index.html")
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "not found")
		}

		c.Response().Header().Set(
			echo.HeaderContentType,
			"text/html; charset=utf-8",
		)

		c.Response().Header().Set(
			"X-Frame-Options",
			"DENY",
		)

		return c.HTML(
			http.StatusOK,
			string(file),
		)
	}
}

func (h *handler) GetManifestJson() echo.HandlerFunc {
	return func(c *echo.Context) error {
		resp := GetManifestJsonResp{
			ManifestVersion: 3,
			Version:         "1",
			Revision:        "1",
			ShortName:       "codenom",
			Name:            "codenon.com",
			StartUrl:        ".",
			Display:         "standalone",
			ThemeColor:      "#000000",
			BackgroundColor: "#ffffff",
			Icons:           CreateManifestJsonIcons("favicon.ico"),
		}

		return c.JSON(http.StatusOK, resp)
	}
}
