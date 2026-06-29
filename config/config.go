package config

import "embed"

//go:embed  *.yaml
var DefaultConfig embed.FS
