package grpcutils

import (
	"errors"
	"time"
	"yes4all/ads-noti-api/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	waitTime     = 100 * time.Millisecond
	maxWaitCount = 100
)

func WaitForReady(conn *grpc.ClientConn) error {
	count := 0
	for conn.GetState() != connectivity.Ready {
		logger.NewLogger().Debugf("grpc conn status", conn.GetState())
		if count >= maxWaitCount {
			return errors.New("exceed max connection retry")
		}
		time.Sleep(waitTime)
		count++
	}
	return nil
}
