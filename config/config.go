package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
	App      App            `yaml:"app"`
}

type PostgresConfig struct {
	Host     string `yaml:"PostgresqlHost" env:"POSTGRESQL_HOST"`
	Port     string `yaml:"PostgresqlPort" env:"POSTGRESQL_PORT"`
	User     string `yaml:"PostgresqlUser" env:"POSTGRESQL_USERNAME"`
	Password string `yaml:"PostgresqlPassword" env:"POSTGRESQL_PASSWORD"`
	Name     string `yaml:"PostgresqlDbname" env:"POSTGRESQL_NAME"`
}

type RedisConfig struct {
	Host     string `yaml:"RedisHost" env:"REDISHOST"`
	Port     string `yaml:"RedisPort" env:"REDISPORT"`
	User     string `yaml:"RedisUser" env:"REDISUSER"`
	Password string `yaml:"RedisPassword" env:"REDISPASSWORD"`
	DB       int    `yaml:"RedisDB" env:"REDISDB"`
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
