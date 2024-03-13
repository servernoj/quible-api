package suite

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"strings"
	"testing"

	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/strmangle"
)

func InsertFromCSV(db *sql.DB, tableName string, csv_as_string string) error {
	isNull := func(s string) bool {
		return len(s) == 0
	}
	isBase64 := func(s string, col string) bool {
		return s[len(s)-1] == '=' && col == "image"
	}
	reader := csv.NewReader(strings.NewReader(csv_as_string))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("unable to read CSV records: %w", err)
	}

	headers := records[0]
	stmt, err := db.Prepare(
		fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES(%s)",
			tableName,
			strings.Join(headers, ","),
			strmangle.Placeholders(true, len(headers), 1, 1),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to prepare SQL insert statement: %w", err)
	}
	for _, data := range records[1:] {
		args := make([]any, len(data))
		for idx := range data {
			switch {
			case isNull(data[idx]):
				args[idx] = nil
			case isBase64(data[idx], headers[idx]):
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
			return fmt.Errorf("unable to execute SQL insert statement: %w", err)
		}
	}
	return nil
}

func GetToken(t *testing.T, db *sql.DB, userId string, action jwt.TokenAction) string {
	user, err := models.FindUser(context.Background(), db, userId)
	if err != nil {
		t.Fatalf("unable to retrieve user record from DB: %q", err)
	}
	token, err := jwt.GenerateToken(user, action, nil)
	if err != nil {
		t.Fatal("unable to generate token")
	}
	return token.String()
}
