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
// @Success		200	{object}	[]RSC.ScheduleItem
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/schedule-season [get]
func ScheduleSeason(c *gin.Context) {
	if c.IsAborted() {
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	data, err := RSC.NewClient().GetScheduleSeason()
	if err != nil {
		log.Printf("failed to use ScheduleSeason API: %q", err)
		SendError(c, http.StatusFailedDependency, Err424_UnknownError)
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
// @Success		200	{object}	[]RSC.ScheduleItem
// @Failure		401	{object}	ErrorResponse
// @Failure		424	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/daily-schedule [get]
func DailySchedule(c *gin.Context) {
	if c.IsAborted() {
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	data, err := RSC.NewClient().GetDailySchedule(c.Request.URL.Query())
	if err != nil {
		log.Printf("failed to use DailySchedule API: %q", err)
		SendError(c, http.StatusFailedDependency, Err424_UnknownError)
		return
	}
	c.JSON(http.StatusOK, data)
}
