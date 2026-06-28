package config

import "embed"

//go:embed  config.yaml
var DefaultConfig embed.FS
