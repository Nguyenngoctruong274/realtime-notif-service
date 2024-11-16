package xerrors

const (
	// Error for cache type -50 -> -100
	CannotGetFromCache    ErrorType = -50
	CannotSaveToCache     ErrorType = -51
	CannotDeleteFromCache ErrorType = -52
	CannotIncrInCache     ErrorType = -53
	CannotDecrInCache     ErrorType = -54
	GetAllKeyRedisError   ErrorType = -55
)
