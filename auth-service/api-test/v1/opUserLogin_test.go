package v1

import (
	"context"
	_ "embed"
	"net/http"

	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/stretchr/testify/assert"
)

//go:embed TestData/users.csv
var users_as_csv string

func (suite *TestCases) TestUserLogin() {
	t := suite.T()
	assert := assert.New(t)
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", users_as_csv)
	// 2. Assert that correct number of users were imported
	count, err := models.Users().CountG(context.Background())
	assert.Nil(err, "users counting should not return error")
	assert.EqualValues(count, 3, "all users imported from CSV should be discoverable")
	// 3. Try to login with correct credentials and expect success
	response := suite.TestAPI.Post("/api/v1/login", map[string]string{
		"email":    "userA@gmail.com",
		"password": "password",
	})
	assert.EqualValues(response.Code, http.StatusOK)
}
