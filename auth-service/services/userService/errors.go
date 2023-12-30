package userService

import "errors"

var (
	ErrWrongCredentials = errors.New("invalid credentials")
	ErrUserNotFound     = errors.New("user not found")
	ErrHashPassword     = errors.New("unable to hash password")
)
