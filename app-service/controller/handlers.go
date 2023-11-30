package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/app-service/RSC"
)

// @Summary		Get Schedule for Season
// @Description	Returns list of games for the selected season
// @Tags			RSC,private
// @Produce		json
// @Success		200	{array}	  RSC.ScheduleItem
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/schedule-season [get]
func ScheduleSeason(c *gin.Context) {
	data, err := RSC.NewClient().GetScheduleSeason()
	if err != nil {
		log.Printf("failed to use ScheduleSeason API: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_UnknownError)
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
	data, err := RSC.NewClient().GetDailySchedule(c.Request.URL.Query())
	if err != nil {
		log.Printf("failed to use DailySchedule API: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_UnknownError)
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary		Get list of teams
// @Description	Returns list of teams or a single team info
// @Tags			RSC,private
// @Produce		json
// @Param			team_id	query		int	false	"Team ID"
// @Success		200	{array}		RSC.TeamInfoItem
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/team-info [get]
func TeamInfo(c *gin.Context) {
	data, err := RSC.NewClient().GetTeamInfo(c.Request.URL.Query())
	if err != nil {
		log.Printf("failed to use TeamInfo API: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_UnknownError)
		return
	}
	c.JSON(http.StatusOK, data)
}
