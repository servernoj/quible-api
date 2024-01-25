package chatService

import "errors"

var (
	ErrChatGroupNotFound               = errors.New("chat group not found")
	ErrChatGroupExists                 = errors.New("chat group exists")
	ErrChannelExists                   = errors.New("channel exists")
	ErrChannelNotFound                 = errors.New("channel not found")
	ErrUserUndefined                   = errors.New("user is not defined")
	ErrPrivateChatGroup                = errors.New("cannot join channel from the private chat group")
	ErrSelfOwnedChatGroup              = errors.New("cannot join channel from the self-owned chat group")
	ErrChannelAlreadyJoined            = errors.New("channel already joined")
	ErrUnableToRetrtieveChatGroups     = errors.New("unable to retrieve chat groups owned by user")
	ErrUnableToRetrtieveChatGroup      = errors.New("unable to retrieve chat group")
	ErrUnableToRetrtieveChannel        = errors.New("unable to retrieve channel")
	ErrUnableToRetrieveChatUserRecords = errors.New("unable to retrieve chat-user associations")
	ErrUnableToCreateChatGroup         = errors.New("unable to create chat group")
	ErrUnableToCreateChannel           = errors.New("unable to create channel")
)
