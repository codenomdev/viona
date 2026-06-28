package plugin

import (
	"github.com/codenomdev/viona/pkg/response"
	"github.com/codenomdev/viona/plugins"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetAllPluginStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := make([]*GetAllPluginStatusResp, 0)

		_ = plugins.CallBase(func(base plugins.Base) error {
			info := base.Info()
			resp = append(resp, &GetAllPluginStatusResp{
				SlugName: info.SlugName,
				Enabled:  plugins.StatusManager.IsEnabled(info.SlugName),
			})
			return nil
		})

		success := response.NewHttpOK(resp)

		return c.JSON(response.ParseHttpResponse(success))
	}
}
