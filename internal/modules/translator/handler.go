package translator

import (
	"encoding/json"

	"github.com/codenomdev/viona/pkg/response"
	"github.com/codenomdev/viona/pkg/translator"
	"github.com/labstack/echo/v5"
	"github.com/segmentfault/pacman/i18n"
)

type Handler struct {
	trans i18n.Translator
}

func NewHandler(
	trans i18n.Translator,
) *Handler {
	return &Handler{
		trans: trans,
	}
}

func (h *Handler) GetLanguage() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data, _ := h.trans.Dump(translator.GetLangByCtx(c.Request().Context()))
		var resp map[string]any
		_ = json.Unmarshal(data, &resp)

		success := response.NewHttpOK(resp)
		return c.JSON(response.ParseHttpResponse(success))
	}
}
