package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DbUrl     string `env:"DB_URL"`
	Port      string `env:"PORT" envDefault:"8080"`
	JWTSecret string `env:"JWT_SECRET"`
}

var Config = &EnvConfig{}

func LoadEnvConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	Config = &EnvConfig{
		DbUrl:     os.Getenv("DB_URL"),
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	log.Printf("Environment configuration loaded: %+v\n", Config)
	return nil
}
