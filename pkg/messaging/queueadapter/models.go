package queueadapter

type SendOpt func(message *Message) error

const (
	ContentType = "Content-Type"
	CreatedAt   = "Created-At"
	MessageID   = "X-Message-ID"
)

type Message struct {
	Topic     string            `json:"topic"`
	Data      []byte            `json:"data"`
	Headers   map[string]string `json:"headers"`
	CreatedAt int64             `json:"created_at"`
}

func (m Message) ContentType() string {
	return m.Headers[ContentType]
}
