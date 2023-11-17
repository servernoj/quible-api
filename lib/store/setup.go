package store

import (
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Setup(dsn string) error {
	if len(dsn) == 0 {
		return errors.New("undefined DSN")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return errors.New("DB server cannot be reached")
	}
	boil.SetDB(db)

	return nil
}

func Close() {
	if db, ok := boil.GetDB().(*sql.DB); ok {
		db.Close()
	}
}
