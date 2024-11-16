package queueadapter

import "context"

//go:generate mockgen -source=$GOFILE -package=gomock -destination=$PWD/mocks/${GOFILE}

type Producer interface {
	Send(ctx context.Context, contentType string, body interface{}, opts ...SendOpt) error
	Close() error
}
