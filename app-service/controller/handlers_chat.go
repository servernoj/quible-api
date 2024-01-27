package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ably/ably-go/ably"
	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/app-service/services/ablyService"
	"github.com/quible-io/quible-api/app-service/services/chatService"
	"github.com/quible-io/quible-api/lib/jwt"
)

// @Summary		Create a chat group owned by the logged in user
// @Tags			chat,private
// @Accept		json
// @Produce		json
// @Param			request	body		chatService.CreateChatGroupDTO	true	"New chat group details"
// @Success		201		{object}	models.Chat
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/chat/groups [post]
func CreateChatGroup(c *gin.Context) {
	var dto chatService.CreateChatGroupDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Println(err)
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}
	user := getUserFromContext(c)
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	chatGroup, err := cs.CreateChatGroup(
		user,
		dto.Name,
		dto.Title,
		dto.Summary,
		dto.IsPrivate,
	)
	if err != nil {
		log.Println(err)
		if errors.Is(err, chatService.ErrChatGroupExists) {
			SendError(c, http.StatusBadRequest, Err400_ChatGroupExists)
		} else {
			SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		}
		return
	}
	c.JSON(http.StatusCreated, chatGroup)
}

// @Summary		Delete chat group (logged in user must be the owner)
// @Tags			chat,private
// @Param     chatGroupId path string true "Chat group ID" format(uuid)
// @Success		204
// @Failure		401		{object}	ErrorResponse
// @Failure		404		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/chat/groups/{chatGroupId} [delete]
func DeleteChatGroup(c *gin.Context) {
	id := c.Param("chatGroupId")
	user := getUserFromContext(c)
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	if err := cs.DeleteChatGroup(user, id); err != nil {
		log.Println(err)
		if errors.Is(err, chatService.ErrChatGroupNotFound) {
			SendError(c, http.StatusNotFound, Err404_ChatGroupNotFound)
		} else {
			SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary		Get new Ably TokenRequest associated with the logged in user
// @Tags			chat,private
// @Produce		json
// @Success		200	{object}	AblyTokenRequest
// @Failure		401	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/chat/token [get]
func GetChatToken(c *gin.Context) {
	user := getUserFromContext(c)
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	capabilities, err := cs.GetCapabilities(user)
	if err != nil {
		log.Println(err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	marshalledCapabilities, _ := json.Marshal(&capabilities)
	token, err := ablyService.CreateTokenRequest(&ably.TokenParams{
		Capability: string(marshalledCapabilities),
		ClientID:   user.ID,
	})
	if err != nil {
		log.Printf("unable to generate ably TokenRequest for user %q: %q", user.ID, err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	c.JSON(http.StatusOK, token)
}

// @Summary		List user's channels and their parent chat groups
// @Tags			chat,private
// @Produce		json
// @Success		200	{array}	chatService.GroupedChannels
// @Failure		401	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/chat/channels [get]
func GetMyGroupedChannels(c *gin.Context) {
	user := getUserFromContext(c)
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	grouppedChannels, err := cs.GetMyGroupedChannels(user)
	if err != nil {
		log.Printf("unable to retrieve grouped channels for user %q: %q", user.ID, err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	c.JSON(http.StatusOK, grouppedChannels)
}

// @Summary		List user's channels
// @Tags			chat,private
// @Produce		json
// @Success		200	{array}	chatService.Channel
// @Failure		401	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/chat/channels/list [get]
func GetMyChannels(c *gin.Context) {
	user := getUserFromContext(c)
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	channels, err := cs.GetMyChannels(user)
	if err != nil {
		log.Printf("unable to retrieve channels for user %q: %q", user.ID, err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	c.JSON(http.StatusOK, channels)
}

// @Summary		List chat groups owned by the logged in user
// @Tags			chat,private
// @Produce		json
// @Success		200		{array}		models.Chat
// @Failure		401		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/chat/groups [get]
func ListChatGroups(c *gin.Context) {
	user := getUserFromContext(c)
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	chatGroups, err := cs.GetChatGroups(user)
	if err != nil {
		log.Println(err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	c.JSON(http.StatusOK, chatGroups)
}

// @Summary		Create a channel within the specified chat group (logged in user must own that chat group)
// @Tags			chat,private
// @Accept		json
// @Produce		json
// @Param			chatGroupId	query		string	false	"ID of the parent `chat group`" format(uuid)
// @Param			request	body		chatService.CreateChannelDTO	true	"New channel details"
// @Success		201		{object}	models.Chat
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/chat/channels [post]
func CreateChannel(c *gin.Context) {
	chatGroupId := c.Request.URL.Query().Get("chatGroupId")
	var dto chatService.CreateChannelDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Println(err)
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	channel, err := cs.CreateChannel(
		chatGroupId,
		dto.Name,
		dto.Title,
		dto.Summary,
	)
	if err != nil {
		log.Println(err)
		if errors.Is(err, chatService.ErrChatGroupNotFound) {
			SendError(c, http.StatusNotFound, Err404_ChatGroupNotFound)
		} else if errors.Is(err, chatService.ErrChannelExists) {
			SendError(c, http.StatusBadRequest, Err400_ChannelExists)
		} else {
			SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		}
		return
	}
	c.JSON(http.StatusCreated, channel)
}

// @Summary		Search public chat groups by partially matched title
// @Tags			chat,public
// @Produce		json
// @Param			q	query		string	false	"partial match for a chat group title, if not provided returns all public results"
// @Success		200		{array}		chatService.SearchResultItem
// @Failure		401		{object}	ErrorResponse
// @Router		/chat/groups/search [get]
func SearchPublicChannelsByChatGroupTitle(c *gin.Context) {
	q := c.Request.URL.Query().Get("q")
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	c.JSON(http.StatusOK, cs.SearchPublicChannelsByChatGroupTitle(q))
}

// @Summary		Join public channel (channel associated with a public chat group)
// @Tags			chat,private
// @Produce		json
// @Param     channelId path string true "Channel ID" format(uuid)
// @Success		200
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		404		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/chat/channels/{channelId} [post]
func JoinPublicChannel(c *gin.Context) {
	id := c.Param("channelId")
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	user := getUserFromContext(c)
	if err := cs.JoinPublicChannel(user, id); err != nil {
		log.Println(err)
		if errors.Is(err, chatService.ErrChannelNotFound) {
			SendError(c, http.StatusNotFound, Err404_ChannelNotFound)
		} else if errors.Is(err, chatService.ErrChatGroupNotFound) {
			SendError(c, http.StatusNotFound, Err404_ChatGroupNotFound)
		} else if errors.Is(err, chatService.ErrPrivateChatGroup) {
			SendError(c, http.StatusBadRequest, Err400_ChatGroupIsPrivate)
		} else if errors.Is(err, chatService.ErrSelfOwnedChatGroup) {
			SendError(c, http.StatusBadRequest, Err400_ChatGroupIsSelfOwned)
		} else if errors.Is(err, chatService.ErrChannelAlreadyJoined) {
			SendError(c, http.StatusBadRequest, Err400_ChannelAlreadyJoined)
		} else {
			SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		}
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Leave previously joined channel
// @Tags			chat,private
// @Param     channelId path string true "Channel ID" format(uuid)
// @Success		204
// @Failure		401		{object}	ErrorResponse
// @Failure		404		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/chat/channels/{channelId} [delete]
func LeaveChannel(c *gin.Context) {
	id := c.Param("channelId")
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	user := getUserFromContext(c)
	if err := cs.LeaveChannel(user, id); err != nil {
		log.Println(err)
		if errors.Is(err, chatService.ErrChannelNotFound) {
			SendError(c, http.StatusNotFound, Err404_ChannelNotFound)
		} else {
			SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary		Invite a user to join private channel
// @Tags			chat,private
// @Param     channelId path string true "Channel ID" format(uuid)
// @Param			request	body		chatService.InviteToPrivateChannelDTO	true	"Contains invitee's email"
// @Success		200
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		404		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/chat/channels/{channelId}/invite [post]
func InviteToPrivateChannel(c *gin.Context) {
	channelId := c.Param("channelId")
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	user := getUserFromContext(c)
	var dto chatService.InviteToPrivateChannelDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Println(err)
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}
	if err := cs.InviteToPrivateChannel(dto.Email, user, channelId); err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, chatService.ErrChannelNotFound):
			SendError(c, http.StatusNotFound, Err404_ChannelNotFound)
		case errors.Is(err, chatService.ErrPublicChatGroup):
			SendError(c, http.StatusBadRequest, Err400_ChatGroupIsPublic)
		case errors.Is(err, chatService.ErrInvalidInviteeEmail):
			SendError(c, http.StatusBadRequest, Err400_EmailNotFound)
		default:
			SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		}
		return
	}
	c.Status(http.StatusOK)
}

func AcceptInvitationToPrivateChannel(c *gin.Context) {
	var dto chatService.TokenDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		SendError(c, http.StatusBadRequest, Err400_InvalidOrMalformedToken)
		return
	}
	tokenClaims, err := jwt.VerifyJWT(dto.Token, jwt.TokenActionInvitationToPrivateChat)
	if err != nil {
		log.Printf("invalid token: %q", err)
		SendError(c, http.StatusExpectationFailed, Err417_InvalidToken)
		return
	}
	userId := tokenClaims["userId"].(string)
	extraClaims := tokenClaims["extraClaims"].(jwt.ExtraClaims)
	inviteeId, ok := extraClaims["inviteeId"].(string)
	if !ok {
		log.Printf("missing inviteeId in extraClaims")
		SendError(c, http.StatusExpectationFailed, Err417_InvalidToken)
		return
	}
	channelId, ok := extraClaims["channelId"].(string)
	if !ok {
		log.Printf("missing channelId in extraClaims")
		SendError(c, http.StatusExpectationFailed, Err417_InvalidToken)
		return
	}
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	if err := cs.AcceptInvitationToPrivateChannel(inviteeId, userId, channelId); err != nil {
		log.Println(err)
		switch {
		default:
			SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		}
		return
	}
	c.Status(http.StatusOK)
}
