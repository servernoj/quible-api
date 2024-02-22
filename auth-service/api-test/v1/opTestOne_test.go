package v1

import (
	_ "embed"

	"github.com/rs/zerolog/log"
)

//go:embed CSV/users.csv
var users_as_csv string

func (suite *TestCases) TestOne() {
	log.Info().Msg("TestOne")
	// t := suite.T()
	// assert := assert.New(t)
	// _, api := humatest.New(t)

	// store.InsertFromCSV(suite.T(), "users", users_as_csv)
	// response := api.Post("/api/v1/login", map[string]string{
	// 	"email":    "userA@gmail.com",
	// 	"password": "pass",
	// })
	// assert.EqualValues(response.Code, http.StatusUnauthorized)
	// assert.Contains(response.Body.String(), "4011001")
}
