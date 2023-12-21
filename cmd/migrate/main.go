package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/quible-io/quible-api/cmd/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {

	envFile := "../.env"
	DSN := os.Getenv("ENV_DSN")
	dir, command, args := "./migrations", os.Args[1], os.Args[2:]

	if DSN == "" && godotenv.Load(envFile) == nil {
		DSN = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			"localhost",
			5432,
			os.Getenv("POSTGRES_DB"),
		)
	}

	db, err := goose.OpenDBWithDriver("pgx", DSN)
	if err != nil {
		log.Fatalf("migrate: failed to open DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("migrate: failed to connect to DB server: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("migrate: failed to close DB: %v\n", err)
		}
	}()

	if err := goose.Run(command, db, dir, args...); err != nil {
		log.Fatalf("migrate %v: %v", command, err)
	}
}
