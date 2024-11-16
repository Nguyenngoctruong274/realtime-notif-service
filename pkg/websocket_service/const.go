package websocket_service

import "time"

const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
	update      = "update"
	add         = "add"
)

// constants for server message
const (
	errInvalidMessage       = "Server: Invalid msg"
	errActionUnrecognizable = "Server: Action unrecognized"
)

const (
	// time to read the next client's pong message
	PongWait = 60 * time.Second
	// time period to send pings to client
	PingPeriod = (PongWait * 9) / 10
	// time allowed to write a message to client
	WriteWait = 60 * time.Second
	// max message size allowed
	MaxMessageSize = 1024
	// I/O read buffer size
	ReadBufferSize = 1024
	// I/O write buffer size
	WriteBufferSize = 1024
)
