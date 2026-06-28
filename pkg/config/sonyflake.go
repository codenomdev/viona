package config

var (
	defaultSonyFlakeStartTimeUnix int64 = 1700000000
	defaultSonyFlakeMachineID     int   = 1
)

type SonyflakeConfig struct {
	START_SERVER_TIME_UNIX int64 `yaml:"START_SERVER_TIME_UNIX" env:"SONYFLAKE_START_TIME_UNIX"`
	MACHINE_ID             int   `yaml:"MACHINE_ID" env:"SONYFLAKE_MACHINE_ID"`
}

func NewSonyFlakeInit() SonyflakeConfig {
	return SonyflakeConfig{
		START_SERVER_TIME_UNIX: defaultSonyFlakeStartTimeUnix,
		MACHINE_ID:             defaultSonyFlakeMachineID,
	}
}
