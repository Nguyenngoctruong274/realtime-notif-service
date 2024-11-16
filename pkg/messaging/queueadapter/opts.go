package queueadapter

import (
	"context"
	"encoding/json"
	"errors"
	"yes4all/ads-noti-api/pkg/messaging/contenttype"

	"google.golang.org/protobuf/proto"
)

func MessageValue(body interface{}) SendOpt {
	return func(m *Message) (err error) {
		contentType, ok := m.Headers[ContentType]
		if !ok {
			return errors.New("invalid content type")
		}
		var value []byte
		switch contentType {
		case contenttype.BINARY:
			value = body.([]byte)
		case contenttype.JSON:
			value, err = json.Marshal(body)
			if err != nil {
				return err
			}
		case contenttype.PROTOBUF:
			e, ok := body.(proto.Message)
			if !ok {
				return errors.New("message must be protobuf type")
			}
			value, err = proto.Marshal(e)
			if err != nil {
				return err
			}
		default:
			return errors.New("invalid content type")
		}

		m.Data = value
		return nil
	}
}

func MessageHeaders(headers map[string]string) SendOpt {
	return func(m *Message) error {
		m.Headers = headers
		return nil
	}
}

func MessageAddHeader(header, value string) SendOpt {
	return func(m *Message) error {
		if m.Headers == nil {
			m.Headers = make(map[string]string)
		}

		m.Headers[header] = value
		return nil
	}
}

func MessageContentType(contentType string) SendOpt {
	return MessageAddHeader(ContentType, contentType)
}

func MessageContextID(ctx context.Context) SendOpt {
	return func(message *Message) error {
		if contextID, ok := ctx.Value(MessageID).(string); ok {
			return MessageAddHeader(MessageID, contextID)(message)
		}
		return nil
	}
}

//
// func InjectSpanContext(spanContext opentracing.SpanContext) SendOpt {
// 	return func(m *Message) error {
// 		if m.Headers == nil {
// 			m.Headers = make(map[string]string)
// 		}
// 		return opentracing.GlobalTracer().Inject(spanContext, opentracing.TextMap, carrier(*m))
// 	}
// }
