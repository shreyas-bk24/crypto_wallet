package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
	Error  string `json:"error"`
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(200, APIResponse{
		Status: "success",
		Data:   data,
		Error:  "",
	})
}

func Failure(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, APIResponse{
		Status: "error",
		Data:   nil,
		Error:  message,
	})
}
