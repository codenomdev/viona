package response

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrInternalServerError = "internal server error"
	ErrNotFound            = "not found"
	ErrBadRequest          = "bad request"
	ErrForbidden           = "forbidden"
	ErrUnauthorized        = "unauthorized"
	ErrTooManyRequest      = "too many request"
	ErrConflict            = "conflict"
	ErrUnprocessEntity     = "unprocessable entity"
)

type RestError interface {
	error
	ResponseMeta() MetaResponse
	ResponseCode() int
	ResponseErrors() []string
	ResponseMessage() string
	ResponsePayload() interface{}
}

func (e ErrorResponse) ResponseCode() int            { return e.Meta.StatusCode }
func (e ErrorResponse) ResponseErrors() []string     { return e.Errors }
func (e ErrorResponse) ResponseMessage() string      { return e.Meta.Message }
func (e ErrorResponse) ResponsePayload() interface{} { return e.Payload }
func (e ErrorResponse) ResponseMeta() MetaResponse   { return e.Meta }
func (e ErrorResponse) Error() string {
	return fmt.Sprintf("status: %d - message: %s - errors: %+v", e.ResponseCode(), e.ResponseMessage(), e.Errors)
}

func NewHttpError(code int, message string, errors []string, payload interface{}) RestError {
	return ErrorResponse{
		Meta: MetaResponse{
			Success:    false,
			Message:    message,
			StatusCode: code,
		},
		Errors:  errors,
		Payload: payload,
	}
}

// Parse error for security response...
func ParseHttpError(err error) RestError {
	switch {
	case strings.Contains(err.Error(), "SQLSTATE"):
		return NewParseSQLError(err)
	case strings.Contains(err.Error(), "record not found"):
		return NewHttpNotFound([]string{err.Error()}, nil)
	case strings.Contains(err.Error(), "strconv.ParseInt"):
		return NewHttpBadRequest([]string{"parameters is bad request"}, nil)
	case strings.Contains(err.Error(), "invalid host"):
		return NewHttpBadRequest([]string{"Host invalid"}, nil)
	case strings.Contains(err.Error(), "UUID"):
		return NewHttpBadRequest([]string{"format ID not valid. Must be UUID v4"}, nil)
	case strings.Contains(err.Error(), "invalid token") || strings.Contains(err.Error(), "token expired"):
		return NewHttpUnauthorized([]string{"token invalid or token expired"}, nil)
	case strings.Contains(err.Error(), "can not parse token") || strings.Contains(err.Error(), "unexpected signing method"):
		return NewHttpBadRequest([]string{err.Error()}, nil)
	case strings.Contains(err.Error(), "strconv.ParseInt"):
		return NewHttpBadRequest([]string{"format invalid, param must int"}, nil)
	case strings.Contains(err.Error(), "session"):
		return NewHttpInternalServerError("Session store something errors")
	case strings.Contains(err.Error(), "ParseUint"):
		return NewHttpBadRequest([]string{"ID must a be int"}, nil)
	default:
		if restErr, ok := err.(ErrorResponse); ok {
			return restErr
		}
		return NewHttpInternalServerError(err.Error())
	}
}

func NewParseSQLError(err error) RestError {
	switch {
	case strings.Contains(err.Error(), "23505"):
		return NewParseSQLUniqueError(err)
	case strings.Contains(err.Error(), "22P02"):
		return NewHttpBadRequest([]string{"parameters is bad request"}, nil)
	default:
		if rest, ok := err.(ErrorResponse); ok {
			return rest
		}
		return NewHttpInternalServerError(err.Error())
	}
}

func NewParseSQLUniqueError(err error) RestError {
	switch {
	case strings.Contains(err.Error(), "email"):
		return NewHttpBadRequest([]string{"email already exists"}, nil)
	case strings.Contains(err.Error(), "phone"):
		return NewHttpBadRequest([]string{"phone number already exists"}, nil)
	default:
		return NewHttpBadRequest([]string{"must unique", "key already exists"}, nil)
	}
}

func NewHttpBadRequest(errors []string, payload interface{}) RestError {
	return NewHttpError(http.StatusBadRequest, ErrBadRequest, errors, payload)
}

func NewHttpValidationError(payload interface{}) RestError {
	return NewHttpError(http.StatusBadRequest, "validation error", []string{"validation error"}, payload)
}

func NewHttpNotFound(errors []string, payload interface{}) RestError {
	return NewHttpError(http.StatusNotFound, ErrNotFound, errors, payload)
}

func NewHttpInternalServerError(message string) RestError {
	return NewHttpError(http.StatusInternalServerError, ErrInternalServerError, []string{message}, nil)
}

func NewHttpUnauthorized(errors []string, payload interface{}) RestError {
	return NewHttpError(http.StatusUnauthorized, ErrUnauthorized, errors, payload)
}

func NewHttpTooManyRequest(errors []string) RestError {
	return NewHttpError(http.StatusTooManyRequests, ErrTooManyRequest, errors, nil)
}

func NewHttpConflict(errors []string) RestError {
	return NewHttpError(http.StatusConflict, ErrConflict, errors, nil)
}

func NewHttpUnprocessedEntity(errors []string) RestError {
	return NewHttpError(http.StatusUnprocessableEntity, ErrUnprocessEntity, errors, nil)
}

func NewHttpForbidden(errors []string) RestError {
	return NewHttpError(http.StatusForbidden, ErrForbidden, errors, nil)
}
