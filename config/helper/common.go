package helper

import (
	"github.com/gin-gonic/gin"
)

// dev failed response
type ErrorStruct struct {
	Err  error
	Code int
}

func (e *ErrorStruct) ErrorJSON() gin.H {
	return gin.H{
		"error":  e.Err.Error(),
		"status": e.Code,
	}
}
