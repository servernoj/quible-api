package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/app-service/RSC"
)

// @Summary		Get Schedule for Season
// @Description	Returns list of gamnes for the selected season
// @Tags			RSC,private
// @Produce		json
// @Success		200	{object}	[]RSC.ScheduleSeasonItem
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
		SendError(c, http.StatusFailedDependency, Err424_UnknownError)
		return
	}
	c.JSON(http.StatusOK, data)
}
