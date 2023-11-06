package domain

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *Response) sendResponse(c *gin.Context, success bool, status int, message string, data interface{}) {
	response := &Response{
		Success: success,
		Status:  status,
		Message: message,
		Data:    data,
	}

	c.JSON(status, response)
}

func (r *Response) SendSuccess(c *gin.Context, message string, data interface{}) {
	r.sendResponse(c, true, http.StatusOK, message, data)
}

func (r *Response) SendError(c *gin.Context, status int, message string, err error) {
	// Log the error
	fmt.Printf("Status %d: %s\n", status, err.Error())
	// Send the error response to the client
	r.sendResponse(c, false, status, message, nil)
}

func (r *Response) SendErrorAbort(c *gin.Context, status int, message string, err error) {
	// Log the error
	fmt.Printf("Status %d: %s\n", status, err.Error())
	// Send the error response to the client
	r.sendResponse(c, false, status, message, nil)
	c.Abort()
}
