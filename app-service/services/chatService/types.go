package chatService

import (
	"github.com/quible-io/quible-api/lib/models"
)

type CreateChatGroupDTO struct {
	// used to form Ably resource identifier for chat group as `chat:<name>`,
	// it should be unique across all chat groups owned by the same user
	Name string `json:"name" binding:"alpha,required"`
	// human-readable "title" of the chat group, will be displayed in UI and used for searching
	Title string `json:"title" binding:"required"`
	// Optional summary, potentially lengthy text
	Summary *string `json:"summary" binding:"omitempty"`
	// Private chat groups require invitation from the owner.
	// Public chat group can be freely joined by using `/join` endpoint
	IsPrivate bool `json:"isPrivate" binding:"boolean,required"`
}

type CreateChannelDTO struct {
	// second part of the channel resource name to be concatenated with
	// the `chat group` name via `:`.
	Name string `json:"name" binding:"alpha,required"`
	// human-readable "title" of the channel
	Title string `json:"title" binding:"required"`
	// Optional summary, potentially lengthy text
	Summary *string `json:"summary" binding:"omitempty"`
}

type SearchResultItem struct {
	// parent public chat group for all listed channels
	Group *models.Chat `json:"chatGroup"`
	// all channels associated with the parent chat group
	Channels models.ChatSlice `json:"channels"`
}
