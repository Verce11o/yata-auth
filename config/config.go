package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	App      App            `yaml:"app"`
}

type PostgresConfig struct {
	Host     string `yaml:"PostgresqlHost" env:"POSTGRESQL_HOST"`
	Port     string `yaml:"PostgresqlPort" env:"POSTGRESQL_PORT"`
	User     string `yaml:"PostgresqlUser" env:"POSTGRESQL_USERNAME"`
	Password string `yaml:"PostgresqlPassword" env:"POSTGRESQL_PASSWORD"`
	Name     string `yaml:"PostgresqlDbname" env:"POSTGRESQL_NAME"`
}

type App struct {
	JWT  JWTConfig `yaml:"jwt"`
	Port string    `yaml:"port"`
}

type JWTConfig struct {
	Secret        string `yaml:"secret"`
	TokenTTLHours int    `yaml:"token_ttl_hours"`
	Salt          string `yaml:"salt"`
}

func LoadConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		log.Fatalf("error while reading config file: %s", err)
	}
	return &cfg

}
