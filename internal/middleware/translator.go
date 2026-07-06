package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/segmentfault/pacman/i18n"
)

type contextKey string

const AcceptLanguageContextKey contextKey = "accept-language"

func I18nMiddleware(defaultLang i18n.Language) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			lang := c.QueryParam("lang")

			if lang == "" {
				lang = c.Request().Header.Get("Accept-Language")

				if idx := strings.Index(lang, ","); idx != -1 {
					lang = lang[:idx]
				}
			}

			if lang == "" {
				lang = string(defaultLang)
			}

			reqCtx := context.WithValue(
				c.Request().Context(),
				AcceptLanguageContextKey,
				i18n.Language(lang),
			)

			c.SetRequest(c.Request().WithContext(reqCtx))

			return next(c)
		}
	}
}
