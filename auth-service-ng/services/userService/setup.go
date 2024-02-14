package userService

import (
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func (s *UserService) ValidatePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *UserService) HashPassword(password string) (string, error) {
	passwordHashCost := 5
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	return string(bytes), err
}
