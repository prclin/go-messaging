package messaging

import "strings"

// Valve frame阀
//
// 对接收到的frame作具体处理
type Valve interface {
	//Valve 处理函数
	Valve(ctx *Context) error
	//SetNext 设置下一个valve
	SetNext(next Valve)
	//Next 返回下一个处理器
	Next() Valve
}

// StandardValve 标准valve，起到抽象类的作用
type StandardValve struct {
	next Valve
}

func (valve *StandardValve) SetNext(nextValve Valve) {
	valve.next = nextValve
}

func (valve *StandardValve) Next() Valve {
	return valve.next
}

// ReceiptValve 作为根阀，对SEND frame做消息确认
type ReceiptValve struct {
	StandardValve
}

func (valve *ReceiptValve) Valve(ctx *Context) error {
	if ctx.Frame.Command != SEND {
		return nil
	}
	//发送RECEIPT frame
	ctx.conn.Write(
		RECEIPT,
		nil,
		"receipt-id", ctx.Frame.Headers["receipt"],
	)
	return nil
}

// SendValve 处理Send frame
type SendValve struct {
	StandardValve
}

func (valve *SendValve) Valve(ctx *Context) error {
	if ctx.Frame.Command != SEND {
		return valve.Next().Valve(ctx)
	}

	if strings.HasPrefix(ctx.Frame.Headers["destination"], ctx.broker.AppDestinationPrefix) {
		ctx.handle()
	} else if strings.HasPrefix(ctx.Frame.Headers["destination"], ctx.broker.BrokerDestinationPrefix) {
		//发送到消息中继
		ctx.broker.relay.Relay(ctx)
	}

	return nil
}

// SubscribeValve 处理Subscribe frame
type SubscribeValve struct {
	StandardValve
}

func (valve *SubscribeValve) Valve(ctx *Context) error {
	if ctx.Frame.Command != SUBSCRIBE {
		return valve.Next().Valve(ctx)
	}
	ctx.broker.addSubscription(&Subscription{Id: ctx.Frame.Headers["id"], Destination: ctx.Frame.Headers["destination"], Conn: ctx.conn})
	ctx.handle()
	return nil
}

// UnsubscribeValve 处理Unsubscribe frame
type UnsubscribeValve struct {
	StandardValve
}

func (valve *UnsubscribeValve) Valve(ctx *Context) error {

	return valve.Next().Valve(ctx)
}

type DisconnectValve struct {
	StandardValve
}

func (valve *DisconnectValve) Valve(ctx *Context) error {
	ctx.conn.disconnected = true
	return valve.Next().Valve(ctx)
}
