package chatService

import (
	"context"
	"fmt"
	"log"

	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	GROUP_PREFIX string = "chat:"
)

var (
	AccessReadOnly  = []string{"subscribe", "history"}
	AccessReadWrite = []string{"subscribe", "publish", "history"}
)

func getErrorWrapper(format string, args ...any) func(...error) error {
	return func(errors ...error) error {
		return fmt.Errorf(
			"%s: %w",
			fmt.Sprintf(format, args...),
			errors,
		)
	}
}

type ChatService struct {
	C context.Context
}

func (cs *ChatService) CreateChatGroup(
	user *models.User,
	name string,
	title string,
	summary *string,
	isPrivate bool,
) (*models.Chat, error) {
	if user == nil {
		return nil, ErrUserUndefined
	}
	userId := user.ID
	errorWrapper := getErrorWrapper("CreateChatGroup for user %q", userId)
	resource := GROUP_PREFIX + name
	chatGroupFound, err := models.Chats(
		models.ChatWhere.OwnerID.EQ(null.StringFrom(userId)),
		models.ChatWhere.ParentID.IsNull(),
		qm.Expr(
			qm.Or2(models.ChatWhere.Resource.EQ(resource)),
			qm.Or2(models.ChatWhere.Title.EQ(title)),
		),
	).ExistsG(cs.C)
	if err != nil {
		return nil, errorWrapper(ErrUnableToRetrtieveChatGroup, err)
	}
	if chatGroupFound {
		return nil, errorWrapper(ErrChatGroupExists)
	}
	chatGroup := models.Chat{
		Resource:  resource,
		ParentID:  null.StringFromPtr(nil),
		IsPrivate: null.BoolFrom(isPrivate),
		OwnerID:   null.StringFrom(userId),
		Summary:   null.StringFromPtr(summary),
		Title:     title,
	}
	if err := chatGroup.InsertG(cs.C, boil.Infer()); err != nil {
		return nil, errorWrapper(ErrUnableToCreateChatGroup, err)
	}
	return &chatGroup, nil
}
func (cs *ChatService) CreateChannel(
	chatGroupId string,
	name string,
	title string,
	summary *string,
) (*models.Chat, error) {
	errorWrapper := getErrorWrapper("CreateChannel")
	chatGroupFound, err := models.Chats(
		models.ChatWhere.ID.EQ(chatGroupId),
		models.ChatWhere.ParentID.IsNull(),
	).ExistsG(cs.C)
	if err != nil || !chatGroupFound {
		return nil, errorWrapper(ErrChatGroupNotFound)
	}
	channelFound, err := models.Chats(
		models.ChatWhere.Resource.EQ(name),
		models.ChatWhere.ParentID.EQ(null.StringFrom(chatGroupId)),
	).ExistsG(cs.C)
	if err != nil || channelFound {
		return nil, errorWrapper(ErrChannelExists)
	}
	channel := models.Chat{
		Resource:  name,
		Title:     title,
		Summary:   null.StringFromPtr(summary),
		ParentID:  null.StringFrom(chatGroupId),
		OwnerID:   null.StringFromPtr(nil),
		IsPrivate: null.BoolFromPtr(nil),
	}
	if err := channel.InsertG(cs.C, boil.Infer()); err != nil {
		return nil, errorWrapper(ErrUnableToCreateChannel, err)
	}
	return &channel, nil
}
func (cs *ChatService) GetChatGroups(user *models.User, extraMods ...qm.QueryMod) (models.ChatSlice, error) {
	errorWrapper := getErrorWrapper("GetChatGroups")
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
	chatGroups, err := models.Chats(mods...).AllG(cs.C)
	if err != nil {
		return nil, errorWrapper(err)
	}
	if len(chatGroups) == 0 {
		return models.ChatSlice{}, nil
	}
	return chatGroups, nil
}
func (cs *ChatService) GetPublicChatGroups(user *models.User) (models.ChatSlice, error) {
	return cs.GetChatGroups(
		user,
		models.ChatWhere.IsPrivate.EQ(null.BoolFrom(false)),
	)
}
func (cs *ChatService) GetPrivateChatGroups(user *models.User) (models.ChatSlice, error) {
	return cs.GetChatGroups(
		user,
		models.ChatWhere.IsPrivate.EQ(null.BoolFrom(true)),
	)
}
func (cs *ChatService) GetCapabilities(user *models.User) (map[string][]string, error) {
	if user == nil {
		return nil, ErrUserUndefined
	}
	userId := user.ID
	errorWrapper := getErrorWrapper("GetCapabilities for user %q", userId)
	// -- initialize empty map
	capabilities := map[string][]string{}
	// -- process implied capabilities from self-owned chat groups
	chatGroups, err := cs.GetChatGroups(user)
	if err != nil {
		return nil, errorWrapper(ErrUnableToRetrtieveChatGroups, err)
	}
	for _, chatGroup := range chatGroups {
		resource := chatGroup.Resource + ":*"
		capabilities[resource] = AccessReadWrite
	}
	// -- process joined channels
	chatUsers, err := user.ChatUsers().AllG(cs.C)
	if err != nil {
		return nil, errorWrapper(ErrUnableToRetrieveChatUserRecords, err)
	}
	for _, item := range chatUsers {
		chatId := item.ChatID
		access := AccessReadWrite
		if item.IsRo {
			access = AccessReadOnly
		}
		chat, err := models.FindChatG(cs.C, chatId)
		if err != nil || chat == nil {
			return nil, errorWrapper(ErrUnableToRetrtieveChannel, err)
		}
		parentChatGroup, err := chat.Parent().OneG(cs.C)
		if err != nil || parentChatGroup == nil {
			return nil, errorWrapper(ErrUnableToRetrtieveChatGroup, err)
		}
		resource := parentChatGroup.Resource + ":" + chat.Resource
		if accessFound, ok := capabilities[resource]; ok && len(accessFound) > len(access) {
			access = accessFound
		}
		capabilities[resource] = access
	}
	return capabilities, nil
}
func (cs *ChatService) GetMyGroupedChannels(user *models.User) ([]GroupedChannels, error) {
	if user == nil {
		return nil, ErrUserUndefined
	}
	userId := user.ID
	errorWrapper := getErrorWrapper("GetMyChannels for user %q", userId)
	// -- initialize empty list
	result := []GroupedChannels{}
	// -- process implied channels from self-owned chat groups
	myChatGroups, err := cs.GetChatGroups(user)
	if err != nil {
		return nil, errorWrapper(ErrUnableToRetrtieveChatGroups)
	}
	for _, chatGroup := range myChatGroups {
		result = append(
			result,
			GroupedChannels{
				ID:      chatGroup.ID,
				Title:   chatGroup.Title,
				Summary: chatGroup.Summary.Ptr(),
				Channels: []Channel{
					{
						ID:       chatGroup.ID,
						Resource: chatGroup.Resource + ":*",
						Title:    chatGroup.Title,
						ReadOnly: false,
					},
				},
			},
		)
	}
	// -- process joined channels
	chatUsers, err := user.ChatUsers().AllG(cs.C)
	if err != nil {
		return nil, errorWrapper(ErrUnableToRetrieveChatUserRecords, err)
	}
	groupedChannelsMap := make(map[string]GroupedChannels)
	for _, chatUser := range chatUsers {
		chat, err := models.FindChatG(cs.C, chatUser.ChatID)
		if err != nil || chat == nil {
			return nil, errorWrapper(ErrUnableToRetrtieveChannel, err)
		}
		parentChatGroup, err := chat.Parent().OneG(cs.C)
		if err != nil || parentChatGroup == nil {
			return nil, errorWrapper(ErrUnableToRetrtieveChatGroup, err)
		}
		if _, ok := groupedChannelsMap[parentChatGroup.ID]; !ok {
			groupedChannelsMap[parentChatGroup.ID] = GroupedChannels{
				ID:       parentChatGroup.ID,
				Title:    parentChatGroup.Title,
				Summary:  parentChatGroup.Summary.Ptr(),
				Channels: []Channel{},
			}
		}
		groupedChannels := groupedChannelsMap[parentChatGroup.ID]
		groupedChannels.Channels = append(groupedChannels.Channels, Channel{
			ID:       chat.ID,
			Resource: parentChatGroup.Resource + ":" + chat.Resource,
			Title:    chat.Title,
			ReadOnly: chatUser.IsRo,
		})
	}
	for _, groupedChannels := range groupedChannelsMap {
		result = append(result, groupedChannels)
	}
	return result, nil
}
func (cs *ChatService) GetMyChannels(user *models.User) ([]Channel, error) {
	if user == nil {
		return nil, ErrUserUndefined
	}
	userId := user.ID
	errorWrapper := getErrorWrapper("GetMyResources for user %q", userId)
	groupedChannels, err := cs.GetMyGroupedChannels(user)
	if err != nil {
		return nil, errorWrapper(err)
	}
	result := []Channel{}
	for idx := range groupedChannels {
		result = append(result, groupedChannels[idx].Channels...)
	}
	return result, nil
}
func (cs *ChatService) GetPublicChannelsByUser(user *models.User) (models.ChatSlice, error) {
	errorWrapper := getErrorWrapper("GetPublicChannels")
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
func (cs *ChatService) JoinPublicChannel(user *models.User, channelId string) error {
	if user == nil {
		return ErrUserUndefined
	}
	errorWrapper := getErrorWrapper("JoinPublicChannel for user %q", user.ID)
	channel, err := models.FindChatG(cs.C, channelId)
	if err != nil || channel == nil {
		return errorWrapper(ErrChannelNotFound)
	}
	chatGroup, err := channel.Parent().OneG(cs.C)
	if err != nil || chatGroup == nil {
		return errorWrapper(ErrChatGroupNotFound)
	}
	if chatGroup.IsPrivate.Bool {
		return errorWrapper(ErrPrivateChatGroup)
	}
	if chatGroup.OwnerID.String == user.ID {
		return errorWrapper(ErrSelfOwnedChatGroup)
	}
	chatUserFound, err := models.ChatUserExistsG(cs.C, channel.ID, user.ID)
	if err != nil || chatUserFound {
		return errorWrapper(ErrChannelAlreadyJoined)
	}

	chatUser := models.ChatUser{
		ChatID: channel.ID,
		UserID: user.ID,
	}
	return chatUser.InsertG(cs.C, boil.Infer())
}
func (cs *ChatService) LeaveChannel(user *models.User, channelId string) error {
	if user == nil {
		return ErrUserUndefined
	}
	errorWrapper := getErrorWrapper("LeaveChannel for user %q", user.ID)
	chatUser, err := models.ChatUsers(
		models.ChatUserWhere.ChatID.EQ(channelId),
		models.ChatUserWhere.UserID.EQ(user.ID),
	).OneG(cs.C)
	if err != nil || chatUser == nil {
		return errorWrapper(ErrChannelNotFound)
	}
	if _, err := chatUser.DeleteG(cs.C); err != nil {
		return errorWrapper(err)
	}
	return nil
}
func (cs *ChatService) SearchPublicChannelsByChatGroupTitle(chatGroupTitle string) []SearchResultItem {
	errorWrapper := getErrorWrapper("SearchPublicChannelsByChatGroupTitle with query %q", chatGroupTitle)
	chatGroups, err := models.Chats(
		models.ChatWhere.Title.ILIKE("%"+chatGroupTitle+"%"),
		models.ChatWhere.ParentID.IsNull(),
		models.ChatWhere.IsPrivate.EQ(null.BoolFrom(false)),
	).AllG(cs.C)
	if err != nil {
		log.Println(
			errorWrapper(
				fmt.Errorf("chat groups not found: %w", err),
			),
		)
		return []SearchResultItem{}
	}
	if len(chatGroups) == 0 {
		return []SearchResultItem{}
	}
	result := []SearchResultItem{}
	for _, chatGroup := range chatGroups {
		channels, err := chatGroup.ParentChats().AllG(cs.C)
		if err != nil {
			log.Println(
				errorWrapper(
					fmt.Errorf("channels for chat group %q not found: %w", chatGroup.ID, err),
				),
			)
			continue
		}
		if len(channels) == 0 {
			continue
		}
		result = append(
			result,
			SearchResultItem{
				Group:    chatGroup,
				Channels: channels,
			},
		)
	}
	return result
}
func (cs *ChatService) DeleteChatGroup(owner *models.User, chatGroupId string) error {
	errorWrapper := getErrorWrapper("DeleteChatGroup of %q", chatGroupId)
	chatGroup, err := models.Chats(
		models.ChatWhere.ID.EQ(chatGroupId),
		models.ChatWhere.ParentID.IsNull(),
		models.ChatWhere.OwnerID.EQ(null.StringFrom(owner.ID)),
	).OneG(cs.C)
	if err != nil || chatGroup == nil {
		return errorWrapper(ErrChatGroupNotFound)
	}
	_, err = chatGroup.ParentChats().DeleteAllG(cs.C)
	if err != nil {
		return errorWrapper(err)
	}
	_, err = chatGroup.DeleteG(cs.C)
	if err != nil {
		return errorWrapper(err)
	}
	return nil
}
