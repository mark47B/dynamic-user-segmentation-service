package settings

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Database struct {
	POSTGRES_DB_NAME  string `env:"POSTGRES_DB_NAME" env-default:"postgres"`
	POSTGRES_HOST     string `env:"POSTGRES_HOST" env-default:"localhost"`
	POSTGRES_PORT     string `env:"POSTGRES_PORT" env-default:"5432"`
	POSTGRES_USER     string `env:"POSTGRES_USER" env-default:"postgres"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD" env-required:"true"`
}

type Server struct {
	ENV            string        `env:"ENV" env-default:"local"`
	SERVER_ADDRESS string        `env:"SERVER_ADDRESS" env-default:"localhost:8080"`
	SERVER_TIMEOUT time.Duration `env:"SERVER_TIMEOUT" env-default:"4s"`
	IDLE_TIMEOUT   time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

type Settings struct {
	Server    Server
	Database  Database
	DebugMode bool
}

func MustLoad() *Settings {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Settings

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
