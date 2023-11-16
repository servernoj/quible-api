package store

import (
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Init(dsn string) error {
	if len(dsn) == 0 {
		return errors.New("undefined DSN")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	boil.SetDB(db)

	return nil
}
