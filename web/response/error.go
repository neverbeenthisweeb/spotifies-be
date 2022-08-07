package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func Error(c *gin.Context, err error) {
	status := http.StatusBadRequest // Default status
	obj := ErrorResponse{
		Message: err.Error(),
	}
	c.JSON(status, obj)
}
