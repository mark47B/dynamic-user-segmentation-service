package settings

import (
	"os"
)

type Database struct {
	POSTGRES_DB_NAME  string `env:"POSTGRES_DB_NAME" env-default:"postgres"`
	POSTGRES_HOST     string `env:"POSTGRES_HOST" env-default:"localhost"`
	POSTGRES_PORT     string `env:"POSTGRES_PORT" env-default:"5432"`
	POSTGRES_USER     string `env:"POSTGRES_USER" env-default:"postgres"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD" env-required:"true"`
}

type Server struct {
	ENV            string `env:"ENV" env-default:"local"`
	SERVER_ADDRESS string `env:"SERVER_ADDRESS" env-default:"localhost:8080"`
}

type Settings struct {
	Server    Server
	Database  Database
	DebugMode bool
}

func MustLoad() *Settings {
	var cfg Settings
	cfg.Server.ENV = os.Getenv("ENV")
	cfg.Server.SERVER_ADDRESS = os.Getenv("SERVER_ADDRESS")

	cfg.Database.POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	cfg.Database.POSTGRES_USER = os.Getenv("POSTGRES_USER")
	cfg.Database.POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	cfg.Database.POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
	cfg.Database.POSTGRES_DB_NAME = os.Getenv("POSTGRES_DB_NAME")

	return &cfg
}
