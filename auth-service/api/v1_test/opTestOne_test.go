package v1_test

import (
	_ "embed"

	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/stretchr/testify/assert"
)

//go:embed CSV/users.csv
var users_as_csv string

func (suite *TestSuite) TestOne() {
	assert := assert.New(suite.T())
	store.InsertFromCSV(suite.T(), "users", users_as_csv)
	numUsers, err := models.Users().CountG(suite.ctx)
	if err != nil {
		suite.T().Fatalf("Couldn't retrieve the number of users")
	}
	assert.EqualValues(3, numUsers)
}
