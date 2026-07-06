package plugins

import "github.com/labstack/echo/v5"

type RenderConfig struct {
	SelectTheme string `json:"select_theme"`
}

// select_theme

type Render interface {
	Base
	GetRenderConfig(ctx echo.Context) (renderConfig *RenderConfig)
}

var (
	// CallRender is a function that calls all registered parsers
	CallRender,
	registerRender = MakePlugin[Render](false)
)
