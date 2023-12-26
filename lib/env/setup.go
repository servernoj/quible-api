package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const envFile = "../.env"

func Setup() {
	if os.Getenv("IS_DOCKER") != "1" {
		if err := godotenv.Load(envFile); err == nil {
			if os.Getenv("ENV_DSN") == "" {
				os.Setenv(
					"ENV_DSN",
					fmt.Sprintf(
						"postgres://%s:%s@%s:%d/%s",
						os.Getenv("POSTGRES_USER"),
						os.Getenv("POSTGRES_PASSWORD"),
						"localhost",
						5432,
						os.Getenv("POSTGRES_DB"),
					),
				)
			}
			if os.Getenv("ENV_URL_AUTH_SERVICE") == "" {
				os.Setenv(
					"ENV_URL_AUTH_SERVICE",
					fmt.Sprintf(
						"http://localhost:%s",
						os.Getenv("AUTH_PORT"),
					),
				)
			}
			if os.Getenv("ENV_URL_APP_SERVICE") == "" {
				os.Setenv(
					"ENV_URL_APP_SERVICE",
					fmt.Sprintf(
						"http://localhost:%s",
						os.Getenv("APP_PORT"),
					),
				)
			}
		}
	} else {
		log.Println("running in docker...")
	}
}
