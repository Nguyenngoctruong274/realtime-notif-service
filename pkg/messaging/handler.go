package messaging

type Handler interface {
	Handle(queueName string, data []byte, headers map[string]string, contentType string)
	Use(handlerFns ...HandlerFunc)
}

type HandlerFunc func(ctx *Context)

type handler struct {
	handlerFns []HandlerFunc
}

func NewHandler() Handler {
	return &handler{
		handlerFns: make([]HandlerFunc, 0),
	}
}

func (h *handler) Use(handlerFns ...HandlerFunc) {
	h.handlerFns = append(h.handlerFns, handlerFns...)
}

func (h *handler) Handle(queueName string, data []byte, headers map[string]string, contentType string) {
	ctx := NewContext(h.handlerFns, queueName, data, headers, contentType)
	ctx.Next()
}
