package errors

import (
	"errors"
	"yes4all/ads-noti-api/pkg/utils/constants"
)

var (
	ErrBadRequest   = errors.New(constants.BadRequestErrMess)
	ErrNotFound     = errors.New(constants.NotFoundErrMess)
	ErrInternal     = errors.New(constants.InternalServerErrMess)
	ErrRedis        = errors.New(constants.InternalServerErrMess)
	ErrNotExistPath = errors.New(constants.NotExistPathErrMess)
)
