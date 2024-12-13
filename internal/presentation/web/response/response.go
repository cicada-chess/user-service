package response

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}

func NewSuccessResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, SuccessResponse{message})
}
