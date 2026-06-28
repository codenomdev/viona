package response

import (
	"fmt"
	"net/http"
)

type ResponseSuccess interface {
	ResponseCode() int
	ResponseData() interface{}
	ResponseErrors() bool
	ResponseMessage() string
	ResponsePayload() interface{}
	ResponseMeta() MetaResponse
}

func (e SuccessResponse) ResponseCode() int            { return e.Meta.StatusCode }
func (e SuccessResponse) ResponseData() interface{}    { return e.Data }
func (e SuccessResponse) ResponseErrors() bool         { return e.Error }
func (e SuccessResponse) ResponseMessage() string      { return e.Meta.Message }
func (e SuccessResponse) ResponsePayload() interface{} { return e.Payload }
func (e SuccessResponse) ResponseMeta() MetaResponse   { return e.Meta }

func (e SuccessResponse) String() string {
	return fmt.Sprintf("status: %d - message: %s - data: %+v", e.Meta.StatusCode, e.Meta.Message, e.Data)
}

func NewHttpOK(data interface{}) ResponseSuccess {
	return SuccessResponse{
		Meta: MetaResponse{
			Success:    true,
			Message:    "success",
			StatusCode: http.StatusOK,
		},
		Error:   false,
		Data:    data,
		Payload: nil,
	}
}

func NewHttpCreated(data interface{}) ResponseSuccess {
	return SuccessResponse{
		Meta: MetaResponse{
			Success:    true,
			Message:    "created",
			StatusCode: http.StatusCreated,
		},
		Error:   false,
		Data:    data,
		Payload: nil,
	}
}
