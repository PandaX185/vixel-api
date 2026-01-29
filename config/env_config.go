package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DbUrl           string `env:"DB_URL"`
	Port            string `env:"PORT" envDefault:"8080"`
	JWTSecret       string `env:"JWT_SECRET"`
	MINIOEndpoint   string `env:"MINIO_ENDPOINT"`
	MINIOAccessKey  string `env:"MINIO_ACCESS_KEY"`
	MINIOSecretKey  string `env:"MINIO_SECRET_KEY"`
	MINIOUseSSL     bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	MINIOBucketName string `env:"MINIO_BUCKET_NAME"`
	MINIORegion     string `env:"MINIO_REGION"`
}

var Config = &EnvConfig{}

func LoadEnvConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	Config = &EnvConfig{
		DbUrl:           os.Getenv("DB_URL"),
		Port:            os.Getenv("PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		MINIOEndpoint:   os.Getenv("MINIO_ENDPOINT"),
		MINIOAccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
		MINIOSecretKey:  os.Getenv("MINIO_SECRET_KEY"),
		MINIOUseSSL:     os.Getenv("MINIO_USE_SSL") == "true",
		MINIOBucketName: os.Getenv("MINIO_BUCKET_NAME"),
		MINIORegion:     os.Getenv("MINIO_REGION"),
	}

	log.Printf("Environment configuration loaded: %+v\n", Config)
	return nil
}
