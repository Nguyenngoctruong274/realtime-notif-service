package messaging

import "os"

type handlerBuilder struct {
	handlerFns []HandlerFunc
}

func NewHandlerBuilder() *handlerBuilder {
	return &handlerBuilder{}
}

func (b *handlerBuilder) withRecovery() *handlerBuilder {
	b.handlerFns = append(b.handlerFns, Recovery(os.Stderr))
	return b
}

func (b *handlerBuilder) AddHandler(handlerFunc HandlerFunc) *handlerBuilder {
	b.handlerFns = append(b.handlerFns, handlerFunc)
	return b
}

func (b *handlerBuilder) BuildDefault() Handler {
	b.withRecovery()

	handler := NewHandler()
	handler.Use(b.handlerFns...)

	return handler
}

func (b *handlerBuilder) Build() Handler {
	handler := NewHandler()
	handler.Use(b.handlerFns...)
	return handler
}
