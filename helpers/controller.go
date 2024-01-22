package helpers

import (
	"github.com/gin-gonic/gin"
)

type APISuccess struct {
	Message *string `json:"message"`
	Data    *any    `json:"data"`
	Status  *int    `json:"status"`
}

type APIFailure struct {
	Message *string `json:"message"`
	Status  *int    `json:"status"`
}

func StrPtr(str string) *string {
	return &str
}
func AnyPtr(data any) *any {
	return &data
}
func IntPtr(value int) *int {
	return &value
}

func FailureResponse(context *gin.Context, error APIFailure) {
	if error.Status == nil {
		error.Status = IntPtr(400)
	}
	if error.Message == nil {
		error.Message = StrPtr("An Error Occurred")
	}
	context.JSON(*error.Status, gin.H{"message": *error.Message})
	context.Abort()
}

func SuccessResponse(context *gin.Context, response APISuccess) {
	if response.Status == nil {
		response.Status = IntPtr(200)
	}
	if response.Message == nil {
		response.Message = StrPtr("Operation is successful")
	}
	if response.Data == nil {
		response.Data = AnyPtr(map[string]string{})
	}
	context.JSON(*response.Status, gin.H{"message": *response.Message,
		"data": *response.Data})
}
