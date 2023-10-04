package messaging

type MessageHandler interface {
	HandleMessage(*Context)
}

type HandlerFunc func(ctx *Context)

func (hf HandlerFunc) HandleMessage(ctx *Context) {
	hf(ctx)
}
