package ginutils

import (
	"bytes"
	"net/http/httputil"
	"time"
	"yes4all/ads-noti-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	_, err := w.body.Write(b)
	if err != nil {
		return 0, err
	}

	return w.ResponseWriter.Write(b)
}

// MiddlewareLogger for log request and response
func Logger(skipPaths ...string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(skipPaths); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range skipPaths {
			skip[path] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		lg := logger.NewLogger().WithKeyword(c, "--) [LOGGER] API middleware logger")
		start := time.Now()
		path := c.FullPath()
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		if _, ok := skip[path]; !ok {
			reqDump, err := httputil.DumpRequest(c.Request, true)
			if err != nil {
				lg.Error()
				return
			}
			reqDumpStr := string(reqDump)
			lg.WithFields(logrus.Fields{
				"request":    reqDumpStr,
				"response":   blw.body.String(),
				"latency.ms": time.Since(start).Milliseconds(),
				"status":     c.Writer.Status(),
				"ip":         c.ClientIP(),
			}).Info()
		}
	}
}
