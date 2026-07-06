package plugins

var (
	DefaultCDNFileType = map[string]bool{
		".ico":   true,
		".json":  true,
		".css":   true,
		".js":    true,
		".webp":  true,
		".woff2": true,
		".woff":  true,
		".jpg":   true,
		".svg":   true,
		".png":   true,
		".map":   true,
		".txt":   true,
	}
)

type CDN interface {
	Base
	GetStaticPrefix() string
}

var (
	// CallCDN is a function that calls all registered parsers
	CallCDN,
	registerCDN = MakePlugin[CDN](false)
)

func coordinatedCDNPlugins(slugName string) (enabledSlugNames []string) {
	isCDN := false
	_ = CallCDN(func(cdn CDN) error {
		name := cdn.Info().SlugName
		if slugName == name {
			isCDN = true
		} else {
			enabledSlugNames = append(enabledSlugNames, name)
		}
		return nil
	})
	if isCDN {
		return enabledSlugNames
	}
	return nil
}
