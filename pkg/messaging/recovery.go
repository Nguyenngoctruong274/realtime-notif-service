package messaging

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"yes4all/ads-noti-api/pkg/utils/stacktrace"
)

const (
	stackSkip = 3
)

var (
	reset = string([]byte{27, 91, 48, 109})
)

func DefaultRecovery() HandlerFunc {
	return Recovery(os.Stderr)
}

func Recovery(out io.Writer) HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var osErr *os.SyscallError
					if ok := errors.As(ne.Err, &osErr); ok {
						if strings.Contains(strings.ToLower(osErr.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(osErr.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if logger != nil {
					stack := stacktrace.Stack(stackSkip)
					data := string(c.data)
					headers := make([]string, 0, len(c.headers))
					for key, value := range c.headers {
						headers = append(headers, key+": "+value)
					}
					if brokenPipe {
						logger.Printf("%s\nqueue:%s\ndata:%s%s", err, c.queueName, data, reset)
					} else {
						logger.Printf("[Recovery] %s panic recovered:\nqueue:%s\nheaders:%s\n%s\n%s\t%s%s",
							timeFormat(time.Now()), c.queueName, strings.Join(headers, "\r\n"), err,
							data, stack, reset)
					}
				}

				// If the connection is dead, we can't write a status to it.
				c.AbortWithError(err.(error))
			}
		}()
		c.Next()
	}
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
