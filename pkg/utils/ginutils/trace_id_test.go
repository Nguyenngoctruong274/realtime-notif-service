package ginutils

import (
	"fmt"
	"testing"

	"github.com/rs/xid"
)

func TestTraceID(t *testing.T) {
	id := xid.New()
	id.Time()
	traceID := fmt.Sprintf("%s_%s", id.Time().Format(formatDate), id.String())
	t.Log(traceID)
}
