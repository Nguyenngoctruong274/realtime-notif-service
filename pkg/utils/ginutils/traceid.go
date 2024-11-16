package ginutils

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"yes4all/ads-noti-api/pkg/xcontext"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

const (
	headerRequestID = "X-Request-ID"
	formatDate      = "0601021504" // YYMMDDHHMM
)

func InjectTraceID(c *gin.Context) {
	id := xid.New()
	uuidID := uuid.New()
	// Format: YYMMDDHHMM_ID => 210914_0050_c4voumo6n88nq9t8dk20
	traceID := fmt.Sprintf("%s_%s", id.Time().Format(formatDate), uuidID.String())
	rid := c.GetHeader(headerRequestID)
	if rid == "" {
		rid = traceID
	}

	ctx := context.WithValue(c.Request.Context(), xcontext.KeyContextID.String(), rid) //nolint: staticcheck
	c.Request = c.Request.WithContext(ctx)
	c.Set(xcontext.KeyContextID.String(), rid)
	c.Next()
}

func GetTraceIDFromCtx(c *gin.Context) string {
	if result, ok := c.Get(xcontext.KeyContextID.String()); ok {
		_, ok := result.(string)
		if !ok {
			return ""
		}
	}
	return ""
}
