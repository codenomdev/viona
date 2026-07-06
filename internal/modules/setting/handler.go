package setting

import (
	"github.com/codenomdev/viona/internal/modules/setting/service"
	"github.com/codenomdev/viona/pkg/response"
	"github.com/labstack/echo/v5"
)

type (
	handler struct {
		settingService service.Service
	}
)

func NewHandler(
	settingService service.Service,
) *handler {
	return &handler{
		settingService: settingService,
	}
}

func (h *handler) GetSettingAll() echo.HandlerFunc {
	return func(c *echo.Context) error {
		rest, err := h.settingService.GetSettingsPerGroup(c.Request().Context())

		if err != nil {
			return c.JSON(response.ParseHttpResponse(err))
		}

		success := response.NewHttpOK(rest)
		return c.JSON(response.ParseHttpResponse(success))
	}
}
