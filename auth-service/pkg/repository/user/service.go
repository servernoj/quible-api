package user

import (
	"errors"
	"strings"
	"time"

	"gitlab.com/quible-backend/auth-service/domain"
	"golang.org/x/crypto/bcrypt"
)

type Impl interface {
	Gets(id int64) (*domain.UserResponse, error)
	GetUserByEmail(email string) (*domain.UserResponse, error)
	GetByUsername(username string) (*domain.UserResponse, error)
	Create(user domain.UserRegisterRequest) (int64, error)
	Update(id int64, user domain.UserUpdateRequest) (int64, error)
	Delete(id int64) (int64, error)
	GetLoginCredential(email string) (*domain.UserLoginResponse, error)
	ValidatePassword(hashedPassword string, password string) error
	HashPassword(password string) (string, error)
}

type Service struct {
	db Database
}

func NewService(db Database) Impl {
	return &Service{db: db}
}

func (s *Service) Create(user domain.UserRegisterRequest) (int64, error) {
	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	new := domain.User{
		Username:       strings.ToLower(user.Username),
		Email:          strings.ToLower(user.Email),
		HashedPassword: string(hashedPassword),
		FullName:       user.FullName,
		Phone:          user.Phone,
		Image:          strings.ToLower(user.Image),
		IsOauth:        false,
	}

	id, err := s.db.Create(new)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) Gets(id int64) (*domain.UserResponse, error) {
	user, err := s.db.Gets(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserByEmail(email string) (*domain.UserResponse, error) {
	user, err := s.db.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetByUsername(username string) (*domain.UserResponse, error) {
	return s.db.GetByUsername(username)
}

func (s *Service) GetLoginCredential(email string) (*domain.UserLoginResponse, error) {
	user, err := s.db.GetLoginCredential(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Update(id int64, user domain.UserUpdateRequest) (int64, error) {
	new := domain.User{
		ID:        id,
		Username:  strings.ToLower(user.Username),
		Email:     strings.ToLower(user.Email),
		FullName:  user.FullName,
		Phone:     user.Phone,
		Image:     user.Image,
		IsOauth:   false,
		UpdatedAt: time.Now(),
	}

	id, err := s.db.Update(new)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) Delete(id int64) (int64, error) {
	id, err := s.db.Delete(id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) ValidatePassword(hashedPassword string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors.New("password not match")
		}
		return err
	}
	return nil
}

func (s *Service) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	return string(bytes), err
}
