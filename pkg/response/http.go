package response

import "net/http"

type MetaResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type SuccessResponse struct {
	Meta    MetaResponse `json:"meta"`
	Error   bool         `json:"error"`
	Payload interface{}  `json:"payload,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
}

type ErrorResponse struct {
	Meta    MetaResponse `json:"meta"`
	Errors  []string     `json:"errors"`
	Payload interface{}  `json:"payload,omitempty"`
}

func ParseHttpResponse(resp interface{}) (int, interface{}) {
	switch v := resp.(type) {
	case SuccessResponse:
		return v.Meta.StatusCode, v
	case ErrorResponse:
		return v.Meta.StatusCode, v
	default:
		return http.StatusOK, resp
	}
}
