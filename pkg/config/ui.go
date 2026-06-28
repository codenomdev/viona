package config

var (
	defaultPublicUrl  string = "/"
	defaultAPIUrl     string = "/"
	defaultBaseUrl    string = ""
	defaultApiBaseUrl string = ""
)

type UIConfig struct {
	PUBLIC_URL   string `yaml:"PUBLIC_URL" env:"UI_PUBLIC_URL"`
	API_URL      string `yaml:"API_URL" env:"UI_API_URL"`
	BASE_URL     string `yaml:"BASE_URL" env:"UI_BASE_URL"`
	API_BASE_URL string `yaml:"API_BASE_URL" env:"API_BASE_URL"`
}

func UIConfigInit() UIConfig {
	return UIConfig{
		PUBLIC_URL:   defaultPublicUrl,
		API_URL:      defaultAPIUrl,
		BASE_URL:     defaultBaseUrl,
		API_BASE_URL: defaultAPIUrl,
	}
}
