package xprom

const (
	// biz's present for general flow, such as: checkout, cart, ...
	bizLabel = "biz"
	// event's present for sub flow, ex: biz = checkout, event = save_order
	eventLabel = "event"
	// event's present for specific flow, ex: biz = checkout, event = save_order, action = call_api_save_order.
	actionLabel = "action"
)

type Biz string

const (
	BizUnknown Biz = "unknown"
	BizGeneral Biz = "general"
	BizESAPI   Biz = "es-api"
)
