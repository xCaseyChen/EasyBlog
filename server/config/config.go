package config

type Configs struct {
	ServerCfg ServerConfig
	PgCfg     PgConfig
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
