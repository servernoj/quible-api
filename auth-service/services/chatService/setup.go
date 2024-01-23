package chatService

import (
	"context"
	"fmt"

	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type ChatService struct {
	C context.Context
}

const (
	GROUP_PREFIX string = "chat:"
)

var (
	AccessReadOnly  = []string{"subscribe", "history"}
	AccessReadWrite = []string{"subscribe", "publish", "history"}
)

func getErrorWrapper(format string, args ...any) func(error) error {
	return func(err error) error {
		return fmt.Errorf(format, append(args, err))
	}
}

func (cs *ChatService) CreateChatGroup(user *models.User, name string, summary string, isPrivate bool) (*models.Chat, error) {
	if user == nil {
		return nil, fmt.Errorf("user is not defined")
	}
	userId := user.ID
	errorWrapper := getErrorWrapper("CreateChatGroup for user %s: %w", userId)
	resource := GROUP_PREFIX + name
	chatGroupFound, err := models.Chats(
		models.ChatWhere.OwnerID.EQ(null.StringFrom(userId)),
		models.ChatWhere.Resource.EQ(resource),
		models.ChatWhere.ParentID.IsNull(),
	).ExistsG(cs.C)
	if err != nil || chatGroupFound {
		return nil, errorWrapper(
			fmt.Errorf("unable to create chat group due to conflict: %w", err),
		)
	}
	chatGroup := models.Chat{
		Resource:  resource,
		ParentID:  null.StringFromPtr(nil),
		IsPrivate: null.BoolFrom(isPrivate),
		OwnerID:   null.StringFrom(userId),
		Summary:   summary,
	}
	if err := chatGroup.InsertG(cs.C, boil.Infer()); err != nil {
		return nil, errorWrapper(
			fmt.Errorf("unable to create chat group: %w", err),
		)
	}
	return &chatGroup, nil
}
func (cs *ChatService) CreateChannel(chatGroup *models.Chat, name string, summary string) (*models.Chat, error) {
	errorWrapper := getErrorWrapper("CreateChannel: %w")
	if chatGroup == nil {
		return nil, errorWrapper(
			fmt.Errorf("parent chat group must be defined, got `nil`"),
		)
	}
	channel := models.Chat{
		Resource:  name,
		Summary:   summary,
		ParentID:  null.StringFrom(chatGroup.ID),
		OwnerID:   null.StringFromPtr(nil),
		IsPrivate: null.BoolFromPtr(nil),
	}
	if err := channel.InsertG(cs.C, boil.Infer()); err != nil {
		return nil, errorWrapper(
			fmt.Errorf("unable to create channel: %w", err),
		)
	}
	return &channel, nil
}
func (cs *ChatService) GetChatGroups(user *models.User, extraMods ...qm.QueryMod) ([]*models.Chat, error) {
	mods := []qm.QueryMod{
		models.ChatWhere.ParentID.IsNull(),
	}
	if user != nil {
		mods = append(
			mods,
			models.ChatWhere.OwnerID.EQ(null.StringFrom(user.ID)),
		)
	}
	mods = append(mods, extraMods...)
	return models.Chats(mods...).AllG(cs.C)
}
func (cs *ChatService) GetPublicChatGroups(user *models.User) ([]*models.Chat, error) {
	return cs.GetChatGroups(
		user,
		models.ChatWhere.IsPrivate.EQ(null.BoolFrom(false)),
	)
}
func (cs *ChatService) GetPrivateChatGroups(user *models.User) ([]*models.Chat, error) {
	return cs.GetChatGroups(
		user,
		models.ChatWhere.IsPrivate.EQ(null.BoolFrom(true)),
	)
}
func (cs *ChatService) GetCapabilities(user *models.User) (map[string][]string, error) {
	userId := user.ID
	errorWrapper := getErrorWrapper("GetCapabilities for user %s: %w", userId)
	// -- initialize empty map
	capabilities := map[string][]string{}
	// -- process implied capabilities from self-owned chat groups
	chatGroups, err := cs.GetChatGroups(user)
	if err != nil {
		return nil, errorWrapper(
			fmt.Errorf("unable to retrieve chat groups owned by user: %w", err),
		)
	}
	for _, chatGroup := range chatGroups {
		resource := chatGroup.Resource + ":*"
		capabilities[resource] = AccessReadWrite
	}
	// -- process explicitly assigned channels
	chatUsers, err := user.ChatUsers().AllG(cs.C)
	if err != nil {
		return nil, errorWrapper(
			fmt.Errorf("unable to retrieve chat-user associations: %w", err),
		)
	}
	for _, item := range chatUsers {
		chatId := item.ChatID
		access := AccessReadWrite
		if item.IsRo {
			access = AccessReadOnly
		}
		chat, err := models.FindChatG(cs.C, chatId)
		if err != nil {
			return nil, errorWrapper(
				fmt.Errorf("unable to retrieve chat channel for chat_id %q: %w", chatId, err),
			)
		}
		if chat == nil {
			return nil, errorWrapper(
				fmt.Errorf("chat channel for chat_id not found %q", chatId),
			)
		}
		chatGroup, err := chat.Parent().OneG(cs.C)
		if err != nil {
			return nil, errorWrapper(
				fmt.Errorf("unable to retrieve chat group for chat_id %q: %w", chatId, err),
			)
		}
		if chatGroup == nil {
			return nil, errorWrapper(
				fmt.Errorf("chat group for chat_id not found %q", chatId),
			)
		}
		resource := chatGroup.Resource + ":" + chat.Resource
		if accessFound, ok := capabilities[resource]; ok && len(accessFound) > len(access) {
			access = accessFound
		}
		capabilities[resource] = access
	}
	return capabilities, nil
}
func (cs *ChatService) GetPublicChannelsByUser(user *models.User) ([]*models.Chat, error) {
	errorWrapper := getErrorWrapper("GetPublicChannels: %w")
	chatGroups, err := cs.GetPublicChatGroups(user)
	if err != nil {
		return nil, errorWrapper(
			fmt.Errorf("unable to retrieve public chat groups: %w", err),
		)
	}
	IDs := make([]string, len(chatGroups))
	for idx, chatGroup := range chatGroups {
		IDs[idx] = chatGroup.ID
	}
	return models.Chats(
		models.ChatWhere.ParentID.IN(IDs),
	).AllG(cs.C)
}
func (cs *ChatService) GetPublicChannelsByGroup(chatGroup *models.Chat) ([]*models.Chat, error) {
	if chatGroup == nil || chatGroup.IsPrivate.Bool {
		return []*models.Chat{}, nil
	}
	return models.Chats(
		models.ChatWhere.ParentID.EQ(null.StringFrom(chatGroup.ID)),
	).AllG(cs.C)
}
func (cs *ChatService) JoinPublicChannel(user *models.User, channel *models.Chat) error {
	if user == nil {
		return fmt.Errorf("user is not defined")
	}
	if channel == nil {
		return fmt.Errorf("channel is not defined")
	}
	errorWrapper := getErrorWrapper("JoinPublicChannel for user %s: %w", user.ID)
	chatGroup, err := channel.Parent().OneG(cs.C)
	if err != nil || chatGroup == nil {
		return errorWrapper(
			fmt.Errorf("unable to locate parent chat group"),
		)
	}
	if chatGroup.IsPrivate.Bool {
		return errorWrapper(
			fmt.Errorf("channel is private"),
		)
	}
	chatUser := models.ChatUser{
		ChatID: channel.ID,
		UserID: user.ID,
	}
	return chatUser.InsertG(cs.C, boil.Infer())
}
