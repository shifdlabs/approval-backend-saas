package utils

import (
	"Microservice/helper"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	responseData := Response{
		Success: true,
		Code:    200,
		Message: "Success",
		Data:    data,
	}
	c.JSON(200, responseData)
}

func ErrorResponse(c *gin.Context, data helper.ErrorModel) {
	responseData := Response{
		Success: false,
		Code:    data.Code,
		Message: data.Message, // Only the status text in the response message
		Data:    nil,
	}

	c.JSON(data.Code, responseData)
}

func ConvertErrorToCustomError(err error, code int, message string) *helper.CustomError {
	fileName, atLine := helper.GetFileAndLine(err)
	return &helper.CustomError{
		Code:     code,
		Message:  message,
		FileName: fileName,
		AtLine:   atLine,
	}
}
