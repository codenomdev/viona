package plugins

import "github.com/labstack/echo/v5"

type EmbedConfig struct {
	Platform string `json:"platform"`
	Enable   bool   `json:"enable"`
}

type Embed interface {
	Base
	GetEmbedConfigs(ctx echo.Context) (embedConfigs []*EmbedConfig, err error)
}

var (
	// CallEmbed is a function that calls all registered parsers
	CallEmbed,
	registerEmbed = MakePlugin[Embed](false)
)
