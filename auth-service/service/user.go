package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
)

type ImageData struct {
	ContentType   string `json:"contentType"`
	BinaryContent []byte `json:"data"`
}

const passwordHashCost = 15

type UserInterface interface {
	GetUserById(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user models.User) (int, error)
	Update(user models.User) error
	Delete(id int) error
	ValidatePassword(hashedPassword string, password string) error
	HashPassword(password string) (string, error)
}

type UserService struct {
	C context.Context
}

func (s *UserService) GetUserById(ID string) (*models.User, error) {
	return models.FindUserG(s.C, ID)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return models.Users(qm.Where("email = $1", email)).OneG(s.C)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return models.Users(qm.Where("username = $1", username)).OneG(s.C)
}

func (s *UserService) CreateUser(dto *UserRegisterDTO) (*models.User, error) {
	hashedPassword, err := s.HashPassword(dto.Password)
	if err != nil {
		return nil, ErrHashPassword
	}
	user := &models.User{
		Email:          dto.Email,
		Username:       dto.Username,
		FullName:       dto.FullName,
		Phone:          dto.Phone,
		HashedPassword: hashedPassword,
	}
	err = user.InsertG(s.C, boil.Blacklist("id", "image", "is_oauth"))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(user *models.User) error {
	userInDB, _ := models.FindUserG(s.C, user.ID)
	if userInDB == nil {
		return ErrUserNotFound
	}
	_, err := user.UpdateG(s.C, boil.Infer())
	return err
}

func (s *UserService) Delete(ID string) error {
	user, _ := models.FindUserG(s.C, ID)
	if user != nil {
		_, err := user.DeleteG(s.C)
		return err
	}
	return nil
}

func (s *UserService) ValidatePassword(hashedPassword string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrWrongCredentials
		}
		return err
	}
	return nil
}

func (s *UserService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	return string(bytes), err
}

func (s *UserService) UpdateUserProfileImage(userID string, imageData *ImageData) error {
	user, err := models.FindUserG(s.C, userID)
	if err != nil {
		return err // User not found or other error
	}

	imageDataBytes, err := json.Marshal(imageData)
	if err != nil {
		return err
	}

	user.Image = null.BytesFrom(imageDataBytes)
	_, err = user.UpdateG(s.C, boil.Whitelist("image"))
	return err
}

func (s *UserService) GetUserImage(userID string) *ImageData {
	user, err := models.FindUserG(s.C, userID)
	if err != nil || user == nil {
		return nil
	}
	imageDataBytesPtr := user.Image.Ptr()
	if imageDataBytesPtr == nil {
		return nil
	}

	var imageData ImageData
	err = json.Unmarshal(*imageDataBytesPtr, &imageData)
	if err != nil {
		return nil
	}

	return &imageData
}
