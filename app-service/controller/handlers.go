package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/app-service/RSC"
)

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
