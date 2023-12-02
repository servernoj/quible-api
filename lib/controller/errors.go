package controller

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorCode interface {
	~int
	String() string
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorMap[T ErrorCode] map[int]map[T]string

func (errorMap ErrorMap[T]) SendError(c *gin.Context, status int, code T) {
	if _, ok := errorMap[status][code]; !ok {
		log.Printf("unable to find mapping for [%d]%d", status, code)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.AbortWithStatusJSON(
		status,
		ErrorResponse{
			Code:    int(code),
			Message: errorMap[status][code],
		},
	)
}

// @Summary		Error Codes UI
// @Description	Renders the list of erros reported by the microservice
// @Tags			docs
// @Produce		text/html
// @Success		200	{string} string
// @Router		/docs/errors [get]
func (errorMap ErrorMap[T]) GetErrorCodes(c *gin.Context) {
	statuses := sort.IntSlice(make([]int, len(errorMap)))
	idx := 0
	for status := range errorMap {
		statuses[idx] = status
		idx++
	}
	sort.Sort(statuses)
	builder := strings.Builder{}
	builder.WriteString("<h2>Error codes</h2>")
	for _, httpStatus := range statuses {
		statusMap := errorMap[httpStatus]
		builder.WriteString(fmt.Sprintf("<h3>Status code: %d</h3>", httpStatus))
		builder.WriteString("<pre>")

		// sort statusMap keys
		errorCodes := sort.IntSlice(make([]int, len(statusMap)))
		idx := 0
		for errorCode := range statusMap {
			errorCodes[idx] = int(errorCode)
			idx++
		}
		sort.Sort(errorCodes)
		for _, errorCode := range errorCodes {
			errorCode := T(errorCode)
			message := statusMap[errorCode]
			builder.WriteString(
				fmt.Sprintf("%-50s%-20d%s\n", errorCode.String(), errorCode, message),
			)
		}
		builder.WriteString("</pre>")
	}
	c.Writer.Header().Add("content-type", gin.MIMEHTML)
	c.String(http.StatusOK, builder.String())
}
