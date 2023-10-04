package messaging

type Context struct {
	Frame  *Frame
	Params map[string]string
	broker *StompBroker
	conn   *Conn
}

func NewContext(frame *Frame, conn *Conn) *Context {
	return &Context{Frame: frame, conn: conn}
}

func (context *Context) handle() {
	handler := context.broker.getMessageHandler(context)
	handler.HandleMessage(context)
}

func (context *Context) Send(destination string, body []byte) {
	headers := make(map[string]string)
	headers["destination"] = destination
	frame := NewFrame(MESSAGE, headers, body)
	context.broker.brokerPipe.Process(NewContext(frame, nil))
}
