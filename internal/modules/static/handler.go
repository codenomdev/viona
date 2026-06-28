package static

import (
	"net/http"
	"strings"

	"github.com/codenomdev/viona/pkg/response"
	"github.com/codenomdev/viona/ui"
	"github.com/labstack/echo/v4"
)

type (
	Handler interface {
		GetFaviconIco() echo.HandlerFunc
		GetManifestJson() echo.HandlerFunc
		SPAHandler() echo.HandlerFunc
	}
	handler struct{}
)

func NewStaticHandler() Handler {
	return &handler{}
}

func (h *handler) GetFaviconIco() echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := ui.Build.ReadFile(
			UIRootFilePath + "/favicon.ico",
		)

		if err != nil {
			return echo.NewHTTPError(
				http.StatusNotFound,
			)
		}

		return c.Blob(
			http.StatusOK,
			"image/vnd.microsoft.icon",
			file,
		)
	}
}

func (h *handler) SPAHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
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

		file, err := ui.Build.ReadFile(
			UIIndexFilePath,
		)

		if err != nil {
			return echo.NewHTTPError(
				http.StatusNotFound,
			)
		}

		content := string(file)

		// optional CDN replacement
		cdnPrefix := ""

		if cdnPrefix != "" {

			cdnPrefix = strings.TrimSuffix(
				cdnPrefix,
				"/",
			)

			content = strings.ReplaceAll(
				content,
				"/static",
				cdnPrefix+"/static",
			)
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
			content,
		)
	}
}

func (h *handler) GetManifestJson() echo.HandlerFunc {
	return func(c echo.Context) error {
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
