package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/app-service/BasketAPI"
	"github.com/quible-io/quible-api/app-service/RSC"
)

// @Summary		Get Schedule for Season
// @Description	Returns list of games for the selected season
// @Tags			RSC,private
// @Produce		json
// @Param			date	query		string	false	"Sport season" default(<current season>) example(2023)
// @Param			team_id	query		int	false	"Team ID"
// @Success		200	{array}	  RSC.ScheduleItem
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/schedule-season [get]
func ScheduleSeason(c *gin.Context) {
	data, err := RSC.GetScheduleSeason(c.Request.URL.Query())
	if err != nil {
		log.Printf("failed to use ScheduleSeason API: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_ScheduleSeason)
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary		Get Daily Schedule
// @Description	Returns list of games for the next 7 days
// @Tags			RSC,private
// @Produce		json
// @Param			team_id	query		int	false	"Team ID"
// @Param			game_id	query		string	false	"Game ID"
// @Param			date	query		string	false	"Report for date and 7 days in advance" format(date) default(now) example(2023-11-23)
// @Success		200	{array}	  RSC.ScheduleItem
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/daily-schedule [get]
func DailySchedule(c *gin.Context) {
	data, err := RSC.GetDailySchedule(c.Request.URL.Query())
	if err != nil {
		log.Printf("failed to use DailySchedule API: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_DailySchedule)
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary		Get list of teams
// @Description	Returns list of teams or a single team info
// @Tags			RSC,private
// @Produce		json
// @Param			team_id	query		int	false	"Team ID"
// @Success		200	{array}		RSC.TeamInfoItemExtended
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/team-info [get]
func TeamInfo(c *gin.Context) {
	data, err := RSC.GetTeamInfo(c.Request.URL.Query())
	if err != nil {
		log.Printf("failed to use TeamInfo API: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_TeamInfo)
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary		Get team(s) stats
// @Description	Returns teams stats for the selected season
// @Tags			RSC,private
// @Produce		json
// @Param			team_id	query		int	false	"Team ID"
// @Param			date	query		string	false	"Beginning of sport season" default(<current season>) example(2023)
// @Success		200	{array}		RSC.TeamSeasonStatItem
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/team-stats [get]
func TeamStats(c *gin.Context) {
	data, err := RSC.GetTeamStats(c.Request.URL.Query())
	if err != nil {
		log.Printf("failed to use TeamStats API: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_TeamStats)
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary		Get list of players
// @Description	Returns list of players of all or selected team
// @Tags			RSC,private
// @Produce		json
// @Param			team_id	query		int	false	"Team ID"
// @Success		200	{array}		RSC.PlayerInfoItemExtended
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/player-info [get]
func PlayerInfo(c *gin.Context) {
	data, err := RSC.GetPlayerInfo(c.Request.URL.Query())
	if err != nil {
		log.Printf("failed to use PlayerInfo API: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_PlayerInfo)
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary		Get player(s) stats
// @Description	Returns players stats for the selected season
// @Tags			RSC,private
// @Produce		json
// @Param			team_id	query		int	false	"Team ID"
// @Param			player_id	query		int	false	"Player ID"
// @Param			date	query		string	false	"Beginning of sport season" default(<current season>) example(2023)
// @Success		200	{array}		RSC.PlayerSeasonStatItem
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/player-stats [get]
func PlayerStats(c *gin.Context) {
	data, err := RSC.GetPlayerStats(c.Request.URL.Query())
	if err != nil {
		log.Printf("failed to use PlayerStats API: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_PlayerStats)
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary		Get players injuries
// @Description	Returns list of recorded players injuries
// @Tags			RSC,private
// @Produce		json
// @Param			team_id	query		int	false	"Team ID"
// @Success		200	{array}		RSC.InjuryItem
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/injuries [get]
func Injuries(c *gin.Context) {
	data, err := RSC.GetInjuries(c.Request.URL.Query())
	if err != nil {
		log.Printf("failed to use Injuries API: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_Injuries)
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary		Get live feed
// @Description	Returns [live] data on current/past game(s)
// @Tags			RSC,private
// @Produce		json
// @Param			date	query		string	false	"Specific date returns started and finished events from that date" format(date) default(now) example(2023-11-23)
// @Success		200	{array}		RSC.LiveFeedItem
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/live [get]
func LiveFeed(c *gin.Context) {
	data, err := RSC.GetLiveFeed(c.Request.URL.Query())
	if err != nil {
		log.Printf("failed to use LiveFeed API: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_LiveFeed)
		return
	}
	c.JSON(http.StatusOK, data)
}

func LivePush(c *gin.Context) {
	var body any
	if err := c.BindJSON(&body); err != nil {
		log.Printf("unable to parse request body")
	} else {
		encoded, _ := json.MarshalIndent(&body, "", "  ")
		log.Println(string(encoded))
	}
	c.Status(http.StatusOK)
}

// @Summary		Get list of games on a specific date
// @Tags			BasketAPI,private
// @Produce		json
// @Param			date	query		string	true	"Specific date to list games for" format(date) example(2024-01-20)
// @Success		200	{array}		BasketAPI.Game
// @Failure		400	{object}	ErrorResponse
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/games [get]
func GetGames(c *gin.Context) {
	var query GetGamesDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		ErrorMap.SendError(c, http.StatusBadRequest, Err400_MissingRequiredQueryParam)
		return
	}
	games, err := BasketAPI.GetGames(c.Request.Context(), query.Date)
	if err != nil {
		log.Println(err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_BasketAPIGetGames)
		return
	}
	c.JSON(http.StatusOK, games)
}
