package messaging

import uuid "github.com/satori/go.uuid"

type Relay interface {
	//Relay 转送消息
	Relay(*Context)
}

// SimpleRelay 直接将消息发送打Outbound Pipeline
type SimpleRelay struct {
}

func (relay *SimpleRelay) Relay(ctx *Context) {
	relay.convertFrame(ctx)
	ctx.broker.outboundPipe.Process(ctx)
}

func (relay *SimpleRelay) convertFrame(ctx *Context) {
	if ctx.Frame.Command == MESSAGE {
		return
	}
	ctx.Frame.Command = MESSAGE
	ctx.Frame.Headers["message-id"] = uuid.NewV1().String()
}
