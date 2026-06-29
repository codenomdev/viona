package config

const (
	defaultI18nDir string = "i18n"
)

type I18NConfig struct {
	BundleDir string `yaml:"BUNDLE_DIR" env:"I18N_BUNDLE_DIR"`
}

func I18NConfigInit() *I18NConfig {
	return &I18NConfig{
		BundleDir: defaultI18nDir,
	}
}
