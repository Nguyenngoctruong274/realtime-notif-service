package messaging

import (
	"context"
)

type contextWrapper struct {
	ctx *Context
}

func NewContextWrapper(ctx *Context) *contextWrapper {
	return &contextWrapper{ctx}
}

func (m *contextWrapper) WithValue(key string,
	value interface{}) *contextWrapper {
	m.ctx.Set(key, value)
	return m
}

func (m *contextWrapper) GetHeader(key string) string {
	value, ok := m.ctx.Headers()[key]
	if !ok {
		return ""
	}

	return value
}

func (m *contextWrapper) Context() context.Context {
	return m.ctx
}
