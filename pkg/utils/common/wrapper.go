package common

import (
	"errors"
	"fmt"
	"net/http"
	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/utils/constants"
	uerr "yes4all/ads-noti-api/pkg/utils/errors"
	"yes4all/ads-noti-api/pkg/utils/tracking"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AppError struct {
	// We don't show root cause to clients.
	RootCause  error  `json:"-"`
	Code       int    `json:"code"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	TraceID    string `json:"trace_id"`
}

func (e AppError) Error() string {
	return e.Message
}

func (e AppError) Status() int {
	return e.StatusCode
}

func NewAppError(err error, statusCode int, messages ...string) AppError {
	var (
		message string
		myErr   = err
	)
	if err != nil {
		message = err.Error()
		if len(messages) > 0 {
			message = messages[0]
		}

		if ae, ok := err.(AppError); ok {
			ae.Message = message
			return ae
		}
	} else {
		if len(messages) > 0 {
			message = messages[0]
		}
		myErr = errors.New(message)
	}

	return AppError{
		RootCause:  myErr,
		Code:       uerr.ErrorMap[statusCode],
		StatusCode: statusCode,
		Message:    message,
	}
}

// nolint: lll
var (
	ErrConflict       = NewConflictError(errors.New(constants.ConflictErrMess), constants.ConflictErrMess)
	ErrBadRequest     = NewBadRequestError(errors.New(constants.BadRequestErrMess), constants.BadRequestErrMess)
	ErrNotFetchData   = NewBadRequestError(errors.New(constants.NotFetchDataErrMess), constants.NotFetchDataErrMess)
	ErrNotFound       = NewNotFoundError(errors.New(constants.NotFoundErrMess), constants.NotFoundErrMess)
	ErrUnAuthorized   = NewUnAuthorizedError(errors.New(constants.UnAuthorizedErrMess), constants.UnAuthorizedErrMess)
	ErrForbidden      = NewForbiddenError(errors.New(constants.ForbiddenErrMess), constants.ForbiddenErrMess)
	ErrInternalServer = NewInternalServerError(errors.New(constants.InternalServerErrMess), constants.InternalServerErrMess)
)

func NewGenericError(err error, statusCode int, messages ...string) AppError {
	return NewAppError(err, statusCode, messages...)
}

func NewGenericErrorWithMessage(message string, statusCode int) AppError {
	return NewAppError(errors.New(message), statusCode, message)
}

func NewBadRequestError(err error, messages ...string) AppError {
	return NewAppError(err, http.StatusBadRequest, messages...)
}

func NewInternalServerError(err error, messages ...string) AppError {
	return NewAppError(err, http.StatusInternalServerError, messages...)
}

func NewNotFoundError(err error, messages ...string) AppError {
	return NewAppError(err, http.StatusNotFound, messages...)
}

func NewConflictError(err error, messages ...string) AppError {
	return NewAppError(err, http.StatusConflict, messages...)
}

func NewForbiddenError(err error, messages ...string) AppError {
	return NewAppError(err, http.StatusForbidden, messages...)
}

func NewUnAuthorizedError(err error, messages ...string) AppError {
	return NewAppError(err, http.StatusUnauthorized, messages...)
}

func NewUnprocessableEntityError(err error, messages ...string) AppError {
	return NewAppError(err, http.StatusUnprocessableEntity, messages...)
}

func NewBadRequestErrorWithMessage(message string) AppError {
	return NewAppError(errors.New(message), http.StatusBadRequest, message)
}

func NewUnprocessableEntityWithMessage(message string) AppError {
	return NewAppError(errors.New(message), http.StatusUnprocessableEntity, message)
}

func NewInternalServerErrorWithMessage(message string) AppError {
	err := errors.New(message)
	return NewAppError(err, http.StatusInternalServerError, message)
}

func NewNotFoundErrorWithMessage(message string) AppError {
	err := errors.New(message)
	return NewAppError(err, http.StatusNotFound, message)
}

func NewConflictErrorWithMessage(message string) AppError {
	err := errors.New(message)
	return NewAppError(err, http.StatusConflict, message)
}

func NewForbiddenErrorWithMessage(message string) AppError {
	err := errors.New(message)
	return NewAppError(err, http.StatusForbidden, message)
}

func NewUnAuthorizedErrorWithMessage(message string) AppError {
	err := errors.New(message)
	return NewAppError(err, http.StatusUnauthorized, message)
}

func HandleBindError(c *gin.Context, err error, keyword string) {
	traceID := tracking.GetTrackIDFromContext(c)
	logger.NewLogger().WithFields(
		logrus.Fields{"keyword": fmt.Sprintf("%s HandleBindError", keyword),
			constants.TrackIDHeader: traceID,
		}).
		WithStatusCode(http.StatusBadRequest).
		WithError(err).
		Errorln()
	response := NewErrorResponse(c, AppError{
		StatusCode: http.StatusBadRequest,
		Message:    constants.BadRequestErrMess,
		TraceID:    traceID,
	})
	c.AbortWithStatusJSON(http.StatusOK, response)
}

func HandleError(c *gin.Context, err error) {
	traceID := tracking.GetTrackIDFromContext(c)
	code := http.StatusInternalServerError
	defer func() {
		if code >= http.StatusNotFound {
			logger.NewLogger().WithFields(
				logrus.Fields{"keyword": "Response AppError",
					constants.TrackIDHeader: traceID,
				}).
				WithStatusCode(code).
				WithError(err).
				Errorln()
		}
	}()
	if ae, ok := err.(AppError); ok {
		ae.TraceID = traceID
		response := NewErrorResponse(c, ae)
		code = ae.StatusCode
		c.AbortWithStatusJSON(ae.StatusCode, response)
		return
	}

	if ae, ok := err.(*AppError); ok {
		ae.TraceID = traceID
		response := NewErrorResponse(c, *ae)
		code = ae.StatusCode
		c.AbortWithStatusJSON(ae.StatusCode, response)
		return
	}

	// c.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError(err))
	c.AbortWithStatusJSON(http.StatusOK, NewErrorResponse(c, AppError{
		Message:    err.Error(),
		Code:       0,
		StatusCode: http.StatusInternalServerError,
	}))
}

func HandleErrorWithInput(c *gin.Context, input interface{}, err error, searchKey string) {
	code := http.StatusInternalServerError
	defer func() {
		if code >= http.StatusNotFound {
			logger.NewLogger().WithFields(
				logrus.Fields{"keyword": fmt.Sprintf("HandleErrorWithInput Response AppError %s", searchKey),
					constants.TrackIDHeader: tracking.GetTrackIDFromContext(c),
				}).
				WithStatusCode(code).
				WithInput(input).
				WithError(err).
				Error()
		}
	}()
	if ae, ok := err.(AppError); ok {
		response := NewErrorResponse(c, ae)
		code = ae.StatusCode
		c.AbortWithStatusJSON(ae.StatusCode, response)
		return
	}

	if ae, ok := err.(*AppError); ok {
		response := NewErrorResponse(c, *ae)
		code = ae.StatusCode
		c.AbortWithStatusJSON(ae.StatusCode, response)
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError(err))
}

func HandlePathNotExistError(c *gin.Context, err error) {
	traceID := tracking.GetTrackIDFromContext(c)
	logger.NewLogger().WithFields(
		logrus.Fields{constants.TrackIDHeader: traceID}).
		WithStatusCode(http.StatusNotFound).
		WithError(err).
		Errorln()
	response := NewErrorResponse(c, AppError{
		StatusCode: http.StatusNotFound,
		Message:    constants.NotExistPathErrMess,
		TraceID:    traceID,
	})
	c.AbortWithStatusJSON(http.StatusNotFound, response)
}
