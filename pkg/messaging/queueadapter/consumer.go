package queueadapter

import (
	"context"
	"yes4all/ads-noti-api/pkg/logger"
)

//go:generate mockgen -source=$GOFILE -package=gomock -destination=$PWD/mocks/${GOFILE}

type Consumer interface {
	Subscribe(ctx context.Context) (<-chan Message, error)
	Close() error
}

type EmptyConsumer struct {
	outCh    chan Message
	isClosed bool
	closeChn chan bool
}

func NewEmptyConsumer() *EmptyConsumer {
	return &EmptyConsumer{
		outCh:    make(chan Message),
		closeChn: make(chan bool),
	}
}

func (c *EmptyConsumer) Subscribe(_ context.Context) (<-chan Message, error) {
	go func() {
		for !c.isClosed {
			select {
			case msg := <-c.outCh:
				c.outCh <- msg
			case <-c.closeChn:
			}
		}
	}()

	return c.outCh, nil
}

func (c *EmptyConsumer) Send(ctx context.Context, contentType string, body interface{},
	optFns ...SendOpt) (err error) {
	msg := &Message{}
	optFns = append(optFns,
		MessageContentType(contentType),
		MessageValue(body))
	for _, fn := range optFns {
		err = fn(msg)
		if err != nil {
			return err
		}
	}

	logger.NewLogger().Debug("send msg")
	c.outCh <- *msg
	logger.NewLogger().Debug("done send msg")

	return nil
}

func (c *EmptyConsumer) Close() error {
	if c.isClosed {
		return nil
	}
	c.isClosed = true
	close(c.outCh)
	c.closeChn <- true

	return nil
}
