package httpx

import "fmt"

type HttpResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
}
type HttpPaginationResponse struct {
	HttpResponse
	Meta PaginationMeta `json:"meta,omitempty"`
}

func (e HttpResponse) Error() string {
	return fmt.Sprintf("description: %s", e.Message)
}

func NewHttpResponse[T any](statusCode int, message string, data T) HttpResponse {
	return HttpResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

func NewHttpPaginationResponse[T any](statusCode int, message string, data T, meta PaginationMeta) HttpPaginationResponse {
	return HttpPaginationResponse{
		HttpResponse{
			StatusCode: statusCode,
			Message:    message,
			Data:       data,
		},
		meta,
	}
}
