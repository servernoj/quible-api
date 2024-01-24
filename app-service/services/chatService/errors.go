package chatService

import "errors"

var (
	ErrChatGroupNotFound    = errors.New("chat group not found")
	ErrChatGroupExists      = errors.New("chat group exists")
	ErrChannelExists        = errors.New("channel exists")
	ErrChannelNotFound      = errors.New("channel not found")
	ErrUserUndefined        = errors.New("user is not defined")
	ErrPrivateChatGroup     = errors.New("cannot join channel from the private chat group")
	ErrSelfOwnedChatGroup   = errors.New("cannot join channel from the self-owned chat group")
	ErrChannelAlreadyJoined = errors.New("channel already joined")
)
