package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/quible-io/quible-api/lib/migrations"
)

func main() {

	envFile := "../.env"
	DSN := os.Getenv("ENV_DSN")
	command, args := os.Args[1], os.Args[2:]

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
	goose.SetBaseFS(migrations.FS)
	if err := goose.Run(command, db, ".", args...); err != nil {
		log.Fatalf("migrate %v: %v", command, err)
	}
}
