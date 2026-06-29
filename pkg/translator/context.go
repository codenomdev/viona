package translator

import (
	"context"

	"github.com/codenomdev/viona/internal/middleware"
	"github.com/segmentfault/pacman/i18n"
)

func GetLangByCtx(ctx context.Context) i18n.Language {
	acceptLanguage, ok := ctx.Value(
		middleware.AcceptLanguageContextKey,
	).(i18n.Language)

	if ok {
		return acceptLanguage
	}

	return i18n.DefaultLanguage
}
