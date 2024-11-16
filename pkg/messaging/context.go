package messaging

import (
	"context"
	"math"
)

type Context struct {
	context.Context

	queueName   string
	key         string
	data        []byte
	headers     map[string]string
	contentType string
	middlewares []HandlerFunc
	result      interface{}
	err         error
	isAborted   bool
	abortIndex  int
	index       int
}

func NewContext(middlewares []HandlerFunc,
	queueName string, data []byte,
	headers map[string]string,
	contentType string) *Context {
	return &Context{
		Context:     context.Background(),
		middlewares: middlewares,
		queueName:   queueName,
		data:        data,
		headers:     headers,
		contentType: contentType,
		index:       -1,
		abortIndex:  math.MaxInt64,
	}
}

func (c *Context) Handle() {
	c.Next()
}

func (c *Context) Next() {
	if c.isAborted {
		return
	}

	c.index++
	for s := len(c.middlewares); c.index < s && !c.isAborted; c.index++ {
		c.middlewares[c.index](c)
	}
}

func (c *Context) Abort() {
	c.isAborted = true
	c.abortIndex = c.index
}

func (c *Context) AbortWithError(err error) {
	c.isAborted = true
	c.abortIndex = c.index
	c.err = err
}

func (c *Context) AbortWithResult(err error, result interface{}) {
	c.AbortWithError(err)
	c.result = result
}

func (c *Context) ReportResult(result interface{}) {
	c.result = result
}

func (c *Context) IsAborted() bool {
	return c.isAborted
}

func (c *Context) QueueName() string {
	return c.queueName
}

func (c *Context) Data() []byte {
	return c.data
}

func (c *Context) Result() interface{} {
	return c.result
}

func (c *Context) Headers() map[string]string {
	return c.headers
}

func (c *Context) Err() error {
	return c.err
}

func (c *Context) Key() string {
	return c.key
}

func (c *Context) Set(key, value interface{}) {
	if c.Context == nil {
		c.Context = context.Background()
	}
	c.Context = context.WithValue(c.Context, key, value)
}

func (c *Context) ContentType() string {
	return c.contentType
}
