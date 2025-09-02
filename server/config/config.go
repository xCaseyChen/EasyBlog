package config

import "time"

type Config struct {
	PgCfg     PgConfig
	ServerCfg ServerConfig
	PgConnCfg PgConnConfig
}

type ServerConfig struct {
	Port uint16 `mapstructure:"port"`
}

type PgConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type PgConnConfig struct {
	ConnMaxRetries    int           `mapstructure:"conn_max_retries"`
	ConnRetryInterval time.Duration `mapstructure:"conn_retry_interval"`
	MaxOpenConns      int           `mapstructure:"max_open_conns"`
	MaxIdleConns      int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime   time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdletime   time.Duration `mapstructure:"conn_max_idletime"`
}
