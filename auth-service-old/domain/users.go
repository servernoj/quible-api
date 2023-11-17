package domain

import "time"

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	FullName       string    `json:"full_name"`
	Phone          string    `json:"phone"`
	Image          string    `json:"image"`
	IsOauth        bool      `json:"is_oauth"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserRegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Image    string `json:"image"`
}

type UserEmailCheckRequest struct {
	Email string `json:"email" validate:"required"`
}

type UserVerifyCodeRequest struct {
	Email string `json:"email" validate:"required"`
	Code  int    `json:"code" validate:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Image    string `json:"image"`
}

type UserChangePasswordRequest struct {
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	IsOauth   bool      `json:"is_oauth"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLoginResponse struct {
	ID             int64  `json:"id" validate:"required"`
	Username       string `json:"username"`
	Email          string `json:"email" validate:"required"`
	HashedPassword string `json:"hashed_password" validate:"required"`
}
