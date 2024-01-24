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
)

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

func SearchPublicChannelsByChatGroupTitle(c *gin.Context) {
	q := c.Request.URL.Query().Get("q")
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	c.JSON(http.StatusOK, cs.SearchPublicChannelsByChatGroupTitle(q))
}

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
