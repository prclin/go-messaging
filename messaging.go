package messaging

import (
	"net/http"
	"sync"
)

// StompBroker 消息代理
//
// 数据流向
//
/*
InboundChannel	---->  MethodMessageHandler

		|       ---->  BrokerMessageHandler	<----       |

OutboundChannel	<----         |
*/
type StompBroker struct {
	Upgrader
	Router
	relay                   Relay           //消息中继
	inboundPipe             Pipeline        //接收消息
	outboundPipe            Pipeline        //发送消息
	brokerPipe              Pipeline        //接收内部消息
	AppDestinationPrefix    string          //应用消息前缀
	BrokerDestinationPrefix string          //代理消息前缀
	lock                    sync.Mutex      //锁
	subscriptions           []*Subscription //订阅
}

func NewStompBroker() *StompBroker {
	broker := &StompBroker{}
	broker.Upgrader = NewUpgrader()
	broker.inboundPipe = DefaultInboundPipeline()
	broker.Router = *(NewRouter())
	broker.brokerPipe = new(BrokerPipeline)
	broker.relay = new(SimpleRelay)
	broker.outboundPipe = new(OutboundPipeline)
	broker.subscriptions = make([]*Subscription, 0)
	return broker
}

// ServeOverHttp 与client建立连接，建立成功后会阻塞在这，直到发生错误或者连接中断
func (sb *StompBroker) ServeOverHttp(w http.ResponseWriter, r *http.Request) error {
	//升级为stomp协议
	conn, err := sb.Upgrade(w, r)
	if err != nil {
		return err
	}
	//监听消息
	for {
		frame, err := conn.ReadFrame()
		if err != nil {
			conn.Close()
			return err
		}
		sb.inboundPipe.Process(NewContext(frame, conn, sb))
	}
}

func (sb *StompBroker) addSubscription(subscription *Subscription) {
	//sb.lock.Lock()
	//defer sb.lock.Unlock()
	sb.subscriptions = append(sb.subscriptions, subscription)
}

type Subscription struct {
	Id          string
	Destination string
	Conn        *Conn
}
