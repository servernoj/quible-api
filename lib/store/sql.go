package store

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Init(dsn string) error {

	// Open handle to database like normal
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	boil.SetDB(db)

	return nil
}
