package service

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.com/quible-backend/lib/models"
)

type UserInterface interface {
	GetUserById(id int64) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user models.User) (int64, error)
	Update(id int64, user models.User) (int64, error)
	Delete(id int64) error
	ValidatePassword(hashedPassword string, password string) error
	HashPassword(password string) (string, error)
}

type UserService struct {
	C context.Context
}

func (s *UserService) GetUserById(id int64) (*models.User, error) {
	return models.Users(qm.Where("id = $1", id)).OneG(s.C)
}
