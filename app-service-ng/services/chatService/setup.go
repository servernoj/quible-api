package chatService

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/quible-io/quible-api/app-service-ng/services/emailService"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/jwt"
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

func getErrorWrapper(prefixFormat string, prefixArgs ...any) func(error, ...error) error {
	return func(mainError error, secondaryErrors ...error) error {
		prefix := fmt.Sprintf(prefixFormat, prefixArgs...)
		if len(secondaryErrors) > 0 {
			return fmt.Errorf("%s:\n%w %v", prefix, mainError, secondaryErrors)
		} else {
			return fmt.Errorf("%s:\n%w", prefix, mainError)
		}
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
	chatUsers, err := user.ChatUsers(
		models.ChatUserWhere.Disabled.EQ(false),
	).AllG(cs.C)
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
	// 0. Initialize empty lists and sets
	result := []GroupedChannels{}
	chatGroupsMap := make(map[string]*models.Chat)
	myChannels := make(map[string]struct{})
	// 1. Get all chat groups owned by the user
	ownedChatGroups, err := cs.GetChatGroups(user)
	if err != nil {
		return nil, errorWrapper(ErrUnableToRetrtieveChatGroups, err)
	}
	for _, chatGroup := range ownedChatGroups {
		chatGroupsMap[chatGroup.ID] = chatGroup
	}
	// 2. Identify chat groups which user has joined
	chatUsers, err := user.ChatUsers(
		models.ChatUserWhere.Disabled.EQ(false),
	).AllG(cs.C)
	if err != nil {
		return nil, errorWrapper(ErrUnableToRetrieveChatUserRecords, err)
	}
	for _, chatUser := range chatUsers {
		chat, err := models.FindChatG(cs.C, chatUser.ChatID)
		if err != nil || chat == nil {
			log.Println(errorWrapper(ErrUnableToRetrtieveChannel, err))
			continue
		}
		myChannels[chat.ID] = struct{}{}
		parentChatGroup, _ := chat.Parent().OneG(cs.C)
		if parentChatGroup != nil {
			chatGroupsMap[parentChatGroup.ID] = parentChatGroup
		}
	}
	// 3. For every chat group get list of associated channels
	for _, chatGroup := range chatGroupsMap {
		chats, err := models.Chats(
			models.ChatWhere.ParentID.EQ(null.StringFrom(chatGroup.ID)),
		).AllG(cs.C)
		if err != nil {
			log.Println(errorWrapper(ErrUnableToRetrtieveChannels, err))
			continue
		}
		channels := make([]Channel, 0, len(chats))
		for _, chat := range chats {
			// -- we add channel to the list if
			// a) the group is owned by user
			// b) if not, then channelId should match on of the joined channels
			if _, ok := myChannels[chat.ID]; ok || chatGroup.OwnerID.String == userId {
				channels = append(channels, Channel{
					ID:       chat.ID,
					Title:    chat.Title,
					Resource: chatGroup.Resource + ":" + chat.Resource,
					ReadOnly: false,
				})
			}
		}
		result = append(result, GroupedChannels{
			ID:       chatGroup.ID,
			Title:    chatGroup.Title,
			Summary:  chatGroup.Summary.Ptr(),
			Channels: channels,
		})
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
func (cs *ChatService) InviteToPrivateChannel(inviteeEmail string, invitor *models.User, channelId string) error {
	if invitor == nil {
		return ErrUserUndefined
	}
	errorWrapper := getErrorWrapper("InviteToPrivateChannel %q for %q by user %q", channelId, inviteeEmail, invitor.ID)
	// 1. Find the channel record targeted for invitation
	channel, err := models.FindChatG(cs.C, channelId)
	if err != nil || channel == nil {
		return errorWrapper(ErrChannelNotFound)
	}
	// 2. Detect if the parent chat group is actually private
	chatGroup, err := channel.Parent().OneG(cs.C)
	if err != nil || chatGroup == nil {
		return errorWrapper(ErrUnableToRetrtieveChatGroup)
	}
	if !chatGroup.IsPrivate.Bool {
		return errorWrapper(ErrPublicChatGroup)
	}
	// 3. Find invitee user by provided email
	invitee, err := models.Users(
		models.UserWhere.Email.EQ(inviteeEmail),
	).OneG(cs.C)
	if err != nil || invitee == nil {
		return errorWrapper(ErrInvalidInviteeEmail)
	}
	// 4. Test if association between user and channel already exists
	foundChatUser, err := models.FindChatUserG(cs.C, channelId, invitee.ID)
	if err != nil && foundChatUser != nil {
		return errorWrapper(ErrUnableToRetrieveChatUserRecords, err)
	}
	if foundChatUser != nil {
		// refresh it
		foundChatUser.Disabled = true
		if _, err := foundChatUser.UpdateG(cs.C, boil.Infer()); err != nil {
			return errorWrapper(ErrUnableToRefreshChatUserAssociation, err)
		}
	} else {
		// create new one
		chatUser := models.ChatUser{
			ChatID:   channelId,
			UserID:   invitee.ID,
			Disabled: true,
		}
		if err := chatUser.InsertG(cs.C, boil.Infer()); err != nil {
			return errorWrapper(ErrUnableToCreateChatUser)
		}
	}
	// 5. Send invitation email
	token, _ := jwt.GenerateToken(
		invitor,
		jwt.TokenActionInvitationToPrivateChat,
		jwt.ExtraClaims{
			"inviteeId": invitee.ID,
			"channelId": channelId,
		},
	)
	var html bytes.Buffer
	emailService.InviteToPrivateChatGroup(
		invitee.FullName,
		channel.Title,
		chatGroup.Title,
		fmt.Sprintf(
			"%s/forms/accept-private-chat-invitation?token=%s",
			os.Getenv("WEB_CLIENT_URL"),
			token.Token,
		),
		&html,
	)
	if err := email.Send(cs.C, email.EmailDTO{
		From:     "no-reply@quible.io",
		To:       invitee.Email,
		Subject:  "Invitation to join private chat channel",
		HTMLBody: html.String(),
	}); err != nil {
		return errorWrapper(ErrUnableToSendInvitationEmail, err)
	}
	return nil
}
func (cs *ChatService) AcceptInvitationToPrivateChannel(inviteeId string, invitorId string, channelId string) error {
	errorWrapper := getErrorWrapper(
		"AcceptInvitationToPrivateChannel %q for %q created by user %q",
		channelId,
		inviteeId,
		invitorId,
	)
	channel, err := models.FindChatG(cs.C, channelId)
	if err != nil || channel == nil || channel.ParentID.Ptr() == nil {
		return errorWrapper(ErrChannelNotFound)
	}
	chatGroup, err := models.FindChatG(cs.C, channel.ParentID.String)
	if err != nil || chatGroup == nil {
		return errorWrapper(ErrChatGroupNotFound)
	}
	if chatGroup.OwnerID.String != invitorId {
		return errorWrapper(ErrWrongInvitor)
	}
	chatUser, err := models.FindChatUserG(cs.C, channelId, inviteeId)
	if err != nil || chatUser == nil {
		return errorWrapper(ErrChatUserNotFound)
	}
	chatUser.Disabled = false
	if _, err := chatUser.UpdateG(cs.C, boil.Whitelist("disabled")); err != nil {
		return errorWrapper(err)
	}
	return nil
}
