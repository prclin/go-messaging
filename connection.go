package messaging

import (
	"errors"
	"github.com/gorilla/websocket"
)

// Conn stomp连接
type Conn struct {
	ProtocolResolver                 //协议解析器
	conn             *websocket.Conn //websocket连接
	disconnected     bool            //是否接收到disconnect帧
}

func NewConn(conn *websocket.Conn) *Conn {
	return &Conn{conn: conn, ProtocolResolver: defaultProtocolResolver}
}

// ReadFrame 读取stomp帧
//
// 此方法会阻塞，直到有可读的stomp帧
//
// 当读取消息错误、消息类型不支持或解析帧失败时返回error
func (c *Conn) ReadFrame() (*Frame, error) {
	//读取消息
	messageType, message, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	//目前只支持text message
	if messageType != websocket.TextMessage {
		return nil, errors.New("unsupported message type")
	}

	//解析frame
	frame, err := c.Parse(message)
	if err != nil {
		return nil, err
	}
	return frame, nil
}

// WriteFrame 写stomp帧
//
// 目前只支持写文本消息类型
func (c *Conn) WriteFrame(frame *Frame) error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte(frame.String()))
}

// Write 写stomp帧
func (c *Conn) Write(command Command, body []byte, headerKeysAndValues ...string) error {
	return c.WriteFrame(NewFrame(command, sliceToMap(headerKeysAndValues...), body))
}

// Close 关闭stomp连接
func (c *Conn) Close() error {
	return c.conn.Close()
}

// sliceToMap 将切片转为map
//
// 如果切片长度是奇数，则最后一个key的值为零值
func sliceToMap(slice ...string) map[string]string {
	m := make(map[string]string, len(slice))
	for i := 0; i < len(slice); i += 2 {
		if i+1 > len(slice) {
			m[slice[i]] = ""
		} else {
			m[slice[i]] = slice[i+1]
		}
	}
	return m
}
