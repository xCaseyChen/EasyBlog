package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// load environment and config file
func init() {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
}

// new config
func New() *Config {
	return &Config{}
}

// load config
func (c *Config) LoadConfig() error {
	if err := c.PgCfg.loadConfig(); err != nil {
		return err
	}
	if err := c.ServerCfg.loadConfig(); err != nil {
		return err
	}
	if err := c.PgConnCfg.loadConfig(); err != nil {
		return err
	}
	return nil
}

// load postgres config from environment
func (p *PgConfig) loadConfig() error {
	missing := []string{}
	require := []string{
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
	}
	for _, env := range require {
		if !viper.IsSet(env) {
			missing = append(missing, env)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("required environment not found %v", missing)
	}
	p.Host = viper.GetString("POSTGRES_HOST")
	p.Port = viper.GetString("POSTGRES_PORT")
	p.User = viper.GetString("POSTGRES_USER")
	p.Password = viper.GetString("POSTGRES_PASSWORD")
	p.Database = viper.GetString("POSTGRES_DB")
	return nil
}

// load server config from file
func (s *ServerConfig) loadConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to load config.yaml: %w", err)
	}
	if err := viper.UnmarshalKey("server", s); err != nil {
		return fmt.Errorf("failed to unmarshal server config: %w", err)
	}
	if !viper.IsSet("JWT_SECRET") {
		return fmt.Errorf("required environment not found: JWT_SECRET")
	}
	s.JwtSecret = viper.GetString("JWT_SECRET")
	return nil
}

// load postgres connection config from file
func (p *PgConnConfig) loadConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to load config.yaml: %w", err)
	}
	if err := viper.UnmarshalKey("pg_conn", p); err != nil {
		return fmt.Errorf("failed to unmarshal postgres connection config: %w", err)
	}
	return nil
}
