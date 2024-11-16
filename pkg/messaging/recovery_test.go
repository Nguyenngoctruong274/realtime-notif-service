package messaging

import (
	"encoding/json"
	"fmt"
	"testing"
	"yes4all/ads-noti-api/pkg/messaging/contenttype"
	"yes4all/ads-noti-api/pkg/messaging/queueadapter"

	"github.com/stretchr/testify/assert"
)

type sample struct {
	Data string
}

func TestRecover(t *testing.T) {
	consumer := queueadapter.NewEmptyConsumer()
	handler := NewQueueHandler(consumer)
	var ctx *Context
	handler.Use(func(c *Context) {
		ctx = c
		c.Next()
	})

	writer := &dumpWriter{}
	handler.Use(Recovery(writer), divisionByZero)
	sampleJSON, _ := json.Marshal(sample{Data: "hello world"})
	queueName := "testqueue"
	headers := map[string]string{
		"content-type": "json",
		"created-at":   "1563845487",
	}
	handler.Handle(queueName, sampleJSON, headers, contenttype.JSON)
	assert.Error(t, ctx.err)
	t.Log(ctx.err)
	// assert.Equal(t, errors.InternalServerError, errors.GetType(ctx.err))
	output := string(writer.bytes)
	for key, value := range headers {
		text := key + ": " + value
		assert.Contains(t, output, text, "not contains %s", text)
	}
	assert.Contains(t, output, "queue:"+queueName, "not contains queue name")
	assert.Contains(t, output, string(sampleJSON), "not contains body")
}

func divisionByZero(ctx *Context) {
	// zero := 1
	zero := 0
	fmt.Println("division by zero", 4/zero)
}

type dumpWriter struct {
	bytes []byte
}

func (w *dumpWriter) Write(p []byte) (n int, err error) {
	w.bytes = append(w.bytes, p...)
	return
}
