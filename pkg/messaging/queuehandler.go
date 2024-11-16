package messaging

import (
	"context"
	"yes4all/ads-noti-api/pkg/messaging/queueadapter"
)

type MessageHandler interface {
	Handler
	Subscribe(context context.Context) error
	Close() error
	IsClosed() bool
}

type queueHandler struct {
	Handler
	consumer queueadapter.Consumer

	isClosed bool
	closeChn chan bool
}

func NewQueueHandler(consumer queueadapter.Consumer) MessageHandler {
	handler := NewHandlerBuilder().BuildDefault()

	return &queueHandler{
		Handler:  handler,
		consumer: consumer,
		closeChn: make(chan bool),
	}
}

func (q *queueHandler) Subscribe(ctx context.Context) error {
	chn, err := q.consumer.Subscribe(ctx)
	if err != nil {
		return err
	}

	go func() {
		for !q.isClosed {
			select {
			case msg := <-chn:
				q.Handle(msg.Topic, msg.Data, msg.Headers, msg.ContentType())
			case <-q.closeChn:
			}
		}
	}()

	return nil
}

func (q *queueHandler) IsClosed() bool {
	return q.isClosed
}

func (q *queueHandler) Close() error {
	if q.isClosed {
		return nil
	}

	q.isClosed = true
	q.closeChn <- true
	return q.consumer.Close()
}
