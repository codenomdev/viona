package static

import (
	"fmt"
	"path/filepath"
	"strings"
)

type GetManifestJsonResp struct {
	ManifestVersion int                `json:"manifest_version"`
	Version         string             `json:"version"`
	Revision        string             `json:"revision"`
	ShortName       string             `json:"short_name"`
	Name            string             `json:"name"`
	Icons           []ManifestJsonIcon `json:"icons"`
	StartUrl        string             `json:"start_url"`
	Display         string             `json:"display"`
	ThemeColor      string             `json:"theme_color"`
	BackgroundColor string             `json:"background_color"`
}
type ManifestJsonIcon struct {
	Src   string `json:"src"`
	Sizes string `json:"sizes"`
	Type  string `json:"type"`
}

func CreateManifestJsonIcons(icon string) []ManifestJsonIcon {
	ext := filepath.Ext(icon)
	if ext == "" {
		ext = "png"
	} else {
		ext = strings.ToLower(ext[1:])
	}
	iconType := fmt.Sprintf("image/%s", ext)
	return []ManifestJsonIcon{
		{
			Src:   icon,
			Sizes: "16x16",
			Type:  iconType,
		},
		{
			Src:   icon,
			Sizes: "32x32",
			Type:  iconType,
		},
		{
			Src:   icon,
			Sizes: "48x48",
			Type:  iconType,
		},
		{
			Src:   icon,
			Sizes: "128x128",
			Type:  iconType,
		},
	}
}
