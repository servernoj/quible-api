package userService

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

type UserService struct {
	C context.Context
}

func (s *UserService) GetUserById(ID string, cols ...string) (*models.User, error) {
	return models.FindUserG(s.C, ID, cols...)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return models.Users(qm.Where("email = $1", email)).OneG(s.C)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return models.Users(qm.Where("username = $1", username)).OneG(s.C)
}

func (s *UserService) GetUserByUsernameOrEmail(dto *UserRegisterDTO) (*models.User, error) {
	return models.Users(
		qm.Or2(models.UserWhere.Email.EQ(dto.Email)),
		qm.Or2(models.UserWhere.Username.EQ(dto.Username)),
	).OneG(s.C)
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

	if err != nil {
		return nil, err
	}

	err = user.InsertG(s.C, boil.Blacklist("id", "image"))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(user *models.User) error {
	if userExists, err := models.UserExistsG(s.C, user.ID); err != nil || !userExists {
		return ErrUserNotFound
	}
	_, err := user.UpdateG(s.C, boil.Infer())
	return err
}

func (s *UserService) UpdateWith(user *models.User, dto *UserRegisterDTO) error {
	if userExists, err := models.UserExistsG(s.C, user.ID); err != nil || !userExists {
		return ErrUserNotFound
	}
	hashedPassword, err := s.HashPassword(dto.Password)
	if err != nil {
		return ErrHashPassword
	}
	user.Email = dto.Email
	user.Username = dto.Username
	user.FullName = dto.FullName
	user.Phone = dto.Phone
	user.HashedPassword = hashedPassword

	_, err = user.UpdateG(s.C, boil.Infer())
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

func (s *UserService) GetUserImage(user *models.User) *ImageData {
	imageDataBytesPtr := user.Image.Ptr()
	var imageData ImageData
	if imageDataBytesPtr == nil {
		return nil
	}
	if json.Unmarshal(*imageDataBytesPtr, &imageData) != nil {
		return nil
	}

	return &imageData
}
