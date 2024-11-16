package tracing

type SpanType int32

const (
	SpanTypeUnknown SpanType = iota
	SpanTypeMiddleware
	SpanTypeDB
	SpanTypeCache
	SpanTypeHTTP
	SpanTypeValidate
	SpanTypeBindData
	SpanTypeQueue
)

var spanTypeMap = map[SpanType]string{
	SpanTypeUnknown:    "unknown",
	SpanTypeMiddleware: "middleware",
	SpanTypeDB:         "db",
	SpanTypeCache:      "cache",
	SpanTypeHTTP:       "http",
	SpanTypeValidate:   "validate",
	SpanTypeBindData:   "bind_data",
	SpanTypeQueue:      "queue",
}

func (s SpanType) String() string {
	return spanTypeMap[s]
}

// ------------------ Transaction Type -------------
type TransactionType int32

const (
	TransTypeUnknown TransactionType = iota
	TransTypeRequest
	TransTypeDB
	TransTypeExternalHTTP
)

var transTypeMap = map[TransactionType]string{
	TransTypeUnknown:      "unknown",
	TransTypeRequest:      "request",
	TransTypeDB:           "db",
	TransTypeExternalHTTP: "external_http",
}

func (t TransactionType) String() string {
	return transTypeMap[t]
}
