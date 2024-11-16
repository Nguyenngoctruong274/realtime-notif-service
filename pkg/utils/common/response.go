package common

import (
	"encoding/json"
	"net/http"
	"yes4all/ads-noti-api/pkg/utils/errors"
	"yes4all/ads-noti-api/pkg/utils/tracking"

	"github.com/gin-gonic/gin"
)

type Response struct {
	IsError    *bool       `json:"is_error,omitempty"`
	StatusCode int         `json:"status_code"`
	Code       int         `json:"code"`
	Data       interface{} `json:"data,omitempty"`
	Message    *string     `json:"message,omitempty"`
	TrackID    string      `json:"trace_id"`
}

func (r *Response) String() string {
	data, _ := json.Marshal(r) //nolint:errchkjson
	return string(data)
}

func NewResponse(c *gin.Context, statusCode int, data interface{}, message *string) *Response {
	return &Response{
		Data:       data,
		Message:    message,
		Code:       errors.ErrorMap[statusCode],
		StatusCode: statusCode,
		TrackID:    tracking.GetTrackIDFromContext(c),
	}
}

func NewSuccessResponse(c *gin.Context, data interface{}) *Response {
	return NewResponse(c, http.StatusOK, data, nil)
}

func NewErrorResponse(c *gin.Context, err AppError) *Response {
	return NewResponse(c, err.StatusCode, nil, &err.Message)
}

func JSONOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, NewSuccessResponse(c, data))
}

func SuccessResponseWithError(c *gin.Context, err error) *Response {
	errString := err.Error()
	resp := NewResponse(c, http.StatusOK, nil, &errString)
	isError := true
	resp.IsError = &isError
	return resp
}

func JSONOKWithError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, SuccessResponseWithError(c, err))
}
