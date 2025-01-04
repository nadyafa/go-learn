package helper

// user failed response
type Response struct {
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
}

func SuccessResponse(message string, data any, statusCode int) Response {
	return Response{
		Data:       data,
		Message:    message,
		StatusCode: statusCode,
	}
}

func FailedResponse(message string, statusCode int) Response {
	return Response{
		Message:    message,
		StatusCode: statusCode,
	}
}
