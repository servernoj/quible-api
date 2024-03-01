package store

import (
	"database/sql"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"strings"
	"testing"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/strmangle"
)

type InsertOptionsFunc func(*InsertOptions)

type InsertOptions struct {
	isNull   func(string) bool
	isBase64 func(string, string) bool
	db       *sql.DB
}

func InsertWithIsNull(isNull func(string) bool) InsertOptionsFunc {
	return func(options *InsertOptions) {
		options.isNull = isNull
	}
}
func InsertWithDB(db *sql.DB) InsertOptionsFunc {
	return func(options *InsertOptions) {
		options.db = db
	}
}

func InsertFromCSV(t *testing.T, tableName string, csv_as_string string, opts ...InsertOptionsFunc) {
	options := InsertOptions{
		isNull: func(s string) bool {
			return len(s) == 0
		},
		isBase64: func(s string, col string) bool {
			return s[len(s)-1] == '=' && col == "image"
		},
		db: boil.GetDB().(*sql.DB),
	}
	for _, opt := range opts {
		opt(&options)
	}
	reader := csv.NewReader(strings.NewReader(csv_as_string))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatal("unable to read CSV records", err)
	}

	headers := records[0]
	stmt, err := options.db.Prepare(
		fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES(%s)",
			tableName,
			strings.Join(headers, ","),
			strmangle.Placeholders(true, len(headers), 1, 1),
		),
	)
	if err != nil {
		t.Fatal("unable to prepare SQL insert statement", err)
	}
	for _, data := range records[1:] {
		args := make([]any, len(data))
		for idx := range data {
			switch {
			case options.isNull(data[idx]):
				args[idx] = nil
			case options.isBase64(data[idx], headers[idx]):
				{
					decoded, err := base64.StdEncoding.DecodeString(data[idx])
					if err != nil {
						args[idx] = data[idx]
					}
					args[idx] = decoded
				}
			default:
				{
					args[idx] = data[idx]
				}
			}
		}
		if _, err := stmt.Exec(args...); err != nil {
			t.Fatal("unable to execute SQL insert statement", err)
		}
	}
}
