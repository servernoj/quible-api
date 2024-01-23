package chatService

import "errors"

var (
	ErrUnknownError = errors.New("unknown error")
	ErrUserNotFound = errors.New("user not found")
)
