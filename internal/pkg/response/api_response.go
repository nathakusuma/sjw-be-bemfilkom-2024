package response

import (
	"github.com/gin-gonic/gin"
	"log"
)

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
	if r.Status >= 500 {
		if gin.Mode() == gin.ReleaseMode {
			r.Data = gin.H{}
		}
		log.Printf("ERROR %d (%s): %v\n", r.Status, r.Message, r.Data)
	}
	ctx.JSON(r.Status, gin.H{
		"message": r.Message,
		"data":    r.Data,
	})
}
