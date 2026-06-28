package config

const (
	// Default Postgre configs.
	PGDefaultMaxOpenConns    = 60
	PGDefaultConnmaxLifetime = 120
	PGDefaultMaxIdleConns    = 30
	PGDefaultConnMaxIdleTime = 20
	PGDefaultDbSslMode       = "disable"
	PGDefaultTimezone        = "Asia/Jakarta"
)

// Postgres config
type PostgresConfig struct {
	Host            string `yaml:"DATABASE_HOST" env:"HOST,secret"`
	Port            string `yaml:"DATABASE_PORT" env:"PORT,secret"`
	Username        string `yaml:"DATABASE_USERNAME" env:"USERNAME,secret"`
	Password        string `yaml:"DATABASE_PASSWORD" env:"PASSWORD,secret"`
	Name            string `yaml:"DATABASE_NAME" env:"NAME,secret"`
	SSLMode         string `yaml:"DATABASE_SSL_MODE" env:"SSL_MODE,secret"`
	MaxOpenConns    int    `yaml:"DATABASE_MAX_OPEN_CONN" env:"MAX_OPEN_CONN,secret"`
	ConnMaxLifetime int    `yaml:"DATABASE_MAX_CONN_LIFETIME" env:"MAX_CONN_LIFETIME,secret"`
	MaxIdleConns    int    `yaml:"DATABASE_MAX_IDLE_CONN" env:"MAX_IDLE_CONN,secret"`
	ConnMaxIdleTime int    `yaml:"DATABASE_MAX_IDLETIME_CONN" env:"MAX_IDLETIME_CONN,secret"`
	Timezone        string `yaml:"DATABASE_TIMEZONE" env:"TIMEZONE,secret"`
}

func DatabaseConfigInit() PostgresConfig {
	return PostgresConfig{
		SSLMode:         PGDefaultDbSslMode,
		MaxOpenConns:    PGDefaultMaxOpenConns,
		ConnMaxLifetime: PGDefaultConnmaxLifetime,
		MaxIdleConns:    PGDefaultMaxIdleConns,
		ConnMaxIdleTime: PGDefaultConnMaxIdleTime,
		Timezone:        PGDefaultTimezone,
	}
}
