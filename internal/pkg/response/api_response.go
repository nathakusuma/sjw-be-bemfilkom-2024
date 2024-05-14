package response

import "github.com/gin-gonic/gin"

type ApiResponse struct {
	Status  int
	Message string
	Data    any
}

func NewApiResponse(status int, message string, data any) ApiResponse {
	return ApiResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func (r ApiResponse) Send(ctx *gin.Context) {
	ctx.JSON(r.Status, gin.H{
		"message": r.Message,
		"data":    r.Data,
	})
}
