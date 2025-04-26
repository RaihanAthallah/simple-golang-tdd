package utils

import (
	dto "simple-golang-tdd/dto"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, status int, message string) {

	response := dto.ErrorResponse{
		Status:  status,
		Message: message,
	}

	c.JSON(status, response)
}

func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {

	response := dto.SuccessResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	c.JSON(status, response)
}
