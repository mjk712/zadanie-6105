package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	Env string `envconfig:"ENV" env-default:"local"`
	HTTPServer
	DBConfig
}

type HTTPServer struct {
	Address string `envconfig:"SERVER_ADDRESS" `
}

type DBConfig struct {
	ConnectionString string `envconfig:"POSTGRES_CONN"`
	JDBCUrl          string `envconfig:"POSTGRES_JDBC_URL"`
	Username         string `envconfig:"POSTGRES_USERNAME"`
	Password         string `envconfig:"POSTGRES_PASSWORD"`
	Host             string `envconfig:"POSTGRES_HOST"`
	Port             int    `envconfig:"POSTGRES_PORT"`
	Database         string `envconfig:"POSTGRES_DATABASE"`
}

func New() *Config {
	const op = "Config Load"
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(op, err)
	}
	return &cfg
}
