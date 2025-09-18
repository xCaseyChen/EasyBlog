package config

import "time"

// config
type Config struct {
	PgCfg     PgConfig     // postgres config
	ServerCfg ServerConfig // server config
	PgConnCfg PgConnConfig // postgres connection config
}

// server config
type ServerConfig struct {
	Port      uint16 `mapstructure:"port"` // server port
	JwtSecret string
}

// postgres config
type PgConfig struct {
	Host     string // postgres host
	Port     string // postgres port
	User     string // postgres user
	Password string // postgres password
	Database string // postgres database
}

// postgres connection config
type PgConnConfig struct {
	ConnMaxRetries    int           `mapstructure:"conn_max_retries"`    // postgres connection maximum retries
	ConnRetryInterval time.Duration `mapstructure:"conn_retry_interval"` // postgres connection retry interval
	MaxOpenConns      int           `mapstructure:"max_open_conns"`      // postgres maximum open connections
	MaxIdleConns      int           `mapstructure:"max_idle_conns"`      // postgres maximum idle connections
	ConnMaxLifetime   time.Duration `mapstructure:"conn_max_lifetime"`   // postgres maximum connection lifetime
	ConnMaxIdletime   time.Duration `mapstructure:"conn_max_idletime"`   // postgres maximum idle lifetime
}
