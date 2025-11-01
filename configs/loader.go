package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func MustLoad() Config {
	cnf := Config{}

	env := os.Getenv("APP_ENV")
	if env == "" || env == "local" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(".env file not found, using environment variables")
		}
	}

	if err := envconfig.Process("", &cnf); err != nil {
		log.Fatal(err)
	}

	return cnf
}
