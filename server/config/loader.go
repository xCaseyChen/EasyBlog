package config

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
}

func (c *Configs) LoadConfig() {
	c.ServerCfg.loadConfig()
	c.PgCfg.loadConfig()
}

func (s *ServerConfig) loadConfig() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to load server.yaml: %v", err)
	}
	if err := viper.Unmarshal(s); err != nil {
		log.Fatalf("Failed to unmarshal server config: %v", err)
	}
}

func (p *PgConfig) loadConfig() {
	p.Host = getEnvStringOrFatal("POSTGRES_HOST")
	p.Port = getEnvStringOrFatal("POSTGRES_PORT")
	p.User = getEnvStringOrFatal("POSTGRES_USER")
	p.Password = getEnvStringOrFatal("POSTGRES_PASSWORD")
	p.Database = getEnvStringOrFatal("POSTGRES_DB")
}

func getEnvStringOrFatal(env string) string {
	if !viper.IsSet(env) {
		log.Fatalf("Environment %s not found", env)
	}
	val := viper.GetString(env)
	return val
}
