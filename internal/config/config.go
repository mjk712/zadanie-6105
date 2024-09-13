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

/*
export SERVER_ADDRESS=0.0.0.0:8080
export POSTGRES_CONN=postgres://postgres:1234@localhost:5432/tender?sslmode=disable
export POSTGRES_JDBC_URL=jdbc:postgresql://localhost:5432/tender
export POSTGRES_USERNAME=postgres
export POSTGRES_PASSWORD=1234
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_DATABASE=PostgreSQL
export ENV=prod
*/
