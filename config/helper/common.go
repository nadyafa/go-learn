package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// dev failed response
type ErrorStruct struct {
	Err  error
	Code int
	Msg  string
}

var Validate = validator.New()

func (e *ErrorStruct) ErrorJSON(message string) gin.H {
	return gin.H{
		"status":  e.Code,
		"message": message,
		"detail":  e.Err.Error(),
	}
}
