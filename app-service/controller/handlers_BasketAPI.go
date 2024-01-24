package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ably/ably-go/ably"
	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/app-service/services/BasketAPI"
	"github.com/quible-io/quible-api/app-service/services/ablyService"
)

// @Summary		Get list of games on a specific date
// @Tags			BasketAPI,private
// @Produce		json
// @Param			date	query		string	true	"Specific date to list games for" format(date) example(2024-01-20)
// @Param			localTimeZoneShift	query		string	false	"Local TZ shift (to UTC) in hours to relate game start time. Defaults to EST/EDT timezone" example(-7)
// @Success		200	{array}		BasketAPI.Game
// @Failure		400	{object}	ErrorResponse
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/games [get]
func GetGames(c *gin.Context) {
	var query BasketAPI.GetGamesDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Printf("unable to parse query: %s", err)
		SendError(c, http.StatusBadRequest, Err400_MissingRequiredQueryParam)
		return
	}
	games, err := BasketAPI.GetGames(c.Request.Context(), query)
	if err != nil {
		log.Println(err)
		SendError(c, http.StatusFailedDependency, Err424_BasketAPIGetGames)
		return
	}
	c.JSON(http.StatusOK, games)
}

// @Summary		Get game details
// @Tags			BasketAPI,private
// @Produce		json
// @Param			gameId	query		string	true	"ID of the BasketAPI match"
// @Success		200	{array}		BasketAPI.GameDetails
// @Failure		400	{object}	ErrorResponse
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/game [get]
func GetGameDetails(c *gin.Context) {
	var query BasketAPI.GetGameDetailsDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Printf("unable to parse query: %s", err)
		SendError(c, http.StatusBadRequest, Err400_MissingRequiredQueryParam)
		return
	}
	result, err := BasketAPI.GetGameDetails(c.Request.Context(), c.Request.URL.Query())
	if err != nil {
		log.Println(err)
		SendError(c, http.StatusFailedDependency, Err424_BasketAPIGetGameDetails)
		return
	}
	c.JSON(http.StatusOK, result)
}

func GetLiveToken(c *gin.Context) {
	capabilities, _ := json.Marshal(&map[string][]string{
		"live:main": {"subscribe", "history"},
	})
	token, err := ablyService.CreateTokenRequest(&ably.TokenParams{
		Capability: string(capabilities),
		ClientID:   "nobody",
	})
	if err != nil {
		log.Printf("unable to generate ably TokenRequest: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	c.JSON(http.StatusOK, token)
}
