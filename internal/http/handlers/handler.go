package handlers

import (
	"net/http"
)

type Response struct {
	Status     string  `json:"status"`
	StatusCode int     `json:"status_code"`
	Message    *string `json:"message"`
	Result     any     `json:"result"`
}

func makeHTTPResponse(statusCode int, message string, result any) (int, Response) {
	resp := Response{
		Status:     http.StatusText(statusCode),
		StatusCode: statusCode,
		Message:    &message,
		Result:     result,
	}

	if message == "" {
		resp.Message = nil
	}

	return statusCode, resp
}
