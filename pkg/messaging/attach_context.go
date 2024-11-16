package messaging

func AttachEnhancedContext() HandlerFunc {
	return func(ctx *Context) {
		// AttachQueueContext(ctx)
		ctx.Next()
	}
}
