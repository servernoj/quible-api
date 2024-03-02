package userService

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	userService := new(UserService)
	password := "password"
	hashedPassword, _ := userService.HashPassword(password)
	assert.Nil(t, userService.ValidatePassword(hashedPassword, password))
	assert.Error(t, userService.ValidatePassword(hashedPassword, "something-else"))
}
