package response

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Status:  "error",
		Message: message,
	})
}
