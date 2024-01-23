package chatService

import "errors"

var (
	ErrChatGroupNotFound = errors.New("chat group not found")
	ErrChatGroupExists   = errors.New("chat group exists")
	ErrChannelExists     = errors.New("channel exists")
)
