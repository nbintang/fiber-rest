package httpx

import (
	"fmt"
	"time"
)

type HttpResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	TimeStamp  int    `json:"timestamp,omitempty"`
}

func (e HttpResponse) Error() string {
	return fmt.Sprintf("description: %s", e.Message)
}

func NewHttpResponse(statusCode int, message string, data any) HttpResponse {
	return HttpResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		TimeStamp:  int(time.Now().Unix()),
	}
}
