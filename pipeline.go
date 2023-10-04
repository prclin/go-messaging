package messaging

// Pipeline frame管道
//
// 对frame作处理
type Pipeline interface {
	//Process 处理frame
	Process(ctx *Context)

	//AddValves 向pipeline中添加valve
	AddValves(...Valve)
}

// InboundPipeline 客户端发送的frame作预处理
type InboundPipeline struct {
	//处理阀
	valveChain Valve
}

func DefaultInboundPipeline() *InboundPipeline {
	ip := &InboundPipeline{valveChain: new(ReceiptValve)}
	ip.AddValves(new(SendValve), new(SubscribeValve), new(UnsubscribeValve), new(DisconnectValve))
	return ip
}

func (pipeline *InboundPipeline) AddValves(valves ...Valve) {
	for _, valve := range valves {
		valve.SetNext(pipeline.valveChain)
		pipeline.valveChain = valve
	}
}

func (pipeline *InboundPipeline) Process(context *Context) {
	//阀处理
	if pipeline.valveChain != nil {
		pipeline.valveChain.Valve(context)
	}
}

// OutboundPipeline 把消息发送到客户端
type OutboundPipeline struct {
	//处理阀
	valveChain Valve
}

func (pipeline *OutboundPipeline) AddValves(valves ...Valve) {
	for _, valve := range valves {
		valve.SetNext(pipeline.valveChain)
		pipeline.valveChain = valve
	}
}

func (pipeline *OutboundPipeline) Process(ctx *Context) {
	for _, subscription := range ctx.broker.subscriptions {
		if subscription.Destination == ctx.Frame.Headers["destination"] {
			ctx.Frame.Headers["subscription"] = subscription.Id
			subscription.Conn.WriteFrame(ctx.Frame)
		}
	}
}

// BrokerPipeline 对系统发出的消息作预处理
type BrokerPipeline struct {
	//处理阀
	valveChain Valve
}

func (pipeline *BrokerPipeline) AddValves(valves ...Valve) {
	for _, valve := range valves {
		valve.SetNext(pipeline.valveChain)
		pipeline.valveChain = valve
	}
}

func (pipeline *BrokerPipeline) Process(ctx *Context) {
	ctx.broker.relay.Relay(ctx)
}
