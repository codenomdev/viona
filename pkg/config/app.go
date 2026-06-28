package config

var (
	defaultCorsAllowOrigins = []string{"*"}
	defaultCorsAllowHeaders = []string{"Origin", "Content-Type", "Accept", "X-Request-Id", "Retry-After", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"}
	defaultCorsAllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
)

const (
	defaultDomain               string = "localhost"
	defaultServerPort           int    = 8000
	defaultContextTimeout       int    = 60
	defaultServerDebug          bool   = false
	defaultWriteTimeout         int    = 5
	defaultReadTimeout          int    = 5
	defaultSSLEnable            bool   = false
	defaultCertFile             string = ""
	defaultKeyFile              string = ""
	defaultCORSWithCredentials  bool   = false
	defaultIdleTimeout          int    = 60
	defaultServerRequestTimeout int    = 10
)

type (
	// Host SSL config
	HostSSLConfig struct {
		SSL_ENABLE bool   `yaml:"ENABLE" env:"SSL_ENABLE"`
		CERT_FILE  string `yaml:"CERT_FILE" env:"CERT_FILE"`
		KEY_FILE   string `yaml:"KEY_FILE" env:"KEY_FILE"`
	}
	// Host CORS Config
	HostCORSConfig struct {
		WITH_CREDENTIALS bool     `yaml:"WITH_CREDENTIALS" env:"CORS_WITH_CREDENTIALS"`
		ALLOW_ORIGINS    []string `yaml:"ALLOW_ORIGINS" env:"CORS_ALLOW_ORIGINS"`
		ALLOW_HEADERS    []string `yaml:"ALLOW_HEADERS" env:"CORS_ALLOW_HEADERS"`
		ALLOW_METHODS    []string `yaml:"ALLOW_METHODS" env:"CORS_ALLOW_METHODS"`
	}

	// Host config
	HostConfig struct {
		// Domain name for host
		DOMAIN string `yaml:"DOMAIN" env:"DOMAIN"`
		// the server port. Defaults to 8080
		PORT                   int            `yaml:"PORT" env:"PORT"`
		CONTEXT_TIMEOUT        int            `yaml:"CONTEXT_TIMEOUT" env:"CONTEXT_TIMEOUT"`
		SERVER_DEBUG           bool           `yaml:"SERVER_DEBUG" env:"SERVER_DEBUG"`
		SERVER_WRITE_TIMEOUT   int            `yaml:"SERVER_WRITE_TIMEOUT" env:"SERVER_WRITE_TIMEOUT"`
		SERVER_READ_TIMEOUT    int            `yaml:"SERVER_READ_TIMEOUT" env:"SERVER_READ_TIMEOUT"`
		SSL                    HostSSLConfig  `yaml:"SSL"`
		CORSConfig             HostCORSConfig `yaml:"CORS"`
		SERVER_IDLE_TIMEOUT    int            `yaml:"SERVER_IDLE_TIMEOUT" env:"SERVER_IDLE_TIMEOUT"`
		SERVER_REQUEST_TIMEOUT int            `yaml:"SERVER_REQUEST_TIMEOUT" env:"SERVER_REQUEST_TIMEOUT"`
	}
)

func HostConfigInit() HostConfig {
	return HostConfig{
		DOMAIN:                 defaultDomain,
		PORT:                   defaultServerPort,
		CONTEXT_TIMEOUT:        defaultContextTimeout,
		SERVER_DEBUG:           defaultServerDebug,
		SERVER_WRITE_TIMEOUT:   defaultWriteTimeout,
		SERVER_READ_TIMEOUT:    defaultReadTimeout,
		SERVER_IDLE_TIMEOUT:    defaultIdleTimeout,
		SERVER_REQUEST_TIMEOUT: defaultServerRequestTimeout,
		SSL: HostSSLConfig{
			SSL_ENABLE: defaultSSLEnable,
			CERT_FILE:  defaultCertFile,
			KEY_FILE:   defaultKeyFile,
		},
		CORSConfig: HostCORSConfig{
			ALLOW_ORIGINS:    defaultCorsAllowOrigins,
			ALLOW_HEADERS:    defaultCorsAllowHeaders,
			WITH_CREDENTIALS: defaultCORSWithCredentials,
			ALLOW_METHODS:    defaultCorsAllowMethods,
		},
	}
}
