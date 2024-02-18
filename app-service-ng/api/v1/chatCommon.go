package v1

import (
	"context"
	"fmt"

	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	GROUP_PREFIX string = "chat:"
)

var (
	AccessReadOnly  = []string{"subscribe", "history"}
	AccessReadWrite = []string{"subscribe", "publish", "history"}
)

type ChatChannel struct {
	ID       string       `json:"id"`
	Title    string       `json:"title"`
	Resource string       `json:"resource"`
	ReadOnly bool         `json:"readOnly"`
	Parent   *models.Chat `json:"-"`
}

func chatChannelsForUser(ctx context.Context, userId string) ([]ChatChannel, error) {
	// 0. Initialize storage
	chatChannelByChatId := map[string]ChatChannel{}
	chatChannels := []ChatChannel{}
	// 1. Identify chat channels associated with groups, owned by user
	ownedChatGroups, err := models.Chats(
		models.ChatWhere.OwnerID.EQ(null.StringFrom(userId)),
		qm.Load(
			qm.Rels(
				models.ChatRels.ParentChats,
				models.ChatRels.Parent,
			),
		),
	).AllG(ctx)
	if err != nil {
		return nil, ErrorMap.GetErrorResponse(
			Err500_UnknownError,
			fmt.Errorf("unable to retrieve chat channels owned by user %q", userId),
			err,
		)
	}
	for _, chatGroup := range ownedChatGroups {
		for _, chat := range chatGroup.R.ParentChats {
			chatId := chat.ID
			if _, ok := chatChannelByChatId[chatId]; !ok {
				chatChannelByChatId[chatId] = ChatChannel{
					ID:       chatId,
					Title:    chat.Title,
					Resource: chat.R.Parent.Resource + ":" + chat.Resource,
					ReadOnly: false,
					Parent:   chat.R.Parent,
				}
			}
		}
	}
	// 2. Identify chat channels that user has joined (been invited)
	chatUsersLoaded, err := models.ChatUsers(
		models.ChatUserWhere.UserID.EQ(userId),
		models.ChatUserWhere.Disabled.EQ(false),
		qm.Load(
			qm.Rels(
				models.ChatUserRels.Chat,
				models.ChatRels.Parent,
			),
		),
	).AllG(ctx)
	if err != nil {
		return nil, ErrorMap.GetErrorResponse(
			Err500_UnknownError,
			fmt.Errorf("unable to chat channels joined by user %q", userId),
			err,
		)
	}
	for _, chatUser := range chatUsersLoaded {
		chat := chatUser.R.Chat
		chatId := chat.ID
		if _, ok := chatChannelByChatId[chatId]; !ok {
			chatChannelByChatId[chatId] = ChatChannel{
				ID:       chatId,
				Title:    chat.Title,
				Resource: chat.R.Parent.Resource + ":" + chat.Resource,
				ReadOnly: chatUser.IsRo,
				Parent:   chat.R.Parent,
			}
		}
	}
	// 3. Compile final result
	for _, chatChannel := range chatChannelByChatId {
		chatChannels = append(chatChannels, chatChannel)
	}
	return chatChannels, err
}
